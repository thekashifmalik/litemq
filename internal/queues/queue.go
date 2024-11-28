package queues

import (
	"context"
	"sync"
)

type Queue struct {
	messages [][]byte
	channels []chan []byte
	lock     sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) LockAndEnqueue(msg []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.channels) > 0 {
		channel := q.channels[0]
		q.channels = q.channels[1:]
		channel <- msg
		close(channel)
	} else {
		q.messages = append(q.messages, msg)
	}
}

func (q *Queue) LockAndDequeue(ctx context.Context) ([]byte, error) {
	data, channel := q.lockAndDequeueOrChannel()
	if channel == nil {
		return data, nil
	}
	select {
	case data := <-channel:
		return data, nil
	case <-ctx.Done():
		// TODO: Figure out if there is a race-condition here with the context channel select and the queue lock being
		// acquired. If a concurrent enqueue request acquires the lock between these operations, it can write a message
		// to the  channel which will then be closed here.
		q.lockAndDisconnect(channel)
		return nil, ctx.Err()
	}
}

func (q *Queue) lockAndDequeueOrChannel() ([]byte, chan []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.messages) > 0 {
		msg := q.messages[0]
		q.messages = q.messages[1:]
		return msg, nil
	}
	channel := make(chan []byte)
	q.channels = append(q.channels, channel)
	return nil, channel
}

func (q *Queue) lockAndDisconnect(channel chan []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()
	channels := []chan []byte{}
	for _, ch := range q.channels {
		if ch != channel {
			channels = append(channels, ch)
		}
	}
	q.channels = channels
	// TODO: Maybe we should close the channel earlier to avoid any race-conditions. See note in LockAndDequeue.
	close(channel)
}

func (q *Queue) Length() int {
	return len(q.messages)
}
