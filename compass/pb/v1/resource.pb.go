// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: pb/v1/resource.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PodsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PodsRequest) Reset() {
	*x = PodsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_v1_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodsRequest) ProtoMessage() {}

func (x *PodsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_v1_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodsRequest.ProtoReflect.Descriptor instead.
func (*PodsRequest) Descriptor() ([]byte, []int) {
	return file_pb_v1_resource_proto_rawDescGZIP(), []int{0}
}

type PodsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PodsResponse) Reset() {
	*x = PodsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_v1_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodsResponse) ProtoMessage() {}

func (x *PodsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_v1_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodsResponse.ProtoReflect.Descriptor instead.
func (*PodsResponse) Descriptor() ([]byte, []int) {
	return file_pb_v1_resource_proto_rawDescGZIP(), []int{1}
}

var File_pb_v1_resource_proto protoreflect.FileDescriptor

var file_pb_v1_resource_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x50, 0x6f, 0x64, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x6f, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x3e, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x50, 0x6f, 0x64, 0x73, 0x12,
	0x0f, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x10, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_v1_resource_proto_rawDescOnce sync.Once
	file_pb_v1_resource_proto_rawDescData = file_pb_v1_resource_proto_rawDesc
)

func file_pb_v1_resource_proto_rawDescGZIP() []byte {
	file_pb_v1_resource_proto_rawDescOnce.Do(func() {
		file_pb_v1_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_v1_resource_proto_rawDescData)
	})
	return file_pb_v1_resource_proto_rawDescData
}

var file_pb_v1_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_v1_resource_proto_goTypes = []interface{}{
	(*PodsRequest)(nil),  // 0: v1.PodsRequest
	(*PodsResponse)(nil), // 1: v1.PodsResponse
}
var file_pb_v1_resource_proto_depIdxs = []int32{
	0, // 0: v1.ResourceService.Pods:input_type -> v1.PodsRequest
	1, // 1: v1.ResourceService.Pods:output_type -> v1.PodsResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_v1_resource_proto_init() }
func file_pb_v1_resource_proto_init() {
	if File_pb_v1_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_v1_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodsRequest); i {
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
		file_pb_v1_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodsResponse); i {
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
			RawDescriptor: file_pb_v1_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_v1_resource_proto_goTypes,
		DependencyIndexes: file_pb_v1_resource_proto_depIdxs,
		MessageInfos:      file_pb_v1_resource_proto_msgTypes,
	}.Build()
	File_pb_v1_resource_proto = out.File
	file_pb_v1_resource_proto_rawDesc = nil
	file_pb_v1_resource_proto_goTypes = nil
	file_pb_v1_resource_proto_depIdxs = nil
}
