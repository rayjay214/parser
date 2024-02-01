// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: device_service.proto

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

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	SendCmd(ctx context.Context, in *SendCmdRequest, opts ...grpc.CallOption) (*SendCmdReply, error)
	OpenShortRecord(ctx context.Context, in *OpenShortRecordRequest, opts ...grpc.CallOption) (*CommonReply, error)
	VorRecordSwitch(ctx context.Context, in *VorRecordSwitchRequest, opts ...grpc.CallOption) (*CommonReply, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) SendCmd(ctx context.Context, in *SendCmdRequest, opts ...grpc.CallOption) (*SendCmdReply, error) {
	out := new(SendCmdReply)
	err := c.cc.Invoke(ctx, "/service.DeviceService/SendCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) OpenShortRecord(ctx context.Context, in *OpenShortRecordRequest, opts ...grpc.CallOption) (*CommonReply, error) {
	out := new(CommonReply)
	err := c.cc.Invoke(ctx, "/service.DeviceService/OpenShortRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) VorRecordSwitch(ctx context.Context, in *VorRecordSwitchRequest, opts ...grpc.CallOption) (*CommonReply, error) {
	out := new(CommonReply)
	err := c.cc.Invoke(ctx, "/service.DeviceService/VorRecordSwitch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility
type DeviceServiceServer interface {
	SendCmd(context.Context, *SendCmdRequest) (*SendCmdReply, error)
	OpenShortRecord(context.Context, *OpenShortRecordRequest) (*CommonReply, error)
	VorRecordSwitch(context.Context, *VorRecordSwitchRequest) (*CommonReply, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServiceServer struct {
}

func (UnimplementedDeviceServiceServer) SendCmd(context.Context, *SendCmdRequest) (*SendCmdReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCmd not implemented")
}
func (UnimplementedDeviceServiceServer) OpenShortRecord(context.Context, *OpenShortRecordRequest) (*CommonReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenShortRecord not implemented")
}
func (UnimplementedDeviceServiceServer) VorRecordSwitch(context.Context, *VorRecordSwitchRequest) (*CommonReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VorRecordSwitch not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_SendCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCmdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).SendCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.DeviceService/SendCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).SendCmd(ctx, req.(*SendCmdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_OpenShortRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenShortRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).OpenShortRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.DeviceService/OpenShortRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).OpenShortRecord(ctx, req.(*OpenShortRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_VorRecordSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VorRecordSwitchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).VorRecordSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.DeviceService/VorRecordSwitch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).VorRecordSwitch(ctx, req.(*VorRecordSwitchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCmd",
			Handler:    _DeviceService_SendCmd_Handler,
		},
		{
			MethodName: "OpenShortRecord",
			Handler:    _DeviceService_OpenShortRecord_Handler,
		},
		{
			MethodName: "VorRecordSwitch",
			Handler:    _DeviceService_VorRecordSwitch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "device_service.proto",
}
