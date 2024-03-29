// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.2
// source: calibrate.proto

package pb

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

const (
	Calibrate_CalibrateOnePort_FullMethodName = "/pb.Calibrate/CalibrateOnePort"
	Calibrate_CalibrateTwoPort_FullMethodName = "/pb.Calibrate/CalibrateTwoPort"
)

// CalibrateClient is the client API for Calibrate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalibrateClient interface {
	CalibrateOnePort(ctx context.Context, in *CalibrateOnePortRequest, opts ...grpc.CallOption) (*CalibrateOnePortResponse, error)
	CalibrateTwoPort(ctx context.Context, in *CalibrateTwoPortRequest, opts ...grpc.CallOption) (*CalibrateTwoPortResponse, error)
}

type calibrateClient struct {
	cc grpc.ClientConnInterface
}

func NewCalibrateClient(cc grpc.ClientConnInterface) CalibrateClient {
	return &calibrateClient{cc}
}

func (c *calibrateClient) CalibrateOnePort(ctx context.Context, in *CalibrateOnePortRequest, opts ...grpc.CallOption) (*CalibrateOnePortResponse, error) {
	out := new(CalibrateOnePortResponse)
	err := c.cc.Invoke(ctx, Calibrate_CalibrateOnePort_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calibrateClient) CalibrateTwoPort(ctx context.Context, in *CalibrateTwoPortRequest, opts ...grpc.CallOption) (*CalibrateTwoPortResponse, error) {
	out := new(CalibrateTwoPortResponse)
	err := c.cc.Invoke(ctx, Calibrate_CalibrateTwoPort_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalibrateServer is the server API for Calibrate service.
// All implementations must embed UnimplementedCalibrateServer
// for forward compatibility
type CalibrateServer interface {
	CalibrateOnePort(context.Context, *CalibrateOnePortRequest) (*CalibrateOnePortResponse, error)
	CalibrateTwoPort(context.Context, *CalibrateTwoPortRequest) (*CalibrateTwoPortResponse, error)
	mustEmbedUnimplementedCalibrateServer()
}

// UnimplementedCalibrateServer must be embedded to have forward compatible implementations.
type UnimplementedCalibrateServer struct {
}

func (UnimplementedCalibrateServer) CalibrateOnePort(context.Context, *CalibrateOnePortRequest) (*CalibrateOnePortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalibrateOnePort not implemented")
}
func (UnimplementedCalibrateServer) CalibrateTwoPort(context.Context, *CalibrateTwoPortRequest) (*CalibrateTwoPortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalibrateTwoPort not implemented")
}
func (UnimplementedCalibrateServer) mustEmbedUnimplementedCalibrateServer() {}

// UnsafeCalibrateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalibrateServer will
// result in compilation errors.
type UnsafeCalibrateServer interface {
	mustEmbedUnimplementedCalibrateServer()
}

func RegisterCalibrateServer(s grpc.ServiceRegistrar, srv CalibrateServer) {
	s.RegisterService(&Calibrate_ServiceDesc, srv)
}

func _Calibrate_CalibrateOnePort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalibrateOnePortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalibrateServer).CalibrateOnePort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calibrate_CalibrateOnePort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalibrateServer).CalibrateOnePort(ctx, req.(*CalibrateOnePortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calibrate_CalibrateTwoPort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalibrateTwoPortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalibrateServer).CalibrateTwoPort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calibrate_CalibrateTwoPort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalibrateServer).CalibrateTwoPort(ctx, req.(*CalibrateTwoPortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calibrate_ServiceDesc is the grpc.ServiceDesc for Calibrate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calibrate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Calibrate",
	HandlerType: (*CalibrateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CalibrateOnePort",
			Handler:    _Calibrate_CalibrateOnePort_Handler,
		},
		{
			MethodName: "CalibrateTwoPort",
			Handler:    _Calibrate_CalibrateTwoPort_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "calibrate.proto",
}
