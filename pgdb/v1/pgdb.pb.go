// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: pgdb/v1/pgdb.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageOptions_Index_IndexMethod int32

const (
	MessageOptions_Index_INDEX_METHOD_UNSPECIFIED MessageOptions_Index_IndexMethod = 0
	MessageOptions_Index_INDEX_METHOD_BTREE       MessageOptions_Index_IndexMethod = 1
	MessageOptions_Index_INDEX_METHOD_GIN         MessageOptions_Index_IndexMethod = 2
	// Requires loading of BTREE_GIN extension:
	// https://www.postgresql.org/docs/current/btree-gin.html
	MessageOptions_Index_INDEX_METHOD_BTREE_GIN MessageOptions_Index_IndexMethod = 3
)

// Enum value maps for MessageOptions_Index_IndexMethod.
var (
	MessageOptions_Index_IndexMethod_name = map[int32]string{
		0: "INDEX_METHOD_UNSPECIFIED",
		1: "INDEX_METHOD_BTREE",
		2: "INDEX_METHOD_GIN",
		3: "INDEX_METHOD_BTREE_GIN",
	}
	MessageOptions_Index_IndexMethod_value = map[string]int32{
		"INDEX_METHOD_UNSPECIFIED": 0,
		"INDEX_METHOD_BTREE":       1,
		"INDEX_METHOD_GIN":         2,
		"INDEX_METHOD_BTREE_GIN":   3,
	}
)

func (x MessageOptions_Index_IndexMethod) Enum() *MessageOptions_Index_IndexMethod {
	p := new(MessageOptions_Index_IndexMethod)
	*p = x
	return p
}

func (x MessageOptions_Index_IndexMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageOptions_Index_IndexMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_pgdb_v1_pgdb_proto_enumTypes[0].Descriptor()
}

func (MessageOptions_Index_IndexMethod) Type() protoreflect.EnumType {
	return &file_pgdb_v1_pgdb_proto_enumTypes[0]
}

func (x MessageOptions_Index_IndexMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageOptions_Index_IndexMethod.Descriptor instead.
func (MessageOptions_Index_IndexMethod) EnumDescriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{0, 0, 0}
}

type FieldOptions_FullTextType int32

const (
	FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED FieldOptions_FullTextType = 0
	FieldOptions_FULL_TEXT_TYPE_EXACT       FieldOptions_FullTextType = 1
	// Best used for short display names
	FieldOptions_FULL_TEXT_TYPE_ENGLISH FieldOptions_FullTextType = 2
	// Removes short-tokens (<3 chars), useful for descirptions
	FieldOptions_FULL_TEXT_TYPE_ENGLISH_LONG FieldOptions_FullTextType = 3
)

// Enum value maps for FieldOptions_FullTextType.
var (
	FieldOptions_FullTextType_name = map[int32]string{
		0: "FULL_TEXT_TYPE_UNSPECIFIED",
		1: "FULL_TEXT_TYPE_EXACT",
		2: "FULL_TEXT_TYPE_ENGLISH",
		3: "FULL_TEXT_TYPE_ENGLISH_LONG",
	}
	FieldOptions_FullTextType_value = map[string]int32{
		"FULL_TEXT_TYPE_UNSPECIFIED":  0,
		"FULL_TEXT_TYPE_EXACT":        1,
		"FULL_TEXT_TYPE_ENGLISH":      2,
		"FULL_TEXT_TYPE_ENGLISH_LONG": 3,
	}
)

func (x FieldOptions_FullTextType) Enum() *FieldOptions_FullTextType {
	p := new(FieldOptions_FullTextType)
	*p = x
	return p
}

func (x FieldOptions_FullTextType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldOptions_FullTextType) Descriptor() protoreflect.EnumDescriptor {
	return file_pgdb_v1_pgdb_proto_enumTypes[1].Descriptor()
}

func (FieldOptions_FullTextType) Type() protoreflect.EnumType {
	return &file_pgdb_v1_pgdb_proto_enumTypes[1]
}

func (x FieldOptions_FullTextType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FieldOptions_FullTextType.Descriptor instead.
func (FieldOptions_FullTextType) EnumDescriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{1, 0}
}

type FieldOptions_FullTextWeight int32

const (
	FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED FieldOptions_FullTextWeight = 0
	FieldOptions_FULL_TEXT_WEIGHT_LOW         FieldOptions_FullTextWeight = 1
	FieldOptions_FULL_TEXT_WEIGHT_MED         FieldOptions_FullTextWeight = 2
	FieldOptions_FULL_TEXT_WEIGHT_HIGH        FieldOptions_FullTextWeight = 3
)

// Enum value maps for FieldOptions_FullTextWeight.
var (
	FieldOptions_FullTextWeight_name = map[int32]string{
		0: "FULL_TEXT_WEIGHT_UNSPECIFIED",
		1: "FULL_TEXT_WEIGHT_LOW",
		2: "FULL_TEXT_WEIGHT_MED",
		3: "FULL_TEXT_WEIGHT_HIGH",
	}
	FieldOptions_FullTextWeight_value = map[string]int32{
		"FULL_TEXT_WEIGHT_UNSPECIFIED": 0,
		"FULL_TEXT_WEIGHT_LOW":         1,
		"FULL_TEXT_WEIGHT_MED":         2,
		"FULL_TEXT_WEIGHT_HIGH":        3,
	}
)

func (x FieldOptions_FullTextWeight) Enum() *FieldOptions_FullTextWeight {
	p := new(FieldOptions_FullTextWeight)
	*p = x
	return p
}

func (x FieldOptions_FullTextWeight) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldOptions_FullTextWeight) Descriptor() protoreflect.EnumDescriptor {
	return file_pgdb_v1_pgdb_proto_enumTypes[2].Descriptor()
}

func (FieldOptions_FullTextWeight) Type() protoreflect.EnumType {
	return &file_pgdb_v1_pgdb_proto_enumTypes[2]
}

func (x FieldOptions_FullTextWeight) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FieldOptions_FullTextWeight.Descriptor instead.
func (FieldOptions_FullTextWeight) EnumDescriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{1, 1}
}

type FieldOptions_MessageBehavoir int32

const (
	FieldOptions_MESSAGE_BEHAVOIR_UNSPECIFIED FieldOptions_MessageBehavoir = 0
	FieldOptions_MESSAGE_BEHAVOIR_EXPAND      FieldOptions_MessageBehavoir = 1
	FieldOptions_MESSAGE_BEHAVOIR_OMIT        FieldOptions_MessageBehavoir = 2
	FieldOptions_MESSAGE_BEHAVOIR_JSONB       FieldOptions_MessageBehavoir = 3
	FieldOptions_MESSAGE_BEHAVOIR_VECTOR      FieldOptions_MessageBehavoir = 4
)

// Enum value maps for FieldOptions_MessageBehavoir.
var (
	FieldOptions_MessageBehavoir_name = map[int32]string{
		0: "MESSAGE_BEHAVOIR_UNSPECIFIED",
		1: "MESSAGE_BEHAVOIR_EXPAND",
		2: "MESSAGE_BEHAVOIR_OMIT",
		3: "MESSAGE_BEHAVOIR_JSONB",
		4: "MESSAGE_BEHAVOIR_VECTOR",
	}
	FieldOptions_MessageBehavoir_value = map[string]int32{
		"MESSAGE_BEHAVOIR_UNSPECIFIED": 0,
		"MESSAGE_BEHAVOIR_EXPAND":      1,
		"MESSAGE_BEHAVOIR_OMIT":        2,
		"MESSAGE_BEHAVOIR_JSONB":       3,
		"MESSAGE_BEHAVOIR_VECTOR":      4,
	}
)

func (x FieldOptions_MessageBehavoir) Enum() *FieldOptions_MessageBehavoir {
	p := new(FieldOptions_MessageBehavoir)
	*p = x
	return p
}

func (x FieldOptions_MessageBehavoir) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldOptions_MessageBehavoir) Descriptor() protoreflect.EnumDescriptor {
	return file_pgdb_v1_pgdb_proto_enumTypes[3].Descriptor()
}

func (FieldOptions_MessageBehavoir) Type() protoreflect.EnumType {
	return &file_pgdb_v1_pgdb_proto_enumTypes[3]
}

func (x FieldOptions_MessageBehavoir) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FieldOptions_MessageBehavoir.Descriptor instead.
func (FieldOptions_MessageBehavoir) EnumDescriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{1, 2}
}

type MessageOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Disabled bool                    `protobuf:"varint,1,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Indexes  []*MessageOptions_Index `protobuf:"bytes,2,rep,name=indexes,proto3" json:"indexes,omitempty"`
	// defaults to `tenant_id`.  Must be set if an object does not have a
	// `tenant_id` field.
	TenantIdField string `protobuf:"bytes,3,opt,name=tenant_id_field,json=tenantIdField,proto3" json:"tenant_id_field,omitempty"`
	// if this message is only used in nested messages, a subset of methods
	// will be generated.
	NestedOnly bool `protobuf:"varint,4,opt,name=nested_only,json=nestedOnly,proto3" json:"nested_only,omitempty"`
	// if this message is used then we create a partitioned table and partition by
	// tenant_id.
	Partitioned bool `protobuf:"varint,5,opt,name=partitioned,proto3" json:"partitioned,omitempty"`
}

func (x *MessageOptions) Reset() {
	*x = MessageOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pgdb_v1_pgdb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOptions) ProtoMessage() {}

func (x *MessageOptions) ProtoReflect() protoreflect.Message {
	mi := &file_pgdb_v1_pgdb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOptions.ProtoReflect.Descriptor instead.
func (*MessageOptions) Descriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{0}
}

func (x *MessageOptions) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *MessageOptions) GetIndexes() []*MessageOptions_Index {
	if x != nil {
		return x.Indexes
	}
	return nil
}

func (x *MessageOptions) GetTenantIdField() string {
	if x != nil {
		return x.TenantIdField
	}
	return ""
}

func (x *MessageOptions) GetNestedOnly() bool {
	if x != nil {
		return x.NestedOnly
	}
	return false
}

func (x *MessageOptions) GetPartitioned() bool {
	if x != nil {
		return x.Partitioned
	}
	return false
}

type FieldOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FullTextType    FieldOptions_FullTextType    `protobuf:"varint,1,opt,name=full_text_type,json=fullTextType,proto3,enum=pgdb.v1.FieldOptions_FullTextType" json:"full_text_type,omitempty"`
	FullTextWeight  FieldOptions_FullTextWeight  `protobuf:"varint,2,opt,name=full_text_weight,json=fullTextWeight,proto3,enum=pgdb.v1.FieldOptions_FullTextWeight" json:"full_text_weight,omitempty"`
	MessageBehavoir FieldOptions_MessageBehavoir `protobuf:"varint,3,opt,name=message_behavoir,json=messageBehavoir,proto3,enum=pgdb.v1.FieldOptions_MessageBehavoir" json:"message_behavoir,omitempty"`
	// vector size options
	VectorSize int32 `protobuf:"varint,4,opt,name=vector_size,json=vectorSize,proto3" json:"vector_size,omitempty"`
}

func (x *FieldOptions) Reset() {
	*x = FieldOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pgdb_v1_pgdb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOptions) ProtoMessage() {}

func (x *FieldOptions) ProtoReflect() protoreflect.Message {
	mi := &file_pgdb_v1_pgdb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOptions.ProtoReflect.Descriptor instead.
func (*FieldOptions) Descriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{1}
}

func (x *FieldOptions) GetFullTextType() FieldOptions_FullTextType {
	if x != nil {
		return x.FullTextType
	}
	return FieldOptions_FULL_TEXT_TYPE_UNSPECIFIED
}

func (x *FieldOptions) GetFullTextWeight() FieldOptions_FullTextWeight {
	if x != nil {
		return x.FullTextWeight
	}
	return FieldOptions_FULL_TEXT_WEIGHT_UNSPECIFIED
}

func (x *FieldOptions) GetMessageBehavoir() FieldOptions_MessageBehavoir {
	if x != nil {
		return x.MessageBehavoir
	}
	return FieldOptions_MESSAGE_BEHAVOIR_UNSPECIFIED
}

func (x *FieldOptions) GetVectorSize() int32 {
	if x != nil {
		return x.VectorSize
	}
	return 0
}

type MessageOptions_Index struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string                           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Method  MessageOptions_Index_IndexMethod `protobuf:"varint,2,opt,name=method,proto3,enum=pgdb.v1.MessageOptions_Index_IndexMethod" json:"method,omitempty"`
	Columns []string                         `protobuf:"bytes,3,rep,name=columns,proto3" json:"columns,omitempty"`
	// used to indicate the index by this name can be dropped
	Dropped bool `protobuf:"varint,4,opt,name=dropped,proto3" json:"dropped,omitempty"`
}

func (x *MessageOptions_Index) Reset() {
	*x = MessageOptions_Index{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pgdb_v1_pgdb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOptions_Index) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOptions_Index) ProtoMessage() {}

func (x *MessageOptions_Index) ProtoReflect() protoreflect.Message {
	mi := &file_pgdb_v1_pgdb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOptions_Index.ProtoReflect.Descriptor instead.
func (*MessageOptions_Index) Descriptor() ([]byte, []int) {
	return file_pgdb_v1_pgdb_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MessageOptions_Index) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MessageOptions_Index) GetMethod() MessageOptions_Index_IndexMethod {
	if x != nil {
		return x.Method
	}
	return MessageOptions_Index_INDEX_METHOD_UNSPECIFIED
}

func (x *MessageOptions_Index) GetColumns() []string {
	if x != nil {
		return x.Columns
	}
	return nil
}

func (x *MessageOptions_Index) GetDropped() bool {
	if x != nil {
		return x.Dropped
	}
	return false
}

var file_pgdb_v1_pgdb_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*MessageOptions)(nil),
		Field:         6010,
		Name:          "pgdb.v1.msg",
		Tag:           "bytes,6010,opt,name=msg",
		Filename:      "pgdb/v1/pgdb.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOptions)(nil),
		Field:         6010,
		Name:          "pgdb.v1.options",
		Tag:           "bytes,6010,opt,name=options",
		Filename:      "pgdb/v1/pgdb.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional pgdb.v1.MessageOptions msg = 6010;
	E_Msg = &file_pgdb_v1_pgdb_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional pgdb.v1.FieldOptions options = 6010;
	E_Options = &file_pgdb_v1_pgdb_proto_extTypes[1]
)

var File_pgdb_v1_pgdb_proto protoreflect.FileDescriptor

var file_pgdb_v1_pgdb_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x67, 0x64, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xdc, 0x03, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x37,
	0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x07,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4f, 0x6e, 0x6c, 0x79,
	0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x65, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x65, 0x64, 0x1a, 0x89, 0x02, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x41, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x29, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x2e,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x22, 0x75, 0x0a, 0x0b, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1c, 0x0a, 0x18, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x5f,
	0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x5f, 0x4d, 0x45,
	0x54, 0x48, 0x4f, 0x44, 0x5f, 0x42, 0x54, 0x52, 0x45, 0x45, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10,
	0x49, 0x4e, 0x44, 0x45, 0x58, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x47, 0x49, 0x4e,
	0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x5f, 0x4d, 0x45, 0x54, 0x48,
	0x4f, 0x44, 0x5f, 0x42, 0x54, 0x52, 0x45, 0x45, 0x5f, 0x47, 0x49, 0x4e, 0x10, 0x03, 0x22, 0xce,
	0x05, 0x0a, 0x0c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x48, 0x0a, 0x0e, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x46,
	0x75, 0x6c, 0x6c, 0x54, 0x65, 0x78, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x66, 0x75, 0x6c,
	0x6c, 0x54, 0x65, 0x78, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x4e, 0x0a, 0x10, 0x66, 0x75, 0x6c,
	0x6c, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x54,
	0x65, 0x78, 0x74, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x0e, 0x66, 0x75, 0x6c, 0x6c, 0x54,
	0x65, 0x78, 0x74, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x50, 0x0a, 0x10, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x6f, 0x69, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x65, 0x68, 0x61, 0x76, 0x6f, 0x69, 0x72, 0x52, 0x0f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x42, 0x65, 0x68, 0x61, 0x76, 0x6f, 0x69, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x85, 0x01, 0x0a,
	0x0c, 0x46, 0x75, 0x6c, 0x6c, 0x54, 0x65, 0x78, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a,
	0x1a, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a,
	0x14, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x45, 0x58, 0x41, 0x43, 0x54, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x46, 0x55, 0x4c, 0x4c, 0x5f,
	0x54, 0x45, 0x58, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x4e, 0x47, 0x4c, 0x49, 0x53,
	0x48, 0x10, 0x02, 0x12, 0x1f, 0x0a, 0x1b, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x4e, 0x47, 0x4c, 0x49, 0x53, 0x48, 0x5f, 0x4c, 0x4f,
	0x4e, 0x47, 0x10, 0x03, 0x22, 0x81, 0x01, 0x0a, 0x0e, 0x46, 0x75, 0x6c, 0x6c, 0x54, 0x65, 0x78,
	0x74, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x20, 0x0a, 0x1c, 0x46, 0x55, 0x4c, 0x4c, 0x5f,
	0x54, 0x45, 0x58, 0x54, 0x5f, 0x57, 0x45, 0x49, 0x47, 0x48, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x46, 0x55, 0x4c,
	0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x5f, 0x57, 0x45, 0x49, 0x47, 0x48, 0x54, 0x5f, 0x4c, 0x4f,
	0x57, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54,
	0x5f, 0x57, 0x45, 0x49, 0x47, 0x48, 0x54, 0x5f, 0x4d, 0x45, 0x44, 0x10, 0x02, 0x12, 0x19, 0x0a,
	0x15, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x5f, 0x57, 0x45, 0x49, 0x47, 0x48,
	0x54, 0x5f, 0x48, 0x49, 0x47, 0x48, 0x10, 0x03, 0x22, 0xa4, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x42, 0x65, 0x68, 0x61, 0x76, 0x6f, 0x69, 0x72, 0x12, 0x20, 0x0a, 0x1c,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x45, 0x48, 0x41, 0x56, 0x4f, 0x49, 0x52,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1b,
	0x0a, 0x17, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x45, 0x48, 0x41, 0x56, 0x4f,
	0x49, 0x52, 0x5f, 0x45, 0x58, 0x50, 0x41, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x4d,
	0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x45, 0x48, 0x41, 0x56, 0x4f, 0x49, 0x52, 0x5f,
	0x4f, 0x4d, 0x49, 0x54, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x42, 0x45, 0x48, 0x41, 0x56, 0x4f, 0x49, 0x52, 0x5f, 0x4a, 0x53, 0x4f, 0x4e, 0x42,
	0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x45,
	0x48, 0x41, 0x56, 0x4f, 0x49, 0x52, 0x5f, 0x56, 0x45, 0x43, 0x54, 0x4f, 0x52, 0x10, 0x04, 0x3a,
	0x4b, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xfa, 0x2e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x3a, 0x4f, 0x0a, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xfa, 0x2e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x81, 0x01,
	0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x50,
	0x67, 0x64, 0x62, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x75, 0x63, 0x74, 0x6f, 0x6e, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x70, 0x67, 0x64, 0x62, 0x2f, 0x70,
	0x67, 0x64, 0x62, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x50,
	0x67, 0x64, 0x62, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x07, 0x50, 0x67, 0x64, 0x62, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x13, 0x50, 0x67, 0x64, 0x62, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x50, 0x67, 0x64, 0x62, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pgdb_v1_pgdb_proto_rawDescOnce sync.Once
	file_pgdb_v1_pgdb_proto_rawDescData = file_pgdb_v1_pgdb_proto_rawDesc
)

func file_pgdb_v1_pgdb_proto_rawDescGZIP() []byte {
	file_pgdb_v1_pgdb_proto_rawDescOnce.Do(func() {
		file_pgdb_v1_pgdb_proto_rawDescData = protoimpl.X.CompressGZIP(file_pgdb_v1_pgdb_proto_rawDescData)
	})
	return file_pgdb_v1_pgdb_proto_rawDescData
}

var file_pgdb_v1_pgdb_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_pgdb_v1_pgdb_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pgdb_v1_pgdb_proto_goTypes = []interface{}{
	(MessageOptions_Index_IndexMethod)(0), // 0: pgdb.v1.MessageOptions.Index.IndexMethod
	(FieldOptions_FullTextType)(0),        // 1: pgdb.v1.FieldOptions.FullTextType
	(FieldOptions_FullTextWeight)(0),      // 2: pgdb.v1.FieldOptions.FullTextWeight
	(FieldOptions_MessageBehavoir)(0),     // 3: pgdb.v1.FieldOptions.MessageBehavoir
	(*MessageOptions)(nil),                // 4: pgdb.v1.MessageOptions
	(*FieldOptions)(nil),                  // 5: pgdb.v1.FieldOptions
	(*MessageOptions_Index)(nil),          // 6: pgdb.v1.MessageOptions.Index
	(*descriptorpb.MessageOptions)(nil),   // 7: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),     // 8: google.protobuf.FieldOptions
}
var file_pgdb_v1_pgdb_proto_depIdxs = []int32{
	6, // 0: pgdb.v1.MessageOptions.indexes:type_name -> pgdb.v1.MessageOptions.Index
	1, // 1: pgdb.v1.FieldOptions.full_text_type:type_name -> pgdb.v1.FieldOptions.FullTextType
	2, // 2: pgdb.v1.FieldOptions.full_text_weight:type_name -> pgdb.v1.FieldOptions.FullTextWeight
	3, // 3: pgdb.v1.FieldOptions.message_behavoir:type_name -> pgdb.v1.FieldOptions.MessageBehavoir
	0, // 4: pgdb.v1.MessageOptions.Index.method:type_name -> pgdb.v1.MessageOptions.Index.IndexMethod
	7, // 5: pgdb.v1.msg:extendee -> google.protobuf.MessageOptions
	8, // 6: pgdb.v1.options:extendee -> google.protobuf.FieldOptions
	4, // 7: pgdb.v1.msg:type_name -> pgdb.v1.MessageOptions
	5, // 8: pgdb.v1.options:type_name -> pgdb.v1.FieldOptions
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	7, // [7:9] is the sub-list for extension type_name
	5, // [5:7] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pgdb_v1_pgdb_proto_init() }
func file_pgdb_v1_pgdb_proto_init() {
	if File_pgdb_v1_pgdb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pgdb_v1_pgdb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOptions); i {
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
		file_pgdb_v1_pgdb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldOptions); i {
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
		file_pgdb_v1_pgdb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOptions_Index); i {
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
			RawDescriptor: file_pgdb_v1_pgdb_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   3,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_pgdb_v1_pgdb_proto_goTypes,
		DependencyIndexes: file_pgdb_v1_pgdb_proto_depIdxs,
		EnumInfos:         file_pgdb_v1_pgdb_proto_enumTypes,
		MessageInfos:      file_pgdb_v1_pgdb_proto_msgTypes,
		ExtensionInfos:    file_pgdb_v1_pgdb_proto_extTypes,
	}.Build()
	File_pgdb_v1_pgdb_proto = out.File
	file_pgdb_v1_pgdb_proto_rawDesc = nil
	file_pgdb_v1_pgdb_proto_goTypes = nil
	file_pgdb_v1_pgdb_proto_depIdxs = nil
}
