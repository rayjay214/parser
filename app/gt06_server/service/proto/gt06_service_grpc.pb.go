// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: gt06_service.proto

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

// Gt06ServiceClient is the client API for Gt06Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Gt06ServiceClient interface {
	SendCmd(ctx context.Context, in *SendGt06CmdRequest, opts ...grpc.CallOption) (*Gt06CommonReply, error)
}

type gt06ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGt06ServiceClient(cc grpc.ClientConnInterface) Gt06ServiceClient {
	return &gt06ServiceClient{cc}
}

func (c *gt06ServiceClient) SendCmd(ctx context.Context, in *SendGt06CmdRequest, opts ...grpc.CallOption) (*Gt06CommonReply, error) {
	out := new(Gt06CommonReply)
	err := c.cc.Invoke(ctx, "/service.Gt06Service/SendCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Gt06ServiceServer is the server API for Gt06Service service.
// All implementations must embed UnimplementedGt06ServiceServer
// for forward compatibility
type Gt06ServiceServer interface {
	SendCmd(context.Context, *SendGt06CmdRequest) (*Gt06CommonReply, error)
	mustEmbedUnimplementedGt06ServiceServer()
}

// UnimplementedGt06ServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGt06ServiceServer struct {
}

func (UnimplementedGt06ServiceServer) SendCmd(context.Context, *SendGt06CmdRequest) (*Gt06CommonReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCmd not implemented")
}
func (UnimplementedGt06ServiceServer) mustEmbedUnimplementedGt06ServiceServer() {}

// UnsafeGt06ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Gt06ServiceServer will
// result in compilation errors.
type UnsafeGt06ServiceServer interface {
	mustEmbedUnimplementedGt06ServiceServer()
}

func RegisterGt06ServiceServer(s grpc.ServiceRegistrar, srv Gt06ServiceServer) {
	s.RegisterService(&Gt06Service_ServiceDesc, srv)
}

func _Gt06Service_SendCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendGt06CmdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Gt06ServiceServer).SendCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Gt06Service/SendCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Gt06ServiceServer).SendCmd(ctx, req.(*SendGt06CmdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gt06Service_ServiceDesc is the grpc.ServiceDesc for Gt06Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gt06Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Gt06Service",
	HandlerType: (*Gt06ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCmd",
			Handler:    _Gt06Service_SendCmd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gt06_service.proto",
}
