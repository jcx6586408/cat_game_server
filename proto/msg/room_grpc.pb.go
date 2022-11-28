// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: pb/room.proto

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

// RoomClient is the client API for Room service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomClient interface {
	Create(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (Room_CreateClient, error)
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Over(ctx context.Context, in *OverRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
}

type roomClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomClient(cc grpc.ClientConnInterface) RoomClient {
	return &roomClient{cc}
}

func (c *roomClient) Create(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (Room_CreateClient, error) {
	stream, err := c.cc.NewStream(ctx, &Room_ServiceDesc.Streams[0], "/msg.Room/Create", opts...)
	if err != nil {
		return nil, err
	}
	x := &roomCreateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Room_CreateClient interface {
	Recv() (*CreateRoomReply, error)
	grpc.ClientStream
}

type roomCreateClient struct {
	grpc.ClientStream
}

func (x *roomCreateClient) Recv() (*CreateRoomReply, error) {
	m := new(CreateRoomReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *roomClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) Over(ctx context.Context, in *OverRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/Over", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServer is the server API for Room service.
// All implementations must embed UnimplementedRoomServer
// for forward compatibility
type RoomServer interface {
	Create(*CreateRoomRequest, Room_CreateServer) error
	Add(context.Context, *AddRequest) (*RoomChangeState, error)
	Leave(context.Context, *LeaveRequest) (*RoomChangeState, error)
	Over(context.Context, *OverRequest) (*RoomChangeState, error)
	mustEmbedUnimplementedRoomServer()
}

// UnimplementedRoomServer must be embedded to have forward compatible implementations.
type UnimplementedRoomServer struct {
}

func (UnimplementedRoomServer) Create(*CreateRoomRequest, Room_CreateServer) error {
	return status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRoomServer) Add(context.Context, *AddRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedRoomServer) Leave(context.Context, *LeaveRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedRoomServer) Over(context.Context, *OverRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Over not implemented")
}
func (UnimplementedRoomServer) mustEmbedUnimplementedRoomServer() {}

// UnsafeRoomServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServer will
// result in compilation errors.
type UnsafeRoomServer interface {
	mustEmbedUnimplementedRoomServer()
}

func RegisterRoomServer(s grpc.ServiceRegistrar, srv RoomServer) {
	s.RegisterService(&Room_ServiceDesc, srv)
}

func _Room_Create_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CreateRoomRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RoomServer).Create(m, &roomCreateServer{stream})
}

type Room_CreateServer interface {
	Send(*CreateRoomReply) error
	grpc.ServerStream
}

type roomCreateServer struct {
	grpc.ServerStream
}

func (x *roomCreateServer) Send(m *CreateRoomReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Room_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_Over_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Over(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/Over",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Over(ctx, req.(*OverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Room_ServiceDesc is the grpc.ServiceDesc for Room service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Room_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "msg.Room",
	HandlerType: (*RoomServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Room_Add_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _Room_Leave_Handler,
		},
		{
			MethodName: "Over",
			Handler:    _Room_Over_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Create",
			Handler:       _Room_Create_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/room.proto",
}