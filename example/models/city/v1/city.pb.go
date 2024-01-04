// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: models/city/v1/city.proto

package v1

import (
	v11 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"
	v1 "github.com/ductone/protoc-gen-pgdb/example/models/zoo/v1"
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

type Attractions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId  string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Id        string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Numid     int32                  `protobuf:"varint,3,opt,name=numid,proto3" json:"numid,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Types that are assignable to What:
	//	*Attractions_Pet
	//	*Attractions_ZooShop
	What   isAttractions_What `protobuf_oneof:"what"`
	Medium *v1.Shop           `protobuf:"bytes,12,opt,name=medium,proto3" json:"medium,omitempty"`
}

func (x *Attractions) Reset() {
	*x = Attractions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_city_v1_city_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attractions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attractions) ProtoMessage() {}

func (x *Attractions) ProtoReflect() protoreflect.Message {
	mi := &file_models_city_v1_city_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attractions.ProtoReflect.Descriptor instead.
func (*Attractions) Descriptor() ([]byte, []int) {
	return file_models_city_v1_city_proto_rawDescGZIP(), []int{0}
}

func (x *Attractions) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *Attractions) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Attractions) GetNumid() int32 {
	if x != nil {
		return x.Numid
	}
	return 0
}

func (x *Attractions) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (m *Attractions) GetWhat() isAttractions_What {
	if m != nil {
		return m.What
	}
	return nil
}

func (x *Attractions) GetPet() *v11.Pet {
	if x, ok := x.GetWhat().(*Attractions_Pet); ok {
		return x.Pet
	}
	return nil
}

func (x *Attractions) GetZooShop() *v1.Shop {
	if x, ok := x.GetWhat().(*Attractions_ZooShop); ok {
		return x.ZooShop
	}
	return nil
}

func (x *Attractions) GetMedium() *v1.Shop {
	if x != nil {
		return x.Medium
	}
	return nil
}

type isAttractions_What interface {
	isAttractions_What()
}

type Attractions_Pet struct {
	Pet *v11.Pet `protobuf:"bytes,10,opt,name=pet,proto3,oneof"`
}

type Attractions_ZooShop struct {
	ZooShop *v1.Shop `protobuf:"bytes,11,opt,name=zoo_shop,json=zooShop,proto3,oneof"`
}

func (*Attractions_Pet) isAttractions_What() {}

func (*Attractions_ZooShop) isAttractions_What() {}

var File_models_city_v1_city_proto protoreflect.FileDescriptor

var file_models_city_v1_city_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x69, 0x74, 0x79, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x64, 0x79, 0x6e,
	0x61, 0x6d, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x61, 0x6e, 0x69,
	0x6d, 0x61, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x7a, 0x6f,
	0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x7a, 0x6f, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12,
	0x70, 0x67, 0x64, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x67, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xe5, 0x04, 0x0a, 0x0b, 0x41, 0x74, 0x74, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xd2, 0xf7, 0x02,
	0x02, 0x08, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x75, 0x6d, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6e, 0x75, 0x6d, 0x69, 0x64, 0x12, 0x39, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2a, 0x0a, 0x03, 0x70, 0x65, 0x74, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x61,
	0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x74, 0x48, 0x00, 0x52,
	0x03, 0x70, 0x65, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x7a, 0x6f, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x70,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x7a, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x6f, 0x70, 0x48, 0x00, 0x52, 0x07, 0x7a,
	0x6f, 0x6f, 0x53, 0x68, 0x6f, 0x70, 0x12, 0x2b, 0x0a, 0x06, 0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x7a, 0x6f, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x6f, 0x70, 0x52, 0x06, 0x6d, 0x65, 0x64,
	0x69, 0x75, 0x6d, 0x3a, 0xbc, 0x02, 0x82, 0xf7, 0x02, 0x18, 0x12, 0x16, 0x0a, 0x09, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x12, 0x02, 0x69, 0x64, 0x12, 0x05, 0x6e, 0x75, 0x6d,
	0x69, 0x64, 0xd2, 0xf7, 0x02, 0x9b, 0x02, 0x12, 0x23, 0x0a, 0x06, 0x66, 0x75, 0x72, 0x72, 0x72,
	0x73, 0x10, 0x03, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0c,
	0x7a, 0x6f, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x66, 0x75, 0x72, 0x12, 0x3e, 0x0a, 0x12,
	0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x6e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x10, 0x03, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a,
	0x1b, 0x7a, 0x6f, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x61, 0x6e, 0x79, 0x74, 0x68, 0x69,
	0x6e, 0x67, 0x2e, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x5f, 0x36, 0x34, 0x12, 0x1a, 0x0a, 0x05,
	0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x1a, 0x04, 0x77, 0x68, 0x61, 0x74, 0x12, 0x2c, 0x0a, 0x0c, 0x6e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0f, 0x7a, 0x6f, 0x6f, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e,
	0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x12, 0x2b, 0x0a, 0x0d, 0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d,
	0x5f, 0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0d, 0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x75, 0x6d, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x65, 0x74, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x10, 0x02, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x0b,
	0x70, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x15, 0x0a, 0x02, 0x69,
	0x64, 0x10, 0x01, 0x1a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x1a, 0x02,
	0x69, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x77, 0x68, 0x61, 0x74, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x75, 0x63, 0x74, 0x6f, 0x6e, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x70, 0x67, 0x64, 0x62,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f,
	0x63, 0x69, 0x74, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_models_city_v1_city_proto_rawDescOnce sync.Once
	file_models_city_v1_city_proto_rawDescData = file_models_city_v1_city_proto_rawDesc
)

func file_models_city_v1_city_proto_rawDescGZIP() []byte {
	file_models_city_v1_city_proto_rawDescOnce.Do(func() {
		file_models_city_v1_city_proto_rawDescData = protoimpl.X.CompressGZIP(file_models_city_v1_city_proto_rawDescData)
	})
	return file_models_city_v1_city_proto_rawDescData
}

var file_models_city_v1_city_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_models_city_v1_city_proto_goTypes = []interface{}{
	(*Attractions)(nil),           // 0: models.city.v1.Attractions
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*v11.Pet)(nil),               // 2: models.animals.v1.Pet
	(*v1.Shop)(nil),               // 3: models.zoo.v1.Shop
}
var file_models_city_v1_city_proto_depIdxs = []int32{
	1, // 0: models.city.v1.Attractions.created_at:type_name -> google.protobuf.Timestamp
	2, // 1: models.city.v1.Attractions.pet:type_name -> models.animals.v1.Pet
	3, // 2: models.city.v1.Attractions.zoo_shop:type_name -> models.zoo.v1.Shop
	3, // 3: models.city.v1.Attractions.medium:type_name -> models.zoo.v1.Shop
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_models_city_v1_city_proto_init() }
func file_models_city_v1_city_proto_init() {
	if File_models_city_v1_city_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_models_city_v1_city_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attractions); i {
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
	file_models_city_v1_city_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Attractions_Pet)(nil),
		(*Attractions_ZooShop)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_models_city_v1_city_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_models_city_v1_city_proto_goTypes,
		DependencyIndexes: file_models_city_v1_city_proto_depIdxs,
		MessageInfos:      file_models_city_v1_city_proto_msgTypes,
	}.Build()
	File_models_city_v1_city_proto = out.File
	file_models_city_v1_city_proto_rawDesc = nil
	file_models_city_v1_city_proto_goTypes = nil
	file_models_city_v1_city_proto_depIdxs = nil
}
