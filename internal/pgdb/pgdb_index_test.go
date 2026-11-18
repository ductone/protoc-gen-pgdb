package pgdb

import (
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
)

func TestValidateLowerColumns(t *testing.T) {
	newIdx := func(b pgdb_v1.MessageOptions_Index_builder) *pgdb_v1.MessageOptions_Index {
		b.Name = "ci"
		return b.Build()
	}

	got := validateLowerColumns("m.v1.Msg", newIdx(pgdb_v1.MessageOptions_Index_builder{
		Method:       pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
		Columns:      []string{"tenant_id", "username"},
		LowerColumns: []string{"username"},
	}))
	require.Equal(t, map[string]struct{}{"username": {}}, got)

	require.Nil(t, validateLowerColumns("m.v1.Msg", newIdx(pgdb_v1.MessageOptions_Index_builder{
		Method:  pgdb_v1.MessageOptions_Index_INDEX_METHOD_GIN,
		Columns: []string{"username"},
	})), "no lower_columns leaves other methods untouched")

	tests := []struct {
		name    string
		idx     pgdb_v1.MessageOptions_Index_builder
		wantErr string
	}{
		{
			name: "not in columns",
			idx: pgdb_v1.MessageOptions_Index_builder{
				Method:       pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:      []string{"tenant_id"},
				LowerColumns: []string{"username"},
			},
			wantErr: "pgdb: index 'ci' on 'm.v1.Msg': lower_columns entry 'username' is not listed in columns",
		},
		{
			name: "non-btree method",
			idx: pgdb_v1.MessageOptions_Index_builder{
				Method:       pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
				Columns:      []string{"username"},
				LowerColumns: []string{"username"},
			},
			wantErr: "pgdb: index 'ci' on 'm.v1.Msg': lower_columns requires method INDEX_METHOD_BTREE, got INDEX_METHOD_BTREE_GIN",
		},
		{
			name: "bit_hamming_ops",
			idx: pgdb_v1.MessageOptions_Index_builder{
				Method:        pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:       []string{"username"},
				LowerColumns:  []string{"username"},
				BitHammingOps: true,
			},
			wantErr: "pgdb: index 'ci' on 'm.v1.Msg': lower_columns cannot be combined with bit_hamming_ops",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.PanicsWithError(t, tt.wantErr, func() {
				validateLowerColumns("m.v1.Msg", newIdx(tt.idx))
			})
		})
	}
}

// lowerColumnsTestMessage builds a pgs AST for a message exercising every
// lower_columns field-resolution case, without running the plugin.
func lowerColumnsTestMessage(t *testing.T) pgs.Message {
	t.Helper()
	str := descriptorpb.FieldDescriptorProto_TYPE_STRING
	optional := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	field := func(name string, num int32, typ descriptorpb.FieldDescriptorProto_Type) *descriptorpb.FieldDescriptorProto {
		return &descriptorpb.FieldDescriptorProto{
			Name:     proto.String(name),
			JsonName: proto.String(name),
			Number:   proto.Int32(num),
			Label:    optional.Enum(),
			Type:     typ.Enum(),
		}
	}
	msgField := func(name string, num int32, typeName string) *descriptorpb.FieldDescriptorProto {
		f := field(name, num, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE)
		f.TypeName = proto.String(typeName)
		return f
	}

	tags := field("tags", 4, str)
	tags.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()

	addr := field("addr", 5, str)
	addr.Options = &descriptorpb.FieldOptions{}
	proto.SetExtension(addr.Options, pgdb_v1.E_Options, pgdb_v1.FieldOptions_builder{
		MessageBehavior: pgdb_v1.FieldOptions_MESSAGE_BEHAVIOR_INET_ADDR,
	}.Build())

	oneofA := field("a", 6, str)
	oneofA.OneofIndex = proto.Int32(0)
	oneofB := field("b", 7, str)
	oneofB.OneofIndex = proto.Int32(0)

	file := &descriptorpb.FileDescriptorProto{
		Name:       proto.String("t/v1/t.proto"),
		Package:    proto.String("t.v1"),
		Syntax:     proto.String("proto3"),
		Dependency: []string{"google/protobuf/wrappers.proto"},
		Options:    &descriptorpb.FileOptions{GoPackage: proto.String("example.com/t/v1")},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name:  proto.String("Profile"),
				Field: []*descriptorpb.FieldDescriptorProto{field("handle", 1, str)},
			},
			{
				Name: proto.String("Msg"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("tenant_id", 1, str),
					field("username", 2, str),
					field("count", 3, descriptorpb.FieldDescriptorProto_TYPE_INT64),
					tags,
					addr,
					oneofA,
					oneofB,
					msgField("profile", 8, ".t.v1.Profile"),
					msgField("wrapped", 9, ".google.protobuf.StringValue"),
				},
				OneofDecl: []*descriptorpb.OneofDescriptorProto{{Name: proto.String("choice")}},
			},
		},
	}
	fdset := &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{
		protodesc.ToFileDescriptorProto(wrapperspb.File_google_protobuf_wrappers_proto),
		file,
	}}
	ast := pgs.ProcessFileDescriptorSet(pgs.InitMockDebugger(), fdset)
	for _, f := range ast.Packages()["t.v1"].Files() {
		for _, m := range f.Messages() {
			if m.Name().String() == "Msg" {
				return m
			}
		}
	}
	require.FailNow(t, "message t.v1.Msg not found")
	return nil
}

func TestExtraIndexesLowerColumns(t *testing.T) {
	m := lowerColumnsTestMessage(t)
	module := &Module{}
	idx := func(columns ...string) *pgdb_v1.MessageOptions_Index {
		return pgdb_v1.MessageOptions_Index_builder{
			Name:         "ci",
			Method:       pgdb_v1.MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:      append([]string{"tenant_id"}, columns...),
			LowerColumns: columns,
		}.Build()
	}

	t.Run("top-level string", func(t *testing.T) {
		got := module.extraIndexes(nil, m, nil, idx("username"))
		require.Equal(t, []string{"tenant_id", "username"}, got.DB.Columns)
		require.Equal(t, `" + io.ColumnName("tenant_id") + ", lower(" + io.ColumnName("username") + ")`, got.DB.OverrideExpression)
		require.Equal(t, map[string]struct{}{"username": {}}, got.LoweredColumns)
	})

	t.Run("nested dotted column", func(t *testing.T) {
		got := module.extraIndexes(nil, m, nil, idx("profile.handle"))
		require.Equal(t, []string{"tenant_id", "8$handle"}, got.DB.Columns)
		require.Equal(t, `" + io.ColumnName("tenant_id") + ", lower(" + io.ColumnName("8$handle") + ")`, got.DB.OverrideExpression)
		require.Equal(t, map[string]struct{}{"8$handle": {}}, got.LoweredColumns)
	})

	failures := []struct {
		column  string
		wantErr string
	}{
		{"choice", "pgdb: index 'ci' on '.t.v1.Msg': lower_columns entry 'choice' must be a string field, not a oneof"},
		{"count", "pgdb: index 'ci' on '.t.v1.Msg': lower_columns field 'count' must be a string field, got TYPE_INT64"},
		{"tags", "pgdb: index 'ci' on '.t.v1.Msg': lower_columns field 'tags' must not be repeated"},
		{"addr", "pgdb: index 'ci' on '.t.v1.Msg': lower_columns field 'addr' must be stored as text, but has message_behavior MESSAGE_BEHAVIOR_INET_ADDR"},
		{"wrapped", "pgdb: index 'ci' on '.t.v1.Msg': lower_columns field 'wrapped' must be a string field, got TYPE_MESSAGE"},
	}
	for _, tt := range failures {
		t.Run("rejects "+tt.column, func(t *testing.T) {
			require.PanicsWithError(t, tt.wantErr, func() {
				module.extraIndexes(nil, m, nil, idx(tt.column))
			})
		})
	}
}
