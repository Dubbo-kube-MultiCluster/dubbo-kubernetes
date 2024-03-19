// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1alpha1

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

// MetadataServiceClient is the client API for MetadataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetadataServiceClient interface {
	MetadataRegister(ctx context.Context, in *MetaDataRegisterRequest, opts ...grpc.CallOption) (*MetaDataRegisterResponse, error)
	MetadataSync(ctx context.Context, opts ...grpc.CallOption) (MetadataService_MetadataSyncClient, error)
}

type metadataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetadataServiceClient(cc grpc.ClientConnInterface) MetadataServiceClient {
	return &metadataServiceClient{cc}
}

func (c *metadataServiceClient) MetadataRegister(ctx context.Context, in *MetaDataRegisterRequest, opts ...grpc.CallOption) (*MetaDataRegisterResponse, error) {
	out := new(MetaDataRegisterResponse)
	err := c.cc.Invoke(ctx, "/dubbo.mesh.v1alpha1.MetadataService/MetadataRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) MetadataSync(ctx context.Context, opts ...grpc.CallOption) (MetadataService_MetadataSyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &MetadataService_ServiceDesc.Streams[0], "/dubbo.mesh.v1alpha1.MetadataService/MetadataSync", opts...)
	if err != nil {
		return nil, err
	}
	x := &metadataServiceMetadataSyncClient{stream}
	return x, nil
}

type MetadataService_MetadataSyncClient interface {
	Send(*MetadataSyncRequest) error
	Recv() (*MetadataSyncResponse, error)
	grpc.ClientStream
}

type metadataServiceMetadataSyncClient struct {
	grpc.ClientStream
}

func (x *metadataServiceMetadataSyncClient) Send(m *MetadataSyncRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *metadataServiceMetadataSyncClient) Recv() (*MetadataSyncResponse, error) {
	m := new(MetadataSyncResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetadataServiceServer is the server API for MetadataService service.
// All implementations must embed UnimplementedMetadataServiceServer
// for forward compatibility
type MetadataServiceServer interface {
	MetadataRegister(context.Context, *MetaDataRegisterRequest) (*MetaDataRegisterResponse, error)
	MetadataSync(MetadataService_MetadataSyncServer) error
	mustEmbedUnimplementedMetadataServiceServer()
}

// UnimplementedMetadataServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetadataServiceServer struct {
}

func (UnimplementedMetadataServiceServer) MetadataRegister(context.Context, *MetaDataRegisterRequest) (*MetaDataRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MetadataRegister not implemented")
}
func (UnimplementedMetadataServiceServer) MetadataSync(MetadataService_MetadataSyncServer) error {
	return status.Errorf(codes.Unimplemented, "method MetadataSync not implemented")
}
func (UnimplementedMetadataServiceServer) mustEmbedUnimplementedMetadataServiceServer() {}

// UnsafeMetadataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetadataServiceServer will
// result in compilation errors.
type UnsafeMetadataServiceServer interface {
	mustEmbedUnimplementedMetadataServiceServer()
}

func RegisterMetadataServiceServer(s grpc.ServiceRegistrar, srv MetadataServiceServer) {
	s.RegisterService(&MetadataService_ServiceDesc, srv)
}

func _MetadataService_MetadataRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetaDataRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).MetadataRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dubbo.mesh.v1alpha1.MetadataService/MetadataRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).MetadataRegister(ctx, req.(*MetaDataRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_MetadataSync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MetadataServiceServer).MetadataSync(&metadataServiceMetadataSyncServer{stream})
}

type MetadataService_MetadataSyncServer interface {
	Send(*MetadataSyncResponse) error
	Recv() (*MetadataSyncRequest, error)
	grpc.ServerStream
}

type metadataServiceMetadataSyncServer struct {
	grpc.ServerStream
}

func (x *metadataServiceMetadataSyncServer) Send(m *MetadataSyncResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *metadataServiceMetadataSyncServer) Recv() (*MetadataSyncRequest, error) {
	m := new(MetadataSyncRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetadataService_ServiceDesc is the grpc.ServiceDesc for MetadataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetadataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dubbo.mesh.v1alpha1.MetadataService",
	HandlerType: (*MetadataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MetadataRegister",
			Handler:    _MetadataService_MetadataRegister_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MetadataSync",
			Handler:       _MetadataService_MetadataSync_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/mesh/v1alpha1/metadata.proto",
}
