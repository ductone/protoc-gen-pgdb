package pgdb

import (
	"testing"

	"github.com/stretchr/testify/require"

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
