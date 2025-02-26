// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: models/food/v1/food.proto

package v1

import (
	v1 "github.com/ductone/protoc-gen-pgdb/example/models/llm/v1"
	_ "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	_ "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Pasta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TenantId      string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pasta) Reset() {
	*x = Pasta{}
	mi := &file_models_food_v1_food_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pasta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pasta) ProtoMessage() {}

func (x *Pasta) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[0]
	if x != nil {
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
	state        protoimpl.MessageState `protogen:"open.v1"`
	TenantId     string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	IngredientId string                 `protobuf:"bytes,2,opt,name=ingredient_id,json=ingredientId,proto3" json:"ingredient_id,omitempty"`
	CreatedAt    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	PastaId      string                 `protobuf:"bytes,6,opt,name=pasta_id,json=pastaId,proto3" json:"pasta_id,omitempty"`
	Id           string                 `protobuf:"bytes,7,opt,name=id,proto3" json:"id,omitempty"`
	// We want to verify
	// 1. Type is repeated
	// 2. Type is a nested message
	// 3. Message type has 2 fields, enum type and repeated float type
	ModelEmbeddings []*PastaIngredient_ModelEmbedding `protobuf:"bytes,8,rep,name=model_embeddings,json=modelEmbeddings,proto3" json:"model_embeddings,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *PastaIngredient) Reset() {
	*x = PastaIngredient{}
	mi := &file_models_food_v1_food_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PastaIngredient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastaIngredient) ProtoMessage() {}

func (x *PastaIngredient) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[1]
	if x != nil {
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

func (x *PastaIngredient) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PastaIngredient) GetModelEmbeddings() []*PastaIngredient_ModelEmbedding {
	if x != nil {
		return x.ModelEmbeddings
	}
	return nil
}

type SauceIngredient struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TenantId      string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	SourceAddr    string                 `protobuf:"bytes,6,opt,name=source_addr,json=sourceAddr,proto3" json:"source_addr,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SauceIngredient) Reset() {
	*x = SauceIngredient{}
	mi := &file_models_food_v1_food_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SauceIngredient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SauceIngredient) ProtoMessage() {}

func (x *SauceIngredient) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SauceIngredient.ProtoReflect.Descriptor instead.
func (*SauceIngredient) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{2}
}

func (x *SauceIngredient) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *SauceIngredient) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SauceIngredient) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *SauceIngredient) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *SauceIngredient) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *SauceIngredient) GetSourceAddr() string {
	if x != nil {
		return x.SourceAddr
	}
	return ""
}

type CheeseIngredient struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TenantId      string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	SourceAddr    string                 `protobuf:"bytes,6,opt,name=source_addr,json=sourceAddr,proto3" json:"source_addr,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheeseIngredient) Reset() {
	*x = CheeseIngredient{}
	mi := &file_models_food_v1_food_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheeseIngredient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheeseIngredient) ProtoMessage() {}

func (x *CheeseIngredient) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheeseIngredient.ProtoReflect.Descriptor instead.
func (*CheeseIngredient) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{3}
}

func (x *CheeseIngredient) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *CheeseIngredient) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CheeseIngredient) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *CheeseIngredient) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *CheeseIngredient) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *CheeseIngredient) GetSourceAddr() string {
	if x != nil {
		return x.SourceAddr
	}
	return ""
}

type PastaIngredient_ModelEmbedding struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Model         v1.Model               `protobuf:"varint,1,opt,name=model,proto3,enum=models.llm.v1.Model" json:"model,omitempty"`
	Embedding     []float32              `protobuf:"fixed32,2,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PastaIngredient_ModelEmbedding) Reset() {
	*x = PastaIngredient_ModelEmbedding{}
	mi := &file_models_food_v1_food_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PastaIngredient_ModelEmbedding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastaIngredient_ModelEmbedding) ProtoMessage() {}

func (x *PastaIngredient_ModelEmbedding) ProtoReflect() protoreflect.Message {
	mi := &file_models_food_v1_food_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastaIngredient_ModelEmbedding.ProtoReflect.Descriptor instead.
func (*PastaIngredient_ModelEmbedding) Descriptor() ([]byte, []int) {
	return file_models_food_v1_food_proto_rawDescGZIP(), []int{1, 0}
}

func (x *PastaIngredient_ModelEmbedding) GetModel() v1.Model {
	if x != nil {
		return x.Model
	}
	return v1.Model(0)
}

func (x *PastaIngredient_ModelEmbedding) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

var File_models_food_v1_food_proto protoreflect.FileDescriptor

var file_models_food_v1_food_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x64, 0x79, 0x6e,
	0x61, 0x6d, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x6c, 0x6c, 0x6d,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x12, 0x70, 0x67, 0x64, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x02, 0x0a, 0x05, 0x50, 0x61, 0x73, 0x74, 0x61, 0x12, 0x1b,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x3a, 0x24, 0x82, 0xf7, 0x02, 0x1a, 0x12, 0x18, 0x0a, 0x09, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x02, 0x69, 0x64, 0x1a, 0x07, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0xd2, 0xf7, 0x02, 0x02, 0x28, 0x01, 0x22, 0xc7, 0x06, 0x0a, 0x0f, 0x50,
	0x61, 0x73, 0x74, 0x61, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x69,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x21, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x74, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01, 0x52, 0x07, 0x70, 0x61, 0x73,
	0x74, 0x61, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x61, 0x0a, 0x10,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x61, 0x49, 0x6e, 0x67,
	0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x45, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x18, 0x04, 0x52, 0x0f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x1a,
	0x62, 0x0a, 0x0e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x2a, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x14, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x6c, 0x6c, 0x6d, 0x2e, 0x76, 0x31,
	0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1c, 0x0a,
	0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x02,
	0x52, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x06, 0xd2, 0xf7, 0x02,
	0x02, 0x20, 0x01, 0x3a, 0xbe, 0x02, 0x82, 0xf7, 0x02, 0x3d, 0x12, 0x3b, 0x0a, 0x09, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x74, 0x61, 0x5f, 0x69,
	0x64, 0x0a, 0x0d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x0a, 0x02, 0x69, 0x64, 0x1a, 0x11, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x69, 0x6e, 0x67,
	0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0xd2, 0xf7, 0x02, 0xf8, 0x01, 0x12, 0x23, 0x0a, 0x06,
	0x70, 0x61, 0x73, 0x74, 0x61, 0x73, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x5f, 0x69, 0x64, 0x1a, 0x08, 0x70, 0x61, 0x73, 0x74, 0x61, 0x5f, 0x69, 0x64, 0x1a, 0x02, 0x69,
	0x64, 0x12, 0x2d, 0x0a, 0x0b, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0d, 0x69,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x02, 0x69, 0x64,
	0x12, 0x2e, 0x0a, 0x0d, 0x65, 0x76, 0x65, 0x72, 0x79, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x67, 0x67,
	0x67, 0x10, 0x01, 0x1a, 0x0d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x1a, 0x08, 0x70, 0x61, 0x73, 0x74, 0x61, 0x5f, 0x69, 0x64, 0x1a, 0x02, 0x69, 0x64,
	0x12, 0x3b, 0x0a, 0x18, 0x65, 0x76, 0x65, 0x72, 0x79, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x67, 0x67,
	0x67, 0x5f, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x10, 0x01, 0x1a, 0x0d,
	0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x08, 0x70,
	0x61, 0x73, 0x74, 0x61, 0x5f, 0x69, 0x64, 0x1a, 0x02, 0x69, 0x64, 0x28, 0x01, 0x28, 0x01, 0x32,
	0x33, 0x0a, 0x17, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x5f, 0x69, 0x6e, 0x67,
	0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x22, 0xf8, 0x02, 0x0a, 0x0f, 0x53, 0x61, 0x75, 0x63, 0x65, 0x49, 0x6e,
	0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x08, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27,
	0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02, 0x02, 0x18, 0x05, 0x52, 0x0a, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x41, 0x64, 0x64, 0x72, 0x3a, 0x56, 0x82, 0xf7, 0x02, 0x1f, 0x12, 0x1d, 0x0a,
	0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x02, 0x69, 0x64, 0x1a, 0x0c,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x61, 0x75, 0x63, 0x65, 0xd2, 0xf7, 0x02, 0x2f,
	0x12, 0x2d, 0x0a, 0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x5f,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x1a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x22,
	0xfe, 0x02, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x65, 0x73, 0x65, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64,
	0x69, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2,
	0xf7, 0x02, 0x02, 0x08, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27, 0x0a, 0x0b, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x06, 0xd2, 0xf7, 0x02, 0x02, 0x18, 0x05, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x41,
	0x64, 0x64, 0x72, 0x3a, 0x5b, 0x82, 0xf7, 0x02, 0x20, 0x12, 0x1e, 0x0a, 0x09, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x0a, 0x02, 0x69, 0x64, 0x1a, 0x0d, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x63, 0x68, 0x65, 0x65, 0x73, 0x65, 0xd2, 0xf7, 0x02, 0x33, 0x12, 0x2d, 0x0a,
	0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a,
	0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x38, 0x01, 0x40, 0x02,
	0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64,
	0x75, 0x63, 0x74, 0x6f, 0x6e, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65,
	0x6e, 0x2d, 0x70, 0x67, 0x64, 0x62, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_models_food_v1_food_proto_rawDescOnce sync.Once
	file_models_food_v1_food_proto_rawDescData []byte
)

func file_models_food_v1_food_proto_rawDescGZIP() []byte {
	file_models_food_v1_food_proto_rawDescOnce.Do(func() {
		file_models_food_v1_food_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_models_food_v1_food_proto_rawDesc), len(file_models_food_v1_food_proto_rawDesc)))
	})
	return file_models_food_v1_food_proto_rawDescData
}

var file_models_food_v1_food_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_models_food_v1_food_proto_goTypes = []any{
	(*Pasta)(nil),                          // 0: models.food.v1.Pasta
	(*PastaIngredient)(nil),                // 1: models.food.v1.PastaIngredient
	(*SauceIngredient)(nil),                // 2: models.food.v1.SauceIngredient
	(*CheeseIngredient)(nil),               // 3: models.food.v1.CheeseIngredient
	(*PastaIngredient_ModelEmbedding)(nil), // 4: models.food.v1.PastaIngredient.ModelEmbedding
	(*timestamppb.Timestamp)(nil),          // 5: google.protobuf.Timestamp
	(v1.Model)(0),                          // 6: models.llm.v1.Model
}
var file_models_food_v1_food_proto_depIdxs = []int32{
	5,  // 0: models.food.v1.Pasta.created_at:type_name -> google.protobuf.Timestamp
	5,  // 1: models.food.v1.Pasta.updated_at:type_name -> google.protobuf.Timestamp
	5,  // 2: models.food.v1.Pasta.deleted_at:type_name -> google.protobuf.Timestamp
	5,  // 3: models.food.v1.PastaIngredient.created_at:type_name -> google.protobuf.Timestamp
	5,  // 4: models.food.v1.PastaIngredient.updated_at:type_name -> google.protobuf.Timestamp
	5,  // 5: models.food.v1.PastaIngredient.deleted_at:type_name -> google.protobuf.Timestamp
	4,  // 6: models.food.v1.PastaIngredient.model_embeddings:type_name -> models.food.v1.PastaIngredient.ModelEmbedding
	5,  // 7: models.food.v1.SauceIngredient.created_at:type_name -> google.protobuf.Timestamp
	5,  // 8: models.food.v1.SauceIngredient.updated_at:type_name -> google.protobuf.Timestamp
	5,  // 9: models.food.v1.SauceIngredient.deleted_at:type_name -> google.protobuf.Timestamp
	5,  // 10: models.food.v1.CheeseIngredient.created_at:type_name -> google.protobuf.Timestamp
	5,  // 11: models.food.v1.CheeseIngredient.updated_at:type_name -> google.protobuf.Timestamp
	5,  // 12: models.food.v1.CheeseIngredient.deleted_at:type_name -> google.protobuf.Timestamp
	6,  // 13: models.food.v1.PastaIngredient.ModelEmbedding.model:type_name -> models.llm.v1.Model
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_models_food_v1_food_proto_init() }
func file_models_food_v1_food_proto_init() {
	if File_models_food_v1_food_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_models_food_v1_food_proto_rawDesc), len(file_models_food_v1_food_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_models_food_v1_food_proto_goTypes,
		DependencyIndexes: file_models_food_v1_food_proto_depIdxs,
		MessageInfos:      file_models_food_v1_food_proto_msgTypes,
	}.Build()
	File_models_food_v1_food_proto = out.File
	file_models_food_v1_food_proto_goTypes = nil
	file_models_food_v1_food_proto_depIdxs = nil
}
