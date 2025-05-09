// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: replica/replica.proto

package replica

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
	Replica_ShareID_FullMethodName         = "/Replica/shareID"
	Replica_NotifyNewLeader_FullMethodName = "/Replica/notifyNewLeader"
	Replica_BecomeLeader_FullMethodName    = "/Replica/becomeLeader"
	Replica_Heartbeat_FullMethodName       = "/Replica/heartbeat"
	Replica_GetStatus_FullMethodName       = "/Replica/getStatus"
)

// ReplicaClient is the client API for Replica service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicaClient interface {
	ShareID(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*IDResponse, error)
	NotifyNewLeader(ctx context.Context, in *LeaderNotification, opts ...grpc.CallOption) (*IDResponse, error)
	BecomeLeader(ctx context.Context, in *LeaderTransfer, opts ...grpc.CallOption) (*IDResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type replicaClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicaClient(cc grpc.ClientConnInterface) ReplicaClient {
	return &replicaClient{cc}
}

func (c *replicaClient) ShareID(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*IDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, Replica_ShareID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) NotifyNewLeader(ctx context.Context, in *LeaderNotification, opts ...grpc.CallOption) (*IDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, Replica_NotifyNewLeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) BecomeLeader(ctx context.Context, in *LeaderTransfer, opts ...grpc.CallOption) (*IDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, Replica_BecomeLeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, Replica_Heartbeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Replica_GetStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicaServer is the server API for Replica service.
// All implementations must embed UnimplementedReplicaServer
// for forward compatibility.
type ReplicaServer interface {
	ShareID(context.Context, *IDRequest) (*IDResponse, error)
	NotifyNewLeader(context.Context, *LeaderNotification) (*IDResponse, error)
	BecomeLeader(context.Context, *LeaderTransfer) (*IDResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
	mustEmbedUnimplementedReplicaServer()
}

// UnimplementedReplicaServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReplicaServer struct{}

func (UnimplementedReplicaServer) ShareID(context.Context, *IDRequest) (*IDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareID not implemented")
}
func (UnimplementedReplicaServer) NotifyNewLeader(context.Context, *LeaderNotification) (*IDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyNewLeader not implemented")
}
func (UnimplementedReplicaServer) BecomeLeader(context.Context, *LeaderTransfer) (*IDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BecomeLeader not implemented")
}
func (UnimplementedReplicaServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedReplicaServer) GetStatus(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedReplicaServer) mustEmbedUnimplementedReplicaServer() {}
func (UnimplementedReplicaServer) testEmbeddedByValue()                 {}

// UnsafeReplicaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicaServer will
// result in compilation errors.
type UnsafeReplicaServer interface {
	mustEmbedUnimplementedReplicaServer()
}

func RegisterReplicaServer(s grpc.ServiceRegistrar, srv ReplicaServer) {
	// If the following call pancis, it indicates UnimplementedReplicaServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Replica_ServiceDesc, srv)
}

func _Replica_ShareID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).ShareID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Replica_ShareID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).ShareID(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_NotifyNewLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderNotification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).NotifyNewLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Replica_NotifyNewLeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).NotifyNewLeader(ctx, req.(*LeaderNotification))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_BecomeLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderTransfer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).BecomeLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Replica_BecomeLeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).BecomeLeader(ctx, req.(*LeaderTransfer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Replica_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Replica_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).GetStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Replica_ServiceDesc is the grpc.ServiceDesc for Replica service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Replica_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Replica",
	HandlerType: (*ReplicaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "shareID",
			Handler:    _Replica_ShareID_Handler,
		},
		{
			MethodName: "notifyNewLeader",
			Handler:    _Replica_NotifyNewLeader_Handler,
		},
		{
			MethodName: "becomeLeader",
			Handler:    _Replica_BecomeLeader_Handler,
		},
		{
			MethodName: "heartbeat",
			Handler:    _Replica_Heartbeat_Handler,
		},
		{
			MethodName: "getStatus",
			Handler:    _Replica_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "replica/replica.proto",
}
