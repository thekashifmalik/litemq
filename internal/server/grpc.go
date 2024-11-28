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
	data     [][]byte
	channels []chan []byte
	lock     sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{}
}

func NewServer() *Server {
	slog.Info("litemq started")
	return &Server{
		queues: map[string]*Queue{},
	}
}

func (s *Server) Enqueue(ctx context.Context, request *gen.EnqueueRequest) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("ENQ %v '%v'", request.Queue, string(request.Data)))
	queue := s.getOrCreateQueue(request.Queue)
	queue.lock.Lock()
	if len(queue.channels) > 0 {
		channel := queue.channels[0]
		queue.channels = queue.channels[1:]
		channel <- request.Data
		close(channel)
	} else {
		queue.data = append(queue.data, request.Data)
	}
	queue.lock.Unlock()
	return &gen.QueueLength{Count: int64(len(queue.data))}, nil
}

func (s *Server) Dequeue(ctx context.Context, request *gen.QueueID) (*gen.DequeueResponse, error) {
	slog.Info(fmt.Sprintf("DEQ %v", request.Queue))
	queue := s.getOrCreateQueue(request.Queue)
	queue.lock.Lock()
	if len(queue.data) > 0 {
		data := queue.data[0]
		queue.data = queue.data[1:]
		queue.lock.Unlock()
		slog.Debug(fmt.Sprintf("< %v bytes", int64(len(data))))
		return &gen.DequeueResponse{Data: data}, nil
	}
	channel := make(chan []byte)
	queue.channels = append(queue.channels, channel)
	queue.lock.Unlock()

	select {
	case data := <-channel:
		slog.Debug(fmt.Sprintf("< %v bytes", int64(len(data))))
		return &gen.DequeueResponse{Data: data}, nil
	case <-ctx.Done():
		// TODO: Figure out if there is a race-condition here with the context channel select and the queue lock being
		// acquired. If a concurrent enqueue request acquires the lock between these operations, it can write a message
		// to the  channel which will then be closed here.
		queue.lock.Lock()
		channels := []chan []byte{}
		for _, ch := range queue.channels {
			if ch != channel {
				channels = append(channels, ch)
			}
		}
		queue.channels = channels
		queue.lock.Unlock()
		close(channel)
		slog.Debug("< disconnected")
		return nil, ctx.Err()
	}
}

func (s *Server) getOrCreateQueue(name string) *Queue {
	queue, ok := s.queues[name]
	if !ok {
		queue = s.acquireLockAndCreateQueue(name)
	}
	return queue
}

func (s *Server) acquireLockAndCreateQueue(name string) *Queue {
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
	slog.Info(fmt.Sprintf("DEL %v", request.Queue))
	s.lock.Lock()
	length := 0
	queue, ok := s.queues[request.Queue]
	if ok {
		delete(s.queues, request.Queue)
		length = len(queue.data)
	}
	s.lock.Unlock()
	return &gen.QueueLength{Count: int64(length)}, nil
}

func (s *Server) Length(ctx context.Context, request *gen.QueueID) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("LEN %v", request.Queue))
	queue, _ := s.queues[request.Queue]
	return &gen.QueueLength{Count: int64(len(queue.data))}, nil
}

func (s *Server) Health(ctx context.Context, request *gen.Nothing) (*gen.Nothing, error) {
	slog.Info(fmt.Sprintf("HEALTH"))
	return &gen.Nothing{}, nil
}
