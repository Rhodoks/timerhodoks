// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.0--rc2
// source: coordinatorpb/coordinator.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CoordinatorClient is the client API for Coordinator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoordinatorClient interface {
	GetCoordinatingInfo(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*CoordinatingInfoReply, error)
	UpdateAllocation(ctx context.Context, in *AllocationRequest, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type coordinatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCoordinatorClient(cc grpc.ClientConnInterface) CoordinatorClient {
	return &coordinatorClient{cc}
}

func (c *coordinatorClient) GetCoordinatingInfo(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*CoordinatingInfoReply, error) {
	out := new(CoordinatingInfoReply)
	err := c.cc.Invoke(ctx, "/proto.Coordinator/GetCoordinatingInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorClient) UpdateAllocation(ctx context.Context, in *AllocationRequest, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/proto.Coordinator/UpdateAllocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoordinatorServer is the server API for Coordinator service.
// All implementations must embed UnimplementedCoordinatorServer
// for forward compatibility
type CoordinatorServer interface {
	GetCoordinatingInfo(context.Context, *EmptyMessage) (*CoordinatingInfoReply, error)
	UpdateAllocation(context.Context, *AllocationRequest) (*EmptyMessage, error)
	mustEmbedUnimplementedCoordinatorServer()
}

// UnimplementedCoordinatorServer must be embedded to have forward compatible implementations.
type UnimplementedCoordinatorServer struct {
}

func (UnimplementedCoordinatorServer) GetCoordinatingInfo(context.Context, *EmptyMessage) (*CoordinatingInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoordinatingInfo not implemented")
}
func (UnimplementedCoordinatorServer) UpdateAllocation(context.Context, *AllocationRequest) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAllocation not implemented")
}
func (UnimplementedCoordinatorServer) mustEmbedUnimplementedCoordinatorServer() {}

// UnsafeCoordinatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoordinatorServer will
// result in compilation errors.
type UnsafeCoordinatorServer interface {
	mustEmbedUnimplementedCoordinatorServer()
}

func RegisterCoordinatorServer(s grpc.ServiceRegistrar, srv CoordinatorServer) {
	s.RegisterService(&Coordinator_ServiceDesc, srv)
}

func _Coordinator_GetCoordinatingInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServer).GetCoordinatingInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Coordinator/GetCoordinatingInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServer).GetCoordinatingInfo(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coordinator_UpdateAllocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServer).UpdateAllocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Coordinator/UpdateAllocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServer).UpdateAllocation(ctx, req.(*AllocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Coordinator_ServiceDesc is the grpc.ServiceDesc for Coordinator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coordinator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Coordinator",
	HandlerType: (*CoordinatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCoordinatingInfo",
			Handler:    _Coordinator_GetCoordinatingInfo_Handler,
		},
		{
			MethodName: "UpdateAllocation",
			Handler:    _Coordinator_UpdateAllocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coordinatorpb/coordinator.proto",
}
