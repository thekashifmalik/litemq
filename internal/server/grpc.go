package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/thekashifmalik/litemq/gen"
)

type Server struct {
	gen.UnimplementedLiteMQServer

	queues map[string]*Queue
	lock   sync.Mutex
}

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

func NewServer() *Server {
	slog.Info("LiteMQ started")
	return &Server{
		queues: map[string]*Queue{},
	}
}

func (s *Server) Enqueue(ctx context.Context, request *gen.EnqueueRequest) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("ENQUEUE %v '%v'", request.Queue, string(request.Data)))
	queue := s.getOrCreateQueue(request.Queue)
	queue.LockAndEnqueue(request.Data)
	return &gen.QueueLength{Count: int64(len(queue.messages))}, nil
}

func (s *Server) Dequeue(ctx context.Context, request *gen.QueueID) (*gen.DequeueResponse, error) {
	slog.Info(fmt.Sprintf("DEQUEUE %v", request.Queue))
	queue := s.getOrCreateQueue(request.Queue)
	data, channel := queue.LockAndDequeueOrChannel()
	if channel == nil {
		slog.Debug(fmt.Sprintf("< %v bytes", int64(len(data))))
		return &gen.DequeueResponse{Data: data}, nil
	}
	select {
	case data := <-channel:
		slog.Debug(fmt.Sprintf("< %v bytes", int64(len(data))))
		return &gen.DequeueResponse{Data: data}, nil
	case <-ctx.Done():
		// TODO: Figure out if there is a race-condition here with the context channel select and the queue lock being
		// acquired. If a concurrent enqueue request acquires the lock between these operations, it can write a message
		// to the  channel which will then be closed here.
		queue.LockAndDisconnect(channel)
		slog.Debug("< disconnected")
		return nil, ctx.Err()
	}
}

func (s *Server) getOrCreateQueue(name string) *Queue {
	queue, ok := s.queues[name]
	if !ok {
		queue = s.lockAndCreateQueue(name)
	}
	return queue
}

func (s *Server) lockAndCreateQueue(name string) *Queue {
	s.lock.Lock()
	defer s.lock.Unlock()
	// Since another request could have created the queue while we were acquiring the lock, we need to check whether the
	// queue exists, even if we checked just before calling this method.
	queue, ok := s.queues[name]
	if !ok {
		queue = NewQueue()
		s.queues[name] = queue
	}
	return queue
}

func (s *Server) Purge(ctx context.Context, request *gen.QueueID) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("PURGE %v", request.Queue))
	s.lock.Lock()
	length := 0
	queue, ok := s.queues[request.Queue]
	if ok {
		delete(s.queues, request.Queue)
		length = len(queue.messages)
		// TODO: Need to delete any active dequeue channels for this purged queue.
	}
	s.lock.Unlock()
	return &gen.QueueLength{Count: int64(length)}, nil
}

func (s *Server) Length(ctx context.Context, request *gen.QueueID) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("LENGTH %v", request.Queue))
	queue, _ := s.queues[request.Queue]
	return &gen.QueueLength{Count: int64(len(queue.messages))}, nil
}

func (s *Server) Health(ctx context.Context, request *gen.Nothing) (*gen.Nothing, error) {
	slog.Info(fmt.Sprintf("HEALTH"))
	return &gen.Nothing{}, nil
}
