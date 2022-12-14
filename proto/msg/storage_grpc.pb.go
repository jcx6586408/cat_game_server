// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: pb/storage.proto

package msg

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

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	LoginStorage(ctx context.Context, in *LoginStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error)
	OffLineStorage(ctx context.Context, in *OffLineStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error)
	UpdateStorage(ctx context.Context, in *UpdateStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error)
	PullStorage(ctx context.Context, in *PullStorageRequest, opts ...grpc.CallOption) (*PullStorageReply, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) LoginStorage(ctx context.Context, in *LoginStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error) {
	out := new(SuccessStateReply)
	err := c.cc.Invoke(ctx, "/msg.Storage/LoginStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) OffLineStorage(ctx context.Context, in *OffLineStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error) {
	out := new(SuccessStateReply)
	err := c.cc.Invoke(ctx, "/msg.Storage/OffLineStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) UpdateStorage(ctx context.Context, in *UpdateStorageRequest, opts ...grpc.CallOption) (*SuccessStateReply, error) {
	out := new(SuccessStateReply)
	err := c.cc.Invoke(ctx, "/msg.Storage/UpdateStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) PullStorage(ctx context.Context, in *PullStorageRequest, opts ...grpc.CallOption) (*PullStorageReply, error) {
	out := new(PullStorageReply)
	err := c.cc.Invoke(ctx, "/msg.Storage/PullStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	LoginStorage(context.Context, *LoginStorageRequest) (*SuccessStateReply, error)
	OffLineStorage(context.Context, *OffLineStorageRequest) (*SuccessStateReply, error)
	UpdateStorage(context.Context, *UpdateStorageRequest) (*SuccessStateReply, error)
	PullStorage(context.Context, *PullStorageRequest) (*PullStorageReply, error)
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) LoginStorage(context.Context, *LoginStorageRequest) (*SuccessStateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginStorage not implemented")
}
func (UnimplementedStorageServer) OffLineStorage(context.Context, *OffLineStorageRequest) (*SuccessStateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OffLineStorage not implemented")
}
func (UnimplementedStorageServer) UpdateStorage(context.Context, *UpdateStorageRequest) (*SuccessStateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStorage not implemented")
}
func (UnimplementedStorageServer) PullStorage(context.Context, *PullStorageRequest) (*PullStorageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullStorage not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_LoginStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).LoginStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Storage/LoginStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).LoginStorage(ctx, req.(*LoginStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_OffLineStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OffLineStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).OffLineStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Storage/OffLineStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).OffLineStorage(ctx, req.(*OffLineStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_UpdateStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).UpdateStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Storage/UpdateStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).UpdateStorage(ctx, req.(*UpdateStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_PullStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).PullStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Storage/PullStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).PullStorage(ctx, req.(*PullStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "msg.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginStorage",
			Handler:    _Storage_LoginStorage_Handler,
		},
		{
			MethodName: "OffLineStorage",
			Handler:    _Storage_OffLineStorage_Handler,
		},
		{
			MethodName: "UpdateStorage",
			Handler:    _Storage_UpdateStorage_Handler,
		},
		{
			MethodName: "PullStorage",
			Handler:    _Storage_PullStorage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/storage.proto",
}
