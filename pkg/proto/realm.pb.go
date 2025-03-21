// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: realm.proto

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

type RealmInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Alias        string `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
	Name         string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description  string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Avatar       string `protobuf:"bytes,6,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Banner       string `protobuf:"bytes,7,opt,name=banner,proto3" json:"banner,omitempty"`
	IsPublic     bool   `protobuf:"varint,9,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	IsCommunity  bool   `protobuf:"varint,10,opt,name=is_community,json=isCommunity,proto3" json:"is_community,omitempty"`
	AccessPolicy []byte `protobuf:"bytes,11,opt,name=access_policy,json=accessPolicy,proto3" json:"access_policy,omitempty"`
}

func (x *RealmInfo) Reset() {
	*x = RealmInfo{}
	mi := &file_realm_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RealmInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmInfo) ProtoMessage() {}

func (x *RealmInfo) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmInfo.ProtoReflect.Descriptor instead.
func (*RealmInfo) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{0}
}

func (x *RealmInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RealmInfo) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *RealmInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RealmInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RealmInfo) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *RealmInfo) GetBanner() string {
	if x != nil {
		return x.Banner
	}
	return ""
}

func (x *RealmInfo) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *RealmInfo) GetIsCommunity() bool {
	if x != nil {
		return x.IsCommunity
	}
	return false
}

func (x *RealmInfo) GetAccessPolicy() []byte {
	if x != nil {
		return x.AccessPolicy
	}
	return nil
}

type ListRealmRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []uint64 `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
}

func (x *ListRealmRequest) Reset() {
	*x = ListRealmRequest{}
	mi := &file_realm_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRealmRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRealmRequest) ProtoMessage() {}

func (x *ListRealmRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRealmRequest.ProtoReflect.Descriptor instead.
func (*ListRealmRequest) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{1}
}

func (x *ListRealmRequest) GetId() []uint64 {
	if x != nil {
		return x.Id
	}
	return nil
}

type LookupUserRealmRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	IncludePublic *bool  `protobuf:"varint,2,opt,name=include_public,json=includePublic,proto3,oneof" json:"include_public,omitempty"`
}

func (x *LookupUserRealmRequest) Reset() {
	*x = LookupUserRealmRequest{}
	mi := &file_realm_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LookupUserRealmRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupUserRealmRequest) ProtoMessage() {}

func (x *LookupUserRealmRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupUserRealmRequest.ProtoReflect.Descriptor instead.
func (*LookupUserRealmRequest) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{2}
}

func (x *LookupUserRealmRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *LookupUserRealmRequest) GetIncludePublic() bool {
	if x != nil && x.IncludePublic != nil {
		return *x.IncludePublic
	}
	return false
}

type LookupRealmRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          *uint64 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Alias       *string `protobuf:"bytes,2,opt,name=alias,proto3,oneof" json:"alias,omitempty"`
	IsPublic    *bool   `protobuf:"varint,3,opt,name=is_public,json=isPublic,proto3,oneof" json:"is_public,omitempty"`
	IsCommunity *bool   `protobuf:"varint,4,opt,name=is_community,json=isCommunity,proto3,oneof" json:"is_community,omitempty"`
}

func (x *LookupRealmRequest) Reset() {
	*x = LookupRealmRequest{}
	mi := &file_realm_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LookupRealmRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupRealmRequest) ProtoMessage() {}

func (x *LookupRealmRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupRealmRequest.ProtoReflect.Descriptor instead.
func (*LookupRealmRequest) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{3}
}

func (x *LookupRealmRequest) GetId() uint64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *LookupRealmRequest) GetAlias() string {
	if x != nil && x.Alias != nil {
		return *x.Alias
	}
	return ""
}

func (x *LookupRealmRequest) GetIsPublic() bool {
	if x != nil && x.IsPublic != nil {
		return *x.IsPublic
	}
	return false
}

func (x *LookupRealmRequest) GetIsCommunity() bool {
	if x != nil && x.IsCommunity != nil {
		return *x.IsCommunity
	}
	return false
}

type ListRealmResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*RealmInfo `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ListRealmResponse) Reset() {
	*x = ListRealmResponse{}
	mi := &file_realm_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRealmResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRealmResponse) ProtoMessage() {}

func (x *ListRealmResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRealmResponse.ProtoReflect.Descriptor instead.
func (*ListRealmResponse) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{4}
}

func (x *ListRealmResponse) GetData() []*RealmInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

type RealmMemberLookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RealmId *uint64 `protobuf:"varint,1,opt,name=realm_id,json=realmId,proto3,oneof" json:"realm_id,omitempty"`
	UserId  *uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
}

func (x *RealmMemberLookupRequest) Reset() {
	*x = RealmMemberLookupRequest{}
	mi := &file_realm_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RealmMemberLookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmMemberLookupRequest) ProtoMessage() {}

func (x *RealmMemberLookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmMemberLookupRequest.ProtoReflect.Descriptor instead.
func (*RealmMemberLookupRequest) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{5}
}

func (x *RealmMemberLookupRequest) GetRealmId() uint64 {
	if x != nil && x.RealmId != nil {
		return *x.RealmId
	}
	return 0
}

func (x *RealmMemberLookupRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type RealmMemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	RealmId    uint64 `protobuf:"varint,2,opt,name=realm_id,json=realmId,proto3" json:"realm_id,omitempty"`
	UserId     uint64 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PowerLevel int32  `protobuf:"varint,4,opt,name=power_level,json=powerLevel,proto3" json:"power_level,omitempty"`
}

func (x *RealmMemberInfo) Reset() {
	*x = RealmMemberInfo{}
	mi := &file_realm_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RealmMemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmMemberInfo) ProtoMessage() {}

func (x *RealmMemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmMemberInfo.ProtoReflect.Descriptor instead.
func (*RealmMemberInfo) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{6}
}

func (x *RealmMemberInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RealmMemberInfo) GetRealmId() uint64 {
	if x != nil {
		return x.RealmId
	}
	return 0
}

func (x *RealmMemberInfo) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *RealmMemberInfo) GetPowerLevel() int32 {
	if x != nil {
		return x.PowerLevel
	}
	return 0
}

type ListRealmMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*RealmMemberInfo `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ListRealmMemberResponse) Reset() {
	*x = ListRealmMemberResponse{}
	mi := &file_realm_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRealmMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRealmMemberResponse) ProtoMessage() {}

func (x *ListRealmMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRealmMemberResponse.ProtoReflect.Descriptor instead.
func (*ListRealmMemberResponse) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{7}
}

func (x *ListRealmMemberResponse) GetData() []*RealmMemberInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

type CheckRealmPermRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RealmId    uint64 `protobuf:"varint,1,opt,name=realm_id,json=realmId,proto3" json:"realm_id,omitempty"`
	UserId     uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PowerLevel int32  `protobuf:"varint,3,opt,name=power_level,json=powerLevel,proto3" json:"power_level,omitempty"`
}

func (x *CheckRealmPermRequest) Reset() {
	*x = CheckRealmPermRequest{}
	mi := &file_realm_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckRealmPermRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRealmPermRequest) ProtoMessage() {}

func (x *CheckRealmPermRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRealmPermRequest.ProtoReflect.Descriptor instead.
func (*CheckRealmPermRequest) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{8}
}

func (x *CheckRealmPermRequest) GetRealmId() uint64 {
	if x != nil {
		return x.RealmId
	}
	return 0
}

func (x *CheckRealmPermRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CheckRealmPermRequest) GetPowerLevel() int32 {
	if x != nil {
		return x.PowerLevel
	}
	return 0
}

type CheckRealmPermResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsSuccess bool `protobuf:"varint,1,opt,name=is_success,json=isSuccess,proto3" json:"is_success,omitempty"`
}

func (x *CheckRealmPermResponse) Reset() {
	*x = CheckRealmPermResponse{}
	mi := &file_realm_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckRealmPermResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRealmPermResponse) ProtoMessage() {}

func (x *CheckRealmPermResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realm_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRealmPermResponse.ProtoReflect.Descriptor instead.
func (*CheckRealmPermResponse) Descriptor() ([]byte, []int) {
	return file_realm_proto_rawDescGZIP(), []int{9}
}

func (x *CheckRealmPermResponse) GetIsSuccess() bool {
	if x != nil {
		return x.IsSuccess
	}
	return false
}

var File_realm_proto protoreflect.FileDescriptor

var file_realm_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfc, 0x01, 0x0a, 0x09, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x61, 0x6e, 0x6e, 0x65, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x12, 0x1b,
	0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x69,
	0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0b, 0x69, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x12, 0x23,
	0x0a, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x22, 0x22, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x70, 0x0a, 0x16, 0x4c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x0e, 0x69, 0x6e,
	0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x00, 0x52, 0x0d, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x88, 0x01, 0x01, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x69, 0x6e, 0x63, 0x6c, 0x75,
	0x64, 0x65, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x22, 0xbe, 0x01, 0x0a, 0x12, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x02,
	0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x88, 0x01, 0x01,
	0x12, 0x20, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x88,
	0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x03, 0x52, 0x0b, 0x69, 0x73, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69,
	0x64, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x69, 0x73,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x22, 0x39, 0x0a, 0x11, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x71, 0x0a, 0x18, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x48, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42,
	0x0b, 0x0a, 0x09, 0x5f, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x76, 0x0a, 0x0f, 0x52, 0x65, 0x61, 0x6c,
	0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x72,
	0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72,
	0x65, 0x61, 0x6c, 0x6d, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c,
	0x22, 0x45, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x6c, 0x0a, 0x15, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x61, 0x6c, 0x6d, 0x50, 0x65, 0x72, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x6f, 0x77, 0x65, 0x72,
	0x4c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0x37, 0x0a, 0x16, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65,
	0x61, 0x6c, 0x6d, 0x50, 0x65, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0xa3,
	0x04, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4f, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x4b, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x77, 0x6e, 0x65, 0x64, 0x52, 0x65, 0x61,
	0x6c, 0x6d, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x40, 0x0a,
	0x09, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x39, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x19, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x65, 0x61, 0x6c, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0f, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x4b, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c,
	0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x55, 0x0a,
	0x14, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x50, 0x65, 0x72, 0x6d, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x50, 0x65, 0x72, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x50, 0x65, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_realm_proto_rawDescOnce sync.Once
	file_realm_proto_rawDescData = file_realm_proto_rawDesc
)

func file_realm_proto_rawDescGZIP() []byte {
	file_realm_proto_rawDescOnce.Do(func() {
		file_realm_proto_rawDescData = protoimpl.X.CompressGZIP(file_realm_proto_rawDescData)
	})
	return file_realm_proto_rawDescData
}

var file_realm_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_realm_proto_goTypes = []any{
	(*RealmInfo)(nil),                // 0: proto.RealmInfo
	(*ListRealmRequest)(nil),         // 1: proto.ListRealmRequest
	(*LookupUserRealmRequest)(nil),   // 2: proto.LookupUserRealmRequest
	(*LookupRealmRequest)(nil),       // 3: proto.LookupRealmRequest
	(*ListRealmResponse)(nil),        // 4: proto.ListRealmResponse
	(*RealmMemberLookupRequest)(nil), // 5: proto.RealmMemberLookupRequest
	(*RealmMemberInfo)(nil),          // 6: proto.RealmMemberInfo
	(*ListRealmMemberResponse)(nil),  // 7: proto.ListRealmMemberResponse
	(*CheckRealmPermRequest)(nil),    // 8: proto.CheckRealmPermRequest
	(*CheckRealmPermResponse)(nil),   // 9: proto.CheckRealmPermResponse
}
var file_realm_proto_depIdxs = []int32{
	0, // 0: proto.ListRealmResponse.data:type_name -> proto.RealmInfo
	6, // 1: proto.ListRealmMemberResponse.data:type_name -> proto.RealmMemberInfo
	2, // 2: proto.RealmService.ListAvailableRealm:input_type -> proto.LookupUserRealmRequest
	2, // 3: proto.RealmService.ListOwnedRealm:input_type -> proto.LookupUserRealmRequest
	1, // 4: proto.RealmService.ListRealm:input_type -> proto.ListRealmRequest
	3, // 5: proto.RealmService.GetRealm:input_type -> proto.LookupRealmRequest
	5, // 6: proto.RealmService.ListRealmMember:input_type -> proto.RealmMemberLookupRequest
	5, // 7: proto.RealmService.GetRealmMember:input_type -> proto.RealmMemberLookupRequest
	8, // 8: proto.RealmService.CheckRealmMemberPerm:input_type -> proto.CheckRealmPermRequest
	4, // 9: proto.RealmService.ListAvailableRealm:output_type -> proto.ListRealmResponse
	4, // 10: proto.RealmService.ListOwnedRealm:output_type -> proto.ListRealmResponse
	4, // 11: proto.RealmService.ListRealm:output_type -> proto.ListRealmResponse
	0, // 12: proto.RealmService.GetRealm:output_type -> proto.RealmInfo
	7, // 13: proto.RealmService.ListRealmMember:output_type -> proto.ListRealmMemberResponse
	6, // 14: proto.RealmService.GetRealmMember:output_type -> proto.RealmMemberInfo
	9, // 15: proto.RealmService.CheckRealmMemberPerm:output_type -> proto.CheckRealmPermResponse
	9, // [9:16] is the sub-list for method output_type
	2, // [2:9] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_realm_proto_init() }
func file_realm_proto_init() {
	if File_realm_proto != nil {
		return
	}
	file_realm_proto_msgTypes[2].OneofWrappers = []any{}
	file_realm_proto_msgTypes[3].OneofWrappers = []any{}
	file_realm_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_realm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_realm_proto_goTypes,
		DependencyIndexes: file_realm_proto_depIdxs,
		MessageInfos:      file_realm_proto_msgTypes,
	}.Build()
	File_realm_proto = out.File
	file_realm_proto_rawDesc = nil
	file_realm_proto_goTypes = nil
	file_realm_proto_depIdxs = nil
}
