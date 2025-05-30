// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: models/llm/v1/models.proto

package v1

import (
	_ "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Map model embedding to their own columns
// for each model embedding, add a column to the table
// xxx_embedding_1
// xxx_embedding_2
// etc for each enum value
// testing
type Model int32

const (
	Model_MODEL_UNSPECIFIED Model = 0
	Model_MODEL_3DIMS       Model = 1
	Model_MODEL_4DIMS       Model = 2
)

// Enum value maps for Model.
var (
	Model_name = map[int32]string{
		0: "MODEL_UNSPECIFIED",
		1: "MODEL_3DIMS",
		2: "MODEL_4DIMS",
	}
	Model_value = map[string]int32{
		"MODEL_UNSPECIFIED": 0,
		"MODEL_3DIMS":       1,
		"MODEL_4DIMS":       2,
	}
)

func (x Model) Enum() *Model {
	p := new(Model)
	*p = x
	return p
}

func (x Model) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Model) Descriptor() protoreflect.EnumDescriptor {
	return file_models_llm_v1_models_proto_enumTypes[0].Descriptor()
}

func (Model) Type() protoreflect.EnumType {
	return &file_models_llm_v1_models_proto_enumTypes[0]
}

func (x Model) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

var File_models_llm_v1_models_proto protoreflect.FileDescriptor

const file_models_llm_v1_models_proto_rawDesc = "" +
	"\n" +
	"\x1amodels/llm/v1/models.proto\x12\rmodels.llm.v1\x1a\x12pgdb/v1/pgdb.proto*P\n" +
	"\x05Model\x12\x15\n" +
	"\x11MODEL_UNSPECIFIED\x10\x00\x12\x17\n" +
	"\vMODEL_3DIMS\x10\x01\x1a\x06\xd2\xf7\x02\x02\b\x03\x12\x17\n" +
	"\vMODEL_4DIMS\x10\x02\x1a\x06\xd2\xf7\x02\x02\b\x04B:Z8github.com/ductone/protoc-gen-pgdb/example/models/llm/v1b\x06proto3"

var file_models_llm_v1_models_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_models_llm_v1_models_proto_goTypes = []any{
	(Model)(0), // 0: models.llm.v1.Model
}
var file_models_llm_v1_models_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_models_llm_v1_models_proto_init() }
func file_models_llm_v1_models_proto_init() {
	if File_models_llm_v1_models_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_models_llm_v1_models_proto_rawDesc), len(file_models_llm_v1_models_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_models_llm_v1_models_proto_goTypes,
		DependencyIndexes: file_models_llm_v1_models_proto_depIdxs,
		EnumInfos:         file_models_llm_v1_models_proto_enumTypes,
	}.Build()
	File_models_llm_v1_models_proto = out.File
	file_models_llm_v1_models_proto_goTypes = nil
	file_models_llm_v1_models_proto_depIdxs = nil
}
