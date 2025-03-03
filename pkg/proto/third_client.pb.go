// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: third_client.proto

package proto

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

type ThirdClientInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Alias       string  `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
	Name        string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description string  `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	IsDraft     bool    `protobuf:"varint,5,opt,name=is_draft,json=isDraft,proto3" json:"is_draft,omitempty"`
	AccountId   *uint64 `protobuf:"varint,6,opt,name=account_id,json=accountId,proto3,oneof" json:"account_id,omitempty"`
}

func (x *ThirdClientInfo) Reset() {
	*x = ThirdClientInfo{}
	mi := &file_third_client_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ThirdClientInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThirdClientInfo) ProtoMessage() {}

func (x *ThirdClientInfo) ProtoReflect() protoreflect.Message {
	mi := &file_third_client_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThirdClientInfo.ProtoReflect.Descriptor instead.
func (*ThirdClientInfo) Descriptor() ([]byte, []int) {
	return file_third_client_proto_rawDescGZIP(), []int{0}
}

func (x *ThirdClientInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ThirdClientInfo) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *ThirdClientInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ThirdClientInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ThirdClientInfo) GetIsDraft() bool {
	if x != nil {
		return x.IsDraft
	}
	return false
}

func (x *ThirdClientInfo) GetAccountId() uint64 {
	if x != nil && x.AccountId != nil {
		return *x.AccountId
	}
	return 0
}

type GetThirdClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     *uint64 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Alias  *string `protobuf:"bytes,2,opt,name=alias,proto3,oneof" json:"alias,omitempty"`
	Secret *string `protobuf:"bytes,3,opt,name=secret,proto3,oneof" json:"secret,omitempty"`
}

func (x *GetThirdClientRequest) Reset() {
	*x = GetThirdClientRequest{}
	mi := &file_third_client_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetThirdClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetThirdClientRequest) ProtoMessage() {}

func (x *GetThirdClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_third_client_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetThirdClientRequest.ProtoReflect.Descriptor instead.
func (*GetThirdClientRequest) Descriptor() ([]byte, []int) {
	return file_third_client_proto_rawDescGZIP(), []int{1}
}

func (x *GetThirdClientRequest) GetId() uint64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *GetThirdClientRequest) GetAlias() string {
	if x != nil && x.Alias != nil {
		return *x.Alias
	}
	return ""
}

func (x *GetThirdClientRequest) GetSecret() string {
	if x != nil && x.Secret != nil {
		return *x.Secret
	}
	return ""
}

type GetThirdClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *ThirdClientInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *GetThirdClientResponse) Reset() {
	*x = GetThirdClientResponse{}
	mi := &file_third_client_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetThirdClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetThirdClientResponse) ProtoMessage() {}

func (x *GetThirdClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_third_client_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetThirdClientResponse.ProtoReflect.Descriptor instead.
func (*GetThirdClientResponse) Descriptor() ([]byte, []int) {
	return file_third_client_proto_rawDescGZIP(), []int{2}
}

func (x *GetThirdClientResponse) GetInfo() *ThirdClientInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

var File_third_client_proto protoreflect.FileDescriptor

var file_third_client_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x68, 0x69, 0x72, 0x64, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x01, 0x0a, 0x0f,
	0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x61, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x69,
	0x73, 0x5f, 0x64, 0x72, 0x61, 0x66, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69,
	0x73, 0x44, 0x72, 0x61, 0x66, 0x74, 0x12, 0x22, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x09, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x80, 0x01, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x48,
	0x00, 0x52, 0x02, 0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73,
	0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x88, 0x01, 0x01,
	0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69, 0x64, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x6c, 0x69, 0x61,
	0x73, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x44, 0x0a, 0x16,
	0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x68, 0x69,
	0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e,
	0x66, 0x6f, 0x32, 0x65, 0x0a, 0x12, 0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54,
	0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_third_client_proto_rawDescOnce sync.Once
	file_third_client_proto_rawDescData = file_third_client_proto_rawDesc
)

func file_third_client_proto_rawDescGZIP() []byte {
	file_third_client_proto_rawDescOnce.Do(func() {
		file_third_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_third_client_proto_rawDescData)
	})
	return file_third_client_proto_rawDescData
}

var file_third_client_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_third_client_proto_goTypes = []any{
	(*ThirdClientInfo)(nil),        // 0: proto.ThirdClientInfo
	(*GetThirdClientRequest)(nil),  // 1: proto.GetThirdClientRequest
	(*GetThirdClientResponse)(nil), // 2: proto.GetThirdClientResponse
}
var file_third_client_proto_depIdxs = []int32{
	0, // 0: proto.GetThirdClientResponse.info:type_name -> proto.ThirdClientInfo
	1, // 1: proto.ThirdClientService.GetThirdClient:input_type -> proto.GetThirdClientRequest
	2, // 2: proto.ThirdClientService.GetThirdClient:output_type -> proto.GetThirdClientResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_third_client_proto_init() }
func file_third_client_proto_init() {
	if File_third_client_proto != nil {
		return
	}
	file_third_client_proto_msgTypes[0].OneofWrappers = []any{}
	file_third_client_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_third_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_third_client_proto_goTypes,
		DependencyIndexes: file_third_client_proto_depIdxs,
		MessageInfos:      file_third_client_proto_msgTypes,
	}.Build()
	File_third_client_proto = out.File
	file_third_client_proto_rawDesc = nil
	file_third_client_proto_goTypes = nil
	file_third_client_proto_depIdxs = nil
}
