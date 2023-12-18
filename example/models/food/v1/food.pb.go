// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: models/food/v1/food.proto

package v1

import (
	_ "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	_ "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Pasta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId  string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id        string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

func (x *Pasta) Reset() {
	*x = Pasta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_food_v1_food_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pasta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pasta) ProtoMessage() {}

func (x *Pasta) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pasta.ProtoReflect.Descriptor instead.
func (*Pasta) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{0}
}

func (x *Pasta) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *Pasta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Pasta) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Pasta) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Pasta) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

type PastaIngredient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId     string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	IngredientId string                 `protobuf:"bytes,2,opt,name=ingredient_id,json=ingredientId,proto3" json:"ingredient_id,omitempty"`
	CreatedAt    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	PastaId      string                 `protobuf:"bytes,6,opt,name=pasta_id,json=pastaId,proto3" json:"pasta_id,omitempty"`
}

func (x *PastaIngredient) Reset() {
	*x = PastaIngredient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_food_v1_food_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PastaIngredient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastaIngredient) ProtoMessage() {}

func (x *PastaIngredient) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastaIngredient.ProtoReflect.Descriptor instead.
func (*PastaIngredient) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{1}
}

func (x *PastaIngredient) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *PastaIngredient) GetIngredientId() string {
	if x != nil {
		return x.IngredientId
	}
	return ""
}

func (x *PastaIngredient) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *PastaIngredient) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *PastaIngredient) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *PastaIngredient) GetPastaId() string {
	if x != nil {
		return x.PastaId
	}
	return ""
}

type SauceIngerdient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId  string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id        string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

func (x *SauceIngerdient) Reset() {
	*x = SauceIngerdient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_food_v1_food_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SauceIngerdient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SauceIngerdient) ProtoMessage() {}

func (x *SauceIngerdient) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SauceIngerdient.ProtoReflect.Descriptor instead.
func (*SauceIngerdient) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{2}
}

func (x *SauceIngerdient) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *SauceIngerdient) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SauceIngerdient) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *SauceIngerdient) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *SauceIngerdient) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

var File_models_food_v1_food_proto protoreflect.FileDescriptor

var file_models_food_v1_food_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x64, 0x79, 0x6e,
	0x61, 0x6d, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x70, 0x67, 0x64, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x67,
	0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x02, 0x0a, 0x05, 0x50, 0x61, 0x73,
	0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02,
	0x02, 0x08, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a,
	0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x3a, 0x24, 0x82, 0xf7, 0x02, 0x1a, 0x12, 0x18,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x02, 0x69, 0x64, 0x1a,
	0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0xd2, 0xf7, 0x02, 0x02, 0x28, 0x01, 0x22, 0xea,
	0x02, 0x0a, 0x0f, 0x50, 0x61, 0x73, 0x74, 0x61, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65,
	0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x21, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x74, 0x61, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01, 0x52,
	0x07, 0x70, 0x61, 0x73, 0x74, 0x61, 0x49, 0x64, 0x3a, 0x41, 0x82, 0xf7, 0x02, 0x39, 0x12, 0x37,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x74, 0x61, 0x5f, 0x69, 0x64, 0x0a, 0x0d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x1a, 0x11, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x69, 0x6e, 0x67,
	0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0xd2, 0xf7, 0x02, 0x00, 0x22, 0x9c, 0x02, 0x0a, 0x0f,
	0x53, 0x61, 0x75, 0x63, 0x65, 0x49, 0x6e, 0x67, 0x65, 0x72, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x3a, 0x23, 0x82, 0xf7, 0x02, 0x1f, 0x12, 0x1d, 0x0a, 0x09, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x02, 0x69, 0x64, 0x1a, 0x0c, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x61, 0x75, 0x63, 0x65, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x75, 0x63, 0x74, 0x6f, 0x6e, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x70, 0x67, 0x64, 0x62,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f,
	0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_models_food_v1_food_proto_rawDescOnce sync.Once
	file_models_food_v1_food_proto_rawDescData = file_models_food_v1_food_proto_rawDesc
)

func file_models_food_v1_food_proto_rawDescGZIP() []byte {
	file_models_food_v1_food_proto_rawDescOnce.Do(func() {
		file_models_food_v1_food_proto_rawDescData = protoimpl.X.CompressGZIP(file_models_food_v1_food_proto_rawDescData)
	})
	return file_models_food_v1_food_proto_rawDescData
}

var file_models_food_v1_food_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_models_food_v1_food_proto_goTypes = []interface{}{
	(*Pasta)(nil),                 // 0: models.food.v1.Pasta
	(*PastaIngredient)(nil),       // 1: models.food.v1.PastaIngredient
	(*SauceIngerdient)(nil),       // 2: models.food.v1.SauceIngerdient
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_models_food_v1_food_proto_depIdxs = []int32{
	3, // 0: models.food.v1.Pasta.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: models.food.v1.Pasta.updated_at:type_name -> google.protobuf.Timestamp
	3, // 2: models.food.v1.Pasta.deleted_at:type_name -> google.protobuf.Timestamp
	3, // 3: models.food.v1.PastaIngredient.created_at:type_name -> google.protobuf.Timestamp
	3, // 4: models.food.v1.PastaIngredient.updated_at:type_name -> google.protobuf.Timestamp
	3, // 5: models.food.v1.PastaIngredient.deleted_at:type_name -> google.protobuf.Timestamp
	3, // 6: models.food.v1.SauceIngerdient.created_at:type_name -> google.protobuf.Timestamp
	3, // 7: models.food.v1.SauceIngerdient.updated_at:type_name -> google.protobuf.Timestamp
	3, // 8: models.food.v1.SauceIngerdient.deleted_at:type_name -> google.protobuf.Timestamp
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_models_food_v1_food_proto_init() }
func file_models_food_v1_food_proto_init() {
	if File_models_food_v1_food_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_models_food_v1_food_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pasta); i {
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
		file_models_food_v1_food_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PastaIngredient); i {
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
		file_models_food_v1_food_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SauceIngerdient); i {
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
			RawDescriptor: file_models_food_v1_food_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_models_food_v1_food_proto_goTypes,
		DependencyIndexes: file_models_food_v1_food_proto_depIdxs,
		MessageInfos:      file_models_food_v1_food_proto_msgTypes,
	}.Build()
	File_models_food_v1_food_proto = out.File
	file_models_food_v1_food_proto_rawDesc = nil
	file_models_food_v1_food_proto_goTypes = nil
	file_models_food_v1_food_proto_depIdxs = nil
}
