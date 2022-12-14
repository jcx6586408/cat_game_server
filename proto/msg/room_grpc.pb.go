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
	Connect(ctx context.Context, in *RoomServerConnectRequest, opts ...grpc.CallOption) (Room_ConnectClient, error)
	Create(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Over(ctx context.Context, in *OverRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	AnswerQuestion(ctx context.Context, in *Answer, opts ...grpc.CallOption) (*RoomChangeState, error)
	MatchRoom(ctx context.Context, in *MatchRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	MatchMember(ctx context.Context, in *MatchMemberRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	MatchRoomCancel(ctx context.Context, in *MatchRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	MatchMemberCancel(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
	Offline(ctx context.Context, in *OfflineRequest, opts ...grpc.CallOption) (*RoomChangeState, error)
}

type roomClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomClient(cc grpc.ClientConnInterface) RoomClient {
	return &roomClient{cc}
}

func (c *roomClient) Connect(ctx context.Context, in *RoomServerConnectRequest, opts ...grpc.CallOption) (Room_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Room_ServiceDesc.Streams[0], "/msg.Room/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &roomConnectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Room_ConnectClient interface {
	Recv() (*CreateRoomReply, error)
	grpc.ClientStream
}

type roomConnectClient struct {
	grpc.ClientStream
}

func (x *roomConnectClient) Recv() (*CreateRoomReply, error) {
	m := new(CreateRoomReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *roomClient) Create(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *roomClient) AnswerQuestion(ctx context.Context, in *Answer, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/AnswerQuestion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) MatchRoom(ctx context.Context, in *MatchRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/MatchRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) MatchMember(ctx context.Context, in *MatchMemberRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/MatchMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) MatchRoomCancel(ctx context.Context, in *MatchRoomRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/MatchRoomCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) MatchMemberCancel(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/MatchMemberCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) Offline(ctx context.Context, in *OfflineRequest, opts ...grpc.CallOption) (*RoomChangeState, error) {
	out := new(RoomChangeState)
	err := c.cc.Invoke(ctx, "/msg.Room/Offline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServer is the server API for Room service.
// All implementations must embed UnimplementedRoomServer
// for forward compatibility
type RoomServer interface {
	Connect(*RoomServerConnectRequest, Room_ConnectServer) error
	Create(context.Context, *CreateRoomRequest) (*RoomChangeState, error)
	Add(context.Context, *AddRequest) (*RoomChangeState, error)
	Leave(context.Context, *LeaveRequest) (*RoomChangeState, error)
	Over(context.Context, *OverRequest) (*RoomChangeState, error)
	AnswerQuestion(context.Context, *Answer) (*RoomChangeState, error)
	MatchRoom(context.Context, *MatchRoomRequest) (*RoomChangeState, error)
	MatchMember(context.Context, *MatchMemberRequest) (*RoomChangeState, error)
	MatchRoomCancel(context.Context, *MatchRoomRequest) (*RoomChangeState, error)
	MatchMemberCancel(context.Context, *LeaveRequest) (*RoomChangeState, error)
	Offline(context.Context, *OfflineRequest) (*RoomChangeState, error)
	mustEmbedUnimplementedRoomServer()
}

// UnimplementedRoomServer must be embedded to have forward compatible implementations.
type UnimplementedRoomServer struct {
}

func (UnimplementedRoomServer) Connect(*RoomServerConnectRequest, Room_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedRoomServer) Create(context.Context, *CreateRoomRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
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
func (UnimplementedRoomServer) AnswerQuestion(context.Context, *Answer) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnswerQuestion not implemented")
}
func (UnimplementedRoomServer) MatchRoom(context.Context, *MatchRoomRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchRoom not implemented")
}
func (UnimplementedRoomServer) MatchMember(context.Context, *MatchMemberRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchMember not implemented")
}
func (UnimplementedRoomServer) MatchRoomCancel(context.Context, *MatchRoomRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchRoomCancel not implemented")
}
func (UnimplementedRoomServer) MatchMemberCancel(context.Context, *LeaveRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchMemberCancel not implemented")
}
func (UnimplementedRoomServer) Offline(context.Context, *OfflineRequest) (*RoomChangeState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Offline not implemented")
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

func _Room_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RoomServerConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RoomServer).Connect(m, &roomConnectServer{stream})
}

type Room_ConnectServer interface {
	Send(*CreateRoomReply) error
	grpc.ServerStream
}

type roomConnectServer struct {
	grpc.ServerStream
}

func (x *roomConnectServer) Send(m *CreateRoomReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Room_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Create(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _Room_AnswerQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Answer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).AnswerQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/AnswerQuestion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).AnswerQuestion(ctx, req.(*Answer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_MatchRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).MatchRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/MatchRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).MatchRoom(ctx, req.(*MatchRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_MatchMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).MatchMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/MatchMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).MatchMember(ctx, req.(*MatchMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_MatchRoomCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).MatchRoomCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/MatchRoomCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).MatchRoomCancel(ctx, req.(*MatchRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_MatchMemberCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).MatchMemberCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/MatchMemberCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).MatchMemberCancel(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_Offline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OfflineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).Offline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/msg.Room/Offline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).Offline(ctx, req.(*OfflineRequest))
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
			MethodName: "Create",
			Handler:    _Room_Create_Handler,
		},
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
		{
			MethodName: "AnswerQuestion",
			Handler:    _Room_AnswerQuestion_Handler,
		},
		{
			MethodName: "MatchRoom",
			Handler:    _Room_MatchRoom_Handler,
		},
		{
			MethodName: "MatchMember",
			Handler:    _Room_MatchMember_Handler,
		},
		{
			MethodName: "MatchRoomCancel",
			Handler:    _Room_MatchRoomCancel_Handler,
		},
		{
			MethodName: "MatchMemberCancel",
			Handler:    _Room_MatchMemberCancel_Handler,
		},
		{
			MethodName: "Offline",
			Handler:    _Room_Offline_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _Room_Connect_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/room.proto",
}
