package server

import "sync"

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

func (q *Queue) LockAndDequeueOrChannel() ([]byte, chan []byte) {
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

func (q *Queue) LockAndDisconnect(channel chan []byte) {
	q.lock.Lock()
	defer q.lock.Unlock()
	channels := []chan []byte{}
	for _, ch := range q.channels {
		if ch != channel {
			channels = append(channels, ch)
		}
	}
	q.channels = channels
	close(channel)
}
