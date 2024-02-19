// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1alpha1

import (
	context "context"
)

import (
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// InterCpPingServiceClient is the client API for InterCpPingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InterCpPingServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type interCpPingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInterCpPingServiceClient(cc grpc.ClientConnInterface) InterCpPingServiceClient {
	return &interCpPingServiceClient{cc}
}

func (c *interCpPingServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/dubbo.system.v1alpha1.InterCpPingService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InterCpPingServiceServer is the server API for InterCpPingService service.
// All implementations must embed UnimplementedInterCpPingServiceServer
// for forward compatibility
type InterCpPingServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedInterCpPingServiceServer()
}

// UnimplementedInterCpPingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInterCpPingServiceServer struct {
}

func (UnimplementedInterCpPingServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedInterCpPingServiceServer) mustEmbedUnimplementedInterCpPingServiceServer() {}

// UnsafeInterCpPingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InterCpPingServiceServer will
// result in compilation errors.
type UnsafeInterCpPingServiceServer interface {
	mustEmbedUnimplementedInterCpPingServiceServer()
}

func RegisterInterCpPingServiceServer(s grpc.ServiceRegistrar, srv InterCpPingServiceServer) {
	s.RegisterService(&InterCpPingService_ServiceDesc, srv)
}

func _InterCpPingService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterCpPingServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dubbo.system.v1alpha1.InterCpPingService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterCpPingServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InterCpPingService_ServiceDesc is the grpc.ServiceDesc for InterCpPingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InterCpPingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dubbo.system.v1alpha1.InterCpPingService",
	HandlerType: (*InterCpPingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _InterCpPingService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/system/v1alpha1/inter_cp_ping.proto",
}
