package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/lmittmann/tint"
	"github.com/thekashifmalik/litemq/gen"
	"github.com/thekashifmalik/litemq/internal/queues"
	"google.golang.org/grpc"
)

const defaultPort = 42090

type Server struct {
	gen.UnimplementedLiteMQServer

	queues map[string]*queues.Queue
	lock   sync.Mutex
}

func NewServer() *Server {
	logLevel := slog.LevelInfo
	level, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		logLevel.UnmarshalText([]byte(level))
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.RFC3339,
		}),
	))
	slog.Info("LiteMQ")
	slog.Debug("debug logging enabled")
	return &Server{
		queues: map[string]*queues.Queue{},
	}
}

func (s *Server) Serve() {
	port := defaultPort
	_port, ok := os.LookupEnv("PORT")
	if ok {
		var err error
		port, err = strconv.Atoi(_port)
		if err != nil {
			slog.Warn(fmt.Sprintf("Invalid PORT value: %v", _port))
			port = defaultPort
		}
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gen.RegisterLiteMQServer(grpcServer, s)
	slog.Info(fmt.Sprintf("listening on %v", addr))
	grpcServer.Serve(lis)
}

func (s *Server) Enqueue(ctx context.Context, request *gen.EnqueueRequest) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("ENQUEUE %v '%v'", request.Queue, string(request.Data)))
	queue := s.getOrCreateQueue(request.Queue)
	queue.LockAndEnqueue(request.Data)
	length := queue.Length()
	return &gen.QueueLength{Count: int64(length)}, nil
}

func (s *Server) Dequeue(ctx context.Context, request *gen.QueueID) (*gen.DequeueResponse, error) {
	slog.Info(fmt.Sprintf("DEQUEUE %v", request.Queue))
	queue := s.getOrCreateQueue(request.Queue)
	data, err := queue.LockAndDequeue(ctx)
	if err != nil {
		slog.Warn("< disconnected")
		return nil, err
	}
	slog.Info(fmt.Sprintf("< %v bytes", int64(len(data))))
	return &gen.DequeueResponse{Data: data}, nil
}

func (s *Server) getOrCreateQueue(name string) *queues.Queue {
	queue, ok := s.queues[name]
	if !ok {
		queue = s.lockAndCreateQueue(name)
	}
	return queue
}

func (s *Server) lockAndCreateQueue(name string) *queues.Queue {
	s.lock.Lock()
	defer s.lock.Unlock()
	// Since another request could have created the queue while we were acquiring the lock, we need to check whether the
	// queue exists, even if we checked just before calling this method.
	queue, ok := s.queues[name]
	if !ok {
		queue = queues.NewQueue()
		s.queues[name] = queue
	}
	return queue
}

func (s *Server) Purge(ctx context.Context, request *gen.QueueID) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("PURGE %v", request.Queue))
	s.lock.Lock()
	defer s.lock.Unlock()
	length := 0
	queue, ok := s.queues[request.Queue]
	if ok {
		delete(s.queues, request.Queue)
		length = queue.Length()
		// TODO: Need to delete any active dequeue channels for this purged queue.
	}
	return &gen.QueueLength{Count: int64(length)}, nil
}

func (s *Server) Length(ctx context.Context, request *gen.QueueID) (*gen.QueueLength, error) {
	slog.Info(fmt.Sprintf("LENGTH %v", request.Queue))
	length := 0
	queue, ok := s.queues[request.Queue]
	if ok {
		length = queue.Length()
	}
	return &gen.QueueLength{Count: int64(length)}, nil
}

func (s *Server) Health(ctx context.Context, request *gen.Nothing) (*gen.Nothing, error) {
	slog.Info(fmt.Sprintf("HEALTH"))
	return &gen.Nothing{}, nil
}
