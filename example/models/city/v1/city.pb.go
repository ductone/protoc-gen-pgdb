// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
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
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Attractions struct {
	state                protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_TenantId  string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3"`
	xxx_hidden_Id        string                 `protobuf:"bytes,2,opt,name=id,proto3"`
	xxx_hidden_Numid     int32                  `protobuf:"varint,3,opt,name=numid,proto3"`
	xxx_hidden_CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3"`
	xxx_hidden_What      isAttractions_What     `protobuf_oneof:"what"`
	xxx_hidden_Medium    *v1.Shop               `protobuf:"bytes,12,opt,name=medium,proto3"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *Attractions) Reset() {
	*x = Attractions{}
	mi := &file_models_city_v1_city_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Attractions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attractions) ProtoMessage() {}

func (x *Attractions) ProtoReflect() protoreflect.Message {
	mi := &file_models_city_v1_city_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Attractions) GetTenantId() string {
	if x != nil {
		return x.xxx_hidden_TenantId
	}
	return ""
}

func (x *Attractions) GetId() string {
	if x != nil {
		return x.xxx_hidden_Id
	}
	return ""
}

func (x *Attractions) GetNumid() int32 {
	if x != nil {
		return x.xxx_hidden_Numid
	}
	return 0
}

func (x *Attractions) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.xxx_hidden_CreatedAt
	}
	return nil
}

func (x *Attractions) GetPet() *v11.Pet {
	if x != nil {
		if x, ok := x.xxx_hidden_What.(*attractions_Pet); ok {
			return x.Pet
		}
	}
	return nil
}

func (x *Attractions) GetZooShop() *v1.Shop {
	if x != nil {
		if x, ok := x.xxx_hidden_What.(*attractions_ZooShop); ok {
			return x.ZooShop
		}
	}
	return nil
}

func (x *Attractions) GetMedium() *v1.Shop {
	if x != nil {
		return x.xxx_hidden_Medium
	}
	return nil
}

func (x *Attractions) SetTenantId(v string) {
	x.xxx_hidden_TenantId = v
}

func (x *Attractions) SetId(v string) {
	x.xxx_hidden_Id = v
}

func (x *Attractions) SetNumid(v int32) {
	x.xxx_hidden_Numid = v
}

func (x *Attractions) SetCreatedAt(v *timestamppb.Timestamp) {
	x.xxx_hidden_CreatedAt = v
}

func (x *Attractions) SetPet(v *v11.Pet) {
	if v == nil {
		x.xxx_hidden_What = nil
		return
	}
	x.xxx_hidden_What = &attractions_Pet{v}
}

func (x *Attractions) SetZooShop(v *v1.Shop) {
	if v == nil {
		x.xxx_hidden_What = nil
		return
	}
	x.xxx_hidden_What = &attractions_ZooShop{v}
}

func (x *Attractions) SetMedium(v *v1.Shop) {
	x.xxx_hidden_Medium = v
}

func (x *Attractions) HasCreatedAt() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_CreatedAt != nil
}

func (x *Attractions) HasWhat() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_What != nil
}

func (x *Attractions) HasPet() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_What.(*attractions_Pet)
	return ok
}

func (x *Attractions) HasZooShop() bool {
	if x == nil {
		return false
	}
	_, ok := x.xxx_hidden_What.(*attractions_ZooShop)
	return ok
}

func (x *Attractions) HasMedium() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Medium != nil
}

func (x *Attractions) ClearCreatedAt() {
	x.xxx_hidden_CreatedAt = nil
}

func (x *Attractions) ClearWhat() {
	x.xxx_hidden_What = nil
}

func (x *Attractions) ClearPet() {
	if _, ok := x.xxx_hidden_What.(*attractions_Pet); ok {
		x.xxx_hidden_What = nil
	}
}

func (x *Attractions) ClearZooShop() {
	if _, ok := x.xxx_hidden_What.(*attractions_ZooShop); ok {
		x.xxx_hidden_What = nil
	}
}

func (x *Attractions) ClearMedium() {
	x.xxx_hidden_Medium = nil
}

const Attractions_What_not_set_case case_Attractions_What = 0
const Attractions_Pet_case case_Attractions_What = 10
const Attractions_ZooShop_case case_Attractions_What = 11

func (x *Attractions) WhichWhat() case_Attractions_What {
	if x == nil {
		return Attractions_What_not_set_case
	}
	switch x.xxx_hidden_What.(type) {
	case *attractions_Pet:
		return Attractions_Pet_case
	case *attractions_ZooShop:
		return Attractions_ZooShop_case
	default:
		return Attractions_What_not_set_case
	}
}

type Attractions_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	TenantId  string
	Id        string
	Numid     int32
	CreatedAt *timestamppb.Timestamp
	// Fields of oneof xxx_hidden_What:
	Pet     *v11.Pet
	ZooShop *v1.Shop
	// -- end of xxx_hidden_What
	Medium *v1.Shop
}

func (b0 Attractions_builder) Build() *Attractions {
	m0 := &Attractions{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_TenantId = b.TenantId
	x.xxx_hidden_Id = b.Id
	x.xxx_hidden_Numid = b.Numid
	x.xxx_hidden_CreatedAt = b.CreatedAt
	if b.Pet != nil {
		x.xxx_hidden_What = &attractions_Pet{b.Pet}
	}
	if b.ZooShop != nil {
		x.xxx_hidden_What = &attractions_ZooShop{b.ZooShop}
	}
	x.xxx_hidden_Medium = b.Medium
	return m0
}

type case_Attractions_What protoreflect.FieldNumber

func (x case_Attractions_What) String() string {
	md := file_models_city_v1_city_proto_msgTypes[0].Descriptor()
	if x == 0 {
		return "not set"
	}
	return protoimpl.X.MessageFieldStringOf(md, protoreflect.FieldNumber(x))
}

type isAttractions_What interface {
	isAttractions_What()
}

type attractions_Pet struct {
	Pet *v11.Pet `protobuf:"bytes,10,opt,name=pet,proto3,oneof"`
}

type attractions_ZooShop struct {
	ZooShop *v1.Shop `protobuf:"bytes,11,opt,name=zoo_shop,json=zooShop,proto3,oneof"`
}

func (*attractions_Pet) isAttractions_What() {}

func (*attractions_ZooShop) isAttractions_What() {}

var File_models_city_v1_city_proto protoreflect.FileDescriptor

const file_models_city_v1_city_proto_rawDesc = "" +
	"\n" +
	"\x19models/city/v1/city.proto\x12\x0emodels.city.v1\x1a\x16dynamo/v1/dynamo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1fmodels/animals/v1/animals.proto\x1a\x17models/zoo/v1/zoo.proto\x1a\x12pgdb/v1/pgdb.proto\"\xe5\x04\n" +
	"\vAttractions\x12\x1b\n" +
	"\ttenant_id\x18\x01 \x01(\tR\btenantId\x12\x16\n" +
	"\x02id\x18\x02 \x01(\tB\x06\xd2\xf7\x02\x02\b\x01R\x02id\x12\x14\n" +
	"\x05numid\x18\x03 \x01(\x05R\x05numid\x129\n" +
	"\n" +
	"created_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x12*\n" +
	"\x03pet\x18\n" +
	" \x01(\v2\x16.models.animals.v1.PetH\x00R\x03pet\x120\n" +
	"\bzoo_shop\x18\v \x01(\v2\x13.models.zoo.v1.ShopH\x00R\azooShop\x12+\n" +
	"\x06medium\x18\f \x01(\v2\x13.models.zoo.v1.ShopR\x06medium:\xbc\x02\x82\xf7\x02\x18\x12\x16\n" +
	"\ttenant_id\x12\x02id\x12\x05numid\xd2\xf7\x02\x9b\x02\x12#\n" +
	"\x06furrrs\x10\x03\x1a\ttenant_id\x1a\fzoo_shop.fur\x12>\n" +
	"\x12nestednestednested\x10\x03\x1a\ttenant_id\x1a\x1bzoo_shop.anything.sfixed_64\x12\x1a\n" +
	"\x05oneof\x10\x01\x1a\ttenant_id\x1a\x04what\x12,\n" +
	"\fnested_oneof\x10\x01\x1a\ttenant_id\x1a\x0fzoo_shop.medium\x12+\n" +
	"\rmedium_medium\x10\x01\x1a\ttenant_id\x1a\rmedium.medium\x12&\n" +
	"\n" +
	"petprofile\x10\x02\x1a\ttenant_id\x1a\vpet.profile\x12\x15\n" +
	"\x02id\x10\x01\x1a\ttenant_id\x1a\x02idB\x06\n" +
	"\x04whatB;Z9github.com/ductone/protoc-gen-pgdb/example/models/city/v1b\x06proto3"

var file_models_city_v1_city_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_models_city_v1_city_proto_goTypes = []any{
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
	file_models_city_v1_city_proto_msgTypes[0].OneofWrappers = []any{
		(*attractions_Pet)(nil),
		(*attractions_ZooShop)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_models_city_v1_city_proto_rawDesc), len(file_models_city_v1_city_proto_rawDesc)),
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
	file_models_city_v1_city_proto_goTypes = nil
	file_models_city_v1_city_proto_depIdxs = nil
}
