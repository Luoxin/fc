// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: nozomi.proto

package nozomi

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

type NotifyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyReq) Reset() {
	*x = NotifyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nozomi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyReq) ProtoMessage() {}

func (x *NotifyReq) ProtoReflect() protoreflect.Message {
	mi := &file_nozomi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyReq.ProtoReflect.Descriptor instead.
func (*NotifyReq) Descriptor() ([]byte, []int) {
	return file_nozomi_proto_rawDescGZIP(), []int{0}
}

type NotifyRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyRsp) Reset() {
	*x = NotifyRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nozomi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyRsp) ProtoMessage() {}

func (x *NotifyRsp) ProtoReflect() protoreflect.Message {
	mi := &file_nozomi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyRsp.ProtoReflect.Descriptor instead.
func (*NotifyRsp) Descriptor() ([]byte, []int) {
	return file_nozomi_proto_rawDescGZIP(), []int{1}
}

var File_nozomi_proto protoreflect.FileDescriptor

var file_nozomi_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69, 0x22, 0x0b, 0x0a, 0x09, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x52, 0x65, 0x71, 0x22, 0x0b, 0x0a, 0x09, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x73, 0x70,
	0x32, 0x3a, 0x0a, 0x06, 0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69, 0x12, 0x30, 0x0a, 0x06, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x79, 0x12, 0x11, 0x2e, 0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69,
	0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x73, 0x70, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b,
	0x62, 0x61, 0x73, 0x65, 0x2f, 0x6e, 0x6f, 0x7a, 0x6f, 0x6d, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_nozomi_proto_rawDescOnce sync.Once
	file_nozomi_proto_rawDescData = file_nozomi_proto_rawDesc
)

func file_nozomi_proto_rawDescGZIP() []byte {
	file_nozomi_proto_rawDescOnce.Do(func() {
		file_nozomi_proto_rawDescData = protoimpl.X.CompressGZIP(file_nozomi_proto_rawDescData)
	})
	return file_nozomi_proto_rawDescData
}

var file_nozomi_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_nozomi_proto_goTypes = []interface{}{
	(*NotifyReq)(nil), // 0: nozomi.NotifyReq
	(*NotifyRsp)(nil), // 1: nozomi.NotifyRsp
}
var file_nozomi_proto_depIdxs = []int32{
	0, // 0: nozomi.nozomi.Notify:input_type -> nozomi.NotifyReq
	1, // 1: nozomi.nozomi.Notify:output_type -> nozomi.NotifyRsp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_nozomi_proto_init() }
func file_nozomi_proto_init() {
	if File_nozomi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nozomi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyReq); i {
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
		file_nozomi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyRsp); i {
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
			RawDescriptor: file_nozomi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nozomi_proto_goTypes,
		DependencyIndexes: file_nozomi_proto_depIdxs,
		MessageInfos:      file_nozomi_proto_msgTypes,
	}.Build()
	File_nozomi_proto = out.File
	file_nozomi_proto_rawDesc = nil
	file_nozomi_proto_goTypes = nil
	file_nozomi_proto_depIdxs = nil
}