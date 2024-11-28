// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: service.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	QueueService_Enqueue_FullMethodName = "/QueueService/Enqueue"
	QueueService_Dequeue_FullMethodName = "/QueueService/Dequeue"
	QueueService_Purge_FullMethodName   = "/QueueService/Purge"
	QueueService_Length_FullMethodName  = "/QueueService/Length"
	QueueService_Health_FullMethodName  = "/QueueService/Health"
)

// QueueServiceClient is the client API for QueueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueueServiceClient interface {
	Enqueue(ctx context.Context, in *EnqueueRequest, opts ...grpc.CallOption) (*Nothing, error)
	Dequeue(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*DequeueResponse, error)
	Purge(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*QueueLength, error)
	Length(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*QueueLength, error)
	Health(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

type queueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueServiceClient(cc grpc.ClientConnInterface) QueueServiceClient {
	return &queueServiceClient{cc}
}

func (c *queueServiceClient) Enqueue(ctx context.Context, in *EnqueueRequest, opts ...grpc.CallOption) (*Nothing, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nothing)
	err := c.cc.Invoke(ctx, QueueService_Enqueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) Dequeue(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*DequeueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DequeueResponse)
	err := c.cc.Invoke(ctx, QueueService_Dequeue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) Purge(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*QueueLength, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueueLength)
	err := c.cc.Invoke(ctx, QueueService_Purge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) Length(ctx context.Context, in *QueueID, opts ...grpc.CallOption) (*QueueLength, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueueLength)
	err := c.cc.Invoke(ctx, QueueService_Length_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) Health(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nothing)
	err := c.cc.Invoke(ctx, QueueService_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueServiceServer is the server API for QueueService service.
// All implementations must embed UnimplementedQueueServiceServer
// for forward compatibility.
type QueueServiceServer interface {
	Enqueue(context.Context, *EnqueueRequest) (*Nothing, error)
	Dequeue(context.Context, *QueueID) (*DequeueResponse, error)
	Purge(context.Context, *QueueID) (*QueueLength, error)
	Length(context.Context, *QueueID) (*QueueLength, error)
	Health(context.Context, *Nothing) (*Nothing, error)
	mustEmbedUnimplementedQueueServiceServer()
}

// UnimplementedQueueServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueueServiceServer struct{}

func (UnimplementedQueueServiceServer) Enqueue(context.Context, *EnqueueRequest) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enqueue not implemented")
}
func (UnimplementedQueueServiceServer) Dequeue(context.Context, *QueueID) (*DequeueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dequeue not implemented")
}
func (UnimplementedQueueServiceServer) Purge(context.Context, *QueueID) (*QueueLength, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Purge not implemented")
}
func (UnimplementedQueueServiceServer) Length(context.Context, *QueueID) (*QueueLength, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Length not implemented")
}
func (UnimplementedQueueServiceServer) Health(context.Context, *Nothing) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedQueueServiceServer) mustEmbedUnimplementedQueueServiceServer() {}
func (UnimplementedQueueServiceServer) testEmbeddedByValue()                      {}

// UnsafeQueueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueueServiceServer will
// result in compilation errors.
type UnsafeQueueServiceServer interface {
	mustEmbedUnimplementedQueueServiceServer()
}

func RegisterQueueServiceServer(s grpc.ServiceRegistrar, srv QueueServiceServer) {
	// If the following call pancis, it indicates UnimplementedQueueServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QueueService_ServiceDesc, srv)
}

func _QueueService_Enqueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnqueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).Enqueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueService_Enqueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).Enqueue(ctx, req.(*EnqueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_Dequeue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).Dequeue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueService_Dequeue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).Dequeue(ctx, req.(*QueueID))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_Purge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).Purge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueService_Purge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).Purge(ctx, req.(*QueueID))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_Length_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).Length(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueService_Length_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).Length(ctx, req.(*QueueID))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueService_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).Health(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

// QueueService_ServiceDesc is the grpc.ServiceDesc for QueueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QueueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "QueueService",
	HandlerType: (*QueueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enqueue",
			Handler:    _QueueService_Enqueue_Handler,
		},
		{
			MethodName: "Dequeue",
			Handler:    _QueueService_Dequeue_Handler,
		},
		{
			MethodName: "Purge",
			Handler:    _QueueService_Purge_Handler,
		},
		{
			MethodName: "Length",
			Handler:    _QueueService_Length_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _QueueService_Health_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
