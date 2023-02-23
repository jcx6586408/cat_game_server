// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: pb/center.proto

package msg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CenterPing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=Url,proto3" json:"Url,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"` // 人数统计
}

func (x *CenterPing) Reset() {
	*x = CenterPing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_center_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CenterPing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CenterPing) ProtoMessage() {}

func (x *CenterPing) ProtoReflect() protoreflect.Message {
	mi := &file_pb_center_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CenterPing.ProtoReflect.Descriptor instead.
func (*CenterPing) Descriptor() ([]byte, []int) {
	return file_pb_center_proto_rawDescGZIP(), []int{0}
}

func (x *CenterPing) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CenterPing) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type CenterPong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CenterPong) Reset() {
	*x = CenterPong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_center_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CenterPong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CenterPong) ProtoMessage() {}

func (x *CenterPong) ProtoReflect() protoreflect.Message {
	mi := &file_pb_center_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CenterPong.ProtoReflect.Descriptor instead.
func (*CenterPong) Descriptor() ([]byte, []int) {
	return file_pb_center_proto_rawDescGZIP(), []int{1}
}

type CenterConnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=Url,proto3" json:"Url,omitempty"` // 业务服的地址
}

func (x *CenterConnectRequest) Reset() {
	*x = CenterConnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_center_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CenterConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CenterConnectRequest) ProtoMessage() {}

func (x *CenterConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_center_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CenterConnectRequest.ProtoReflect.Descriptor instead.
func (*CenterConnectRequest) Descriptor() ([]byte, []int) {
	return file_pb_center_proto_rawDescGZIP(), []int{2}
}

func (x *CenterConnectRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type CenterConnectReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CenterConnectReply) Reset() {
	*x = CenterConnectReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_center_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CenterConnectReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CenterConnectReply) ProtoMessage() {}

func (x *CenterConnectReply) ProtoReflect() protoreflect.Message {
	mi := &file_pb_center_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CenterConnectReply.ProtoReflect.Descriptor instead.
func (*CenterConnectReply) Descriptor() ([]byte, []int) {
	return file_pb_center_proto_rawDescGZIP(), []int{3}
}

var File_pb_center_proto protoreflect.FileDescriptor

var file_pb_center_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x62, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x34, 0x0a, 0x0a, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x50, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x55, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x0c, 0x0a, 0x0a,
	0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x28, 0x0a, 0x14, 0x43, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x55, 0x72, 0x6c, 0x22, 0x14, 0x0a, 0x12, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x7c, 0x0a, 0x06, 0x43, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61,
	0x74, 0x12, 0x0f, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x50, 0x69,
	0x6e, 0x67, 0x1a, 0x0f, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x50,
	0x6f, 0x6e, 0x67, 0x22, 0x00, 0x28, 0x01, 0x12, 0x3f, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x12, 0x19, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x6d, 0x73, 0x67, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x6d, 0x73, 0x67,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_center_proto_rawDescOnce sync.Once
	file_pb_center_proto_rawDescData = file_pb_center_proto_rawDesc
)

func file_pb_center_proto_rawDescGZIP() []byte {
	file_pb_center_proto_rawDescOnce.Do(func() {
		file_pb_center_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_center_proto_rawDescData)
	})
	return file_pb_center_proto_rawDescData
}

var file_pb_center_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pb_center_proto_goTypes = []interface{}{
	(*CenterPing)(nil),           // 0: msg.CenterPing
	(*CenterPong)(nil),           // 1: msg.CenterPong
	(*CenterConnectRequest)(nil), // 2: msg.CenterConnectRequest
	(*CenterConnectReply)(nil),   // 3: msg.CenterConnectReply
}
var file_pb_center_proto_depIdxs = []int32{
	0, // 0: msg.Center.Heartbeat:input_type -> msg.CenterPing
	2, // 1: msg.Center.Connect:input_type -> msg.CenterConnectRequest
	1, // 2: msg.Center.Heartbeat:output_type -> msg.CenterPong
	3, // 3: msg.Center.Connect:output_type -> msg.CenterConnectReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_center_proto_init() }
func file_pb_center_proto_init() {
	if File_pb_center_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_center_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CenterPing); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_center_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CenterPong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_center_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CenterConnectRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_center_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CenterConnectReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_center_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_center_proto_goTypes,
		DependencyIndexes: file_pb_center_proto_depIdxs,
		MessageInfos:      file_pb_center_proto_msgTypes,
	}.Build()
	File_pb_center_proto = out.File
	file_pb_center_proto_rawDesc = nil
	file_pb_center_proto_goTypes = nil
	file_pb_center_proto_depIdxs = nil
}
