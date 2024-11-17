package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"

	"github.com/thekashifmalik/kandy/gen"
	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedQueueServiceServer

	queues map[string]*Queue
	lock   sync.Mutex
}

type Queue struct {
	name     string
	data     [][]byte
	channels []chan []byte
	lock     sync.Mutex
}

func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
	}
}

func NewServer() *Server {
	slog.Info("HTTPQ started")
	return &Server{
		queues: map[string]*Queue{},
	}
}

func (s *Server) Enqueue(ctx context.Context, request *gen.EnqueueRequest) (*gen.Nothing, error) {
	slog.Info(fmt.Sprintf("ENQ %v '%v'", request.Queue, string(request.Data)))
	queue := s.getOrCreateQueue(request.Queue)
	queue.lock.Lock()
	if len(queue.channels) > 0 {
		channel := queue.channels[0]
		queue.channels = queue.channels[1:]
		channel <- request.Data
		// TODO: Figure out if we need to close or remove this channel explicitly here. It is removed from the
		// queue.channels slice but donunclear if it will be garbade-collected.
	} else {
		queue.data = append(queue.data, request.Data)
	}
	queue.lock.Unlock()

	return &gen.Nothing{}, nil
}

func (s *Server) Dequeue(ctx context.Context, request *gen.QueueID) (*gen.DequeueResponse, error) {
	slog.Info(fmt.Sprintf("DEQ %v", request.Queue))
	queue := s.getOrCreateQueue(request.Queue)
	queue.lock.Lock()
	if len(queue.data) > 0 {
		data := queue.data[0]
		queue.data = queue.data[1:]
		queue.lock.Unlock()
		return &gen.DequeueResponse{Data: data}, nil
	}
	channel := make(chan []byte)
	queue.channels = append(queue.channels, channel)
	queue.lock.Unlock()

	select {
	case data := <-channel:
		return &gen.DequeueResponse{Data: data}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *Server) getOrCreateQueue(name string) *Queue {
	s.lock.Lock()
	queue, ok := s.queues[name]
	if !ok {
		queue = NewQueue(name)
		s.queues[name] = queue
	}
	s.lock.Unlock()
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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 42090))
	if err != nil {
		log.Fatalf(err.Error())
	}
	server := NewServer()
	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	slog.Info(fmt.Sprintf("HEALTH"))
	// 	w.WriteHeader(http.StatusOK)
	// })

	// mux.HandleFunc("GET /{queue}", func(w http.ResponseWriter, r *http.Request) {
	// })

	// mux.HandleFunc("POST /{queue}/purge", func(w http.ResponseWriter, r *http.Request) {
	// 	queueName := r.PathValue("queue")
	// 	slog.Info(fmt.Sprintf("DEL %v", queueName))
	// 	// queue, _ := queues[queueName]
	// 	delete(server.queues, queueName)
	// 	w.WriteHeader(http.StatusOK)
	// })

	// mux.HandleFunc("POST /{queue}", func(w http.ResponseWriter, r *http.Request) {
	// 	queueName := r.PathValue("queue")
	// 	data, _ := io.ReadAll(r.Body)
	// 	request := &gen.EnqueueRequest{
	// 		Queue: queueName,
	// 		Data:  data,
	// 	}
	// 	_, err := server.Enqueue(r.Context(), request)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.WriteHeader(http.StatusOK)
	// })

	// mux.HandleFunc("DELETE /{queue}", func(w http.ResponseWriter, r *http.Request) {
	// 	queueName := r.PathValue("queue")
	// 	queue := server.queues[queueName]
	// 	slog.Info(fmt.Sprintf("DEQ %v", queueName))
	// 	for {
	// 		if len(queue.data) > 0 {
	// 			break
	// 		}
	// 		select {
	// 		case <-time.After(time.Second / 10):
	// 			queue = server.queues[queueName]
	// 		case <-r.Context().Done():
	// 			return
	// 		}
	// 	}
	// 	data := queue.data[0]
	// 	queue.data = queue.data[1:]
	// 	w.Write(data)
	// })

	// go http.ListenAndServe(":42080", mux)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gen.RegisterQueueServiceServer(grpcServer, server)
	grpcServer.Serve(lis)

}
