package v1

import (
	"bytes"
	"testing"
)

type mockDescriptor struct {
	tableName                   string
	isPartitioned               bool
	isPartitionedByCreatedAt    bool
	partitionedByKsuidFieldName string
	partitionDateRange          MessageOptions_PartitionedByDateRange
}

func (m *mockDescriptor) TableName() string {
	return m.tableName
}

func (m *mockDescriptor) Fields(opts ...DescriptorFieldOptionFunc) []*Column {
	return nil
}

func (m *mockDescriptor) PKSKField() *Column {
	return &Column{
		Name: "pb$pksk",
	}
}

func (m *mockDescriptor) DataField() *Column {
	return nil
}

func (m *mockDescriptor) SearchField() *Column {
	return nil
}

func (m *mockDescriptor) VersioningField() *Column {
	return &Column{
		Name: "pb$updated_at",
	}
}

func (m *mockDescriptor) TenantField() *Column {
	return nil
}

func (m *mockDescriptor) IsPartitioned() bool {
	return m.isPartitioned
}

func (m *mockDescriptor) IsPartitionedByCreatedAt() bool {
	return m.isPartitionedByCreatedAt
}

func (m *mockDescriptor) GetPartitionedByKsuidFieldName() string {
	return m.partitionedByKsuidFieldName
}

func (m *mockDescriptor) Indexes(opts ...IndexOptionsFunc) []*Index {
	return nil
}

func (m *mockDescriptor) IndexPrimaryKey(opts ...IndexOptionsFunc) *Index {
	return &Index{
		Name: "pbidx_" + m.tableName,
	}
}

func (m *mockDescriptor) Statistics(opts ...StatisticOptionsFunc) []*Statistic {
	return nil
}

func (m *mockDescriptor) GetPartitionDateRange() MessageOptions_PartitionedByDateRange {
	return m.partitionDateRange
}

func (m *mockDescriptor) GetStorageParameters() *MessageOptions_StorageParameters {
	return nil
}

func TestPgWriteString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "test",
			expected: `"test"`,
		},
		{
			name:     "String with spaces",
			input:    "test table",
			expected: `"test table"`,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: `""`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			pgWriteString(buf, test.input)
			if buf.String() != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, buf.String())
			}
		})
	}
}
func TestCol2alter(t *testing.T) {
	tests := []struct {
		name     string
		desc     Descriptor
		col      *Column
		expected string
	}{
		{
			name: "Add simple column",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			col: &Column{
				Name:     "col1",
				Type:     "TEXT",
				Nullable: true,
			},
			expected: `ALTER TABLE "test_table"
ADD COLUMN IF NOT EXISTS
  "col1" TEXT`,
		},
		{
			name: "Add NOT NULL column",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			col: &Column{
				Name:     "col2",
				Type:     "INTEGER",
				Nullable: false,
			},
			expected: `ALTER TABLE "test_table"
ADD COLUMN IF NOT EXISTS
  "col2" INTEGER NOT NULL`,
		},
		{
			name: "Add column with default",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			col: &Column{
				Name:     "col3",
				Type:     "TIMESTAMP",
				Nullable: false,
				Default:  "NOW()",
			},
			expected: `ALTER TABLE "test_table"
ADD COLUMN IF NOT EXISTS
  "col3" TIMESTAMP NOT NULL DEFAULT NOW()`,
		},
		{
			name: "Add column with collation",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			col: &Column{
				Name:      "col4",
				Type:      "TEXT",
				Nullable:  true,
				Collation: "en_US",
			},
			expected: `ALTER TABLE "test_table"
ADD COLUMN IF NOT EXISTS
  "col4" TEXT COLLATE "en_US"`,
		},
		{
			name: "Add column with override expression",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			col: &Column{
				Name:               "col5",
				OverrideExpression: "JSONB NOT NULL DEFAULT '{}'::jsonb",
			},
			expected: `ALTER TABLE "test_table"
ADD COLUMN IF NOT EXISTS
  "col5" JSONB NOT NULL DEFAULT '{}'::jsonb`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := col2alter(test.desc, test.col)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestCol2spec(t *testing.T) {
	tests := []struct {
		name     string
		col      *Column
		expected string
	}{
		{
			name: "Simple column",
			col: &Column{
				Name:     "col1",
				Type:     "TEXT",
				Nullable: true,
			},
			expected: `  "col1" TEXT`,
		},
		{
			name: "NOT NULL column",
			col: &Column{
				Name:     "col2",
				Type:     "INTEGER",
				Nullable: false,
			},
			expected: `  "col2" INTEGER NOT NULL`,
		},
		{
			name: "Column with default",
			col: &Column{
				Name:     "col3",
				Type:     "TIMESTAMP",
				Nullable: false,
				Default:  "NOW()",
			},
			expected: `  "col3" TIMESTAMP NOT NULL DEFAULT NOW()`,
		},
		{
			name: "Column with collation",
			col: &Column{
				Name:      "col4",
				Type:      "TEXT",
				Nullable:  true,
				Collation: "en_US",
			},
			expected: `  "col4" TEXT COLLATE "en_US"`,
		},
		{
			name: "Column with override expression",
			col: &Column{
				Name:               "col5",
				OverrideExpression: "JSONB NOT NULL DEFAULT '{}'::jsonb",
			},
			expected: `  "col5" JSONB NOT NULL DEFAULT '{}'::jsonb`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := col2spec(test.col)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestKsuidColOverrideExpression(t *testing.T) {
	tests := []struct {
		name     string
		col      *Column
		expected string
	}{
		{
			name: "Simple KSUID column",
			col: &Column{
				Type:     "TEXT",
				Nullable: false,
			},
			expected: `TEXT NOT NULL COLLATE "C"`,
		},
		{
			name: "KSUID column with default",
			col: &Column{
				Type:     "TEXT",
				Nullable: false,
				Default:  "gen_random_uuid()",
			},
			expected: `TEXT NOT NULL DEFAULT gen_random_uuid() COLLATE "C"`,
		},
		{
			name: "KSUID column with collation", // please don't try using Collation with KSUID columns
			col: &Column{
				Type:      "TEXT",
				Nullable:  true,
				Collation: "en_US",
			},
			expected: `TEXT COLLATE "C"`,
		},
		{
			name: "Nullable KSUID column",
			col: &Column{
				Type:     "TEXT",
				Nullable: true,
			},
			expected: `TEXT COLLATE "C"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ksuidColOverrideExpression(test.col)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

type mockDescriptorWithStorageParams struct {
	mockDescriptor
	storageParams *MessageOptions_StorageParameters
}

func (m *mockDescriptorWithStorageParams) GetStorageParameters() *MessageOptions_StorageParameters {
	return m.storageParams
}

func TestStorageParams2with(t *testing.T) {
	makeWithVacuumThreshold := func(threshold int32) *MessageOptions_StorageParameters {
		sp := &MessageOptions_StorageParameters{}
		sp.SetAutovacuumVacuumThreshold(threshold)
		return sp
	}

	makeMultiple := func() *MessageOptions_StorageParameters {
		sp := &MessageOptions_StorageParameters{}
		sp.SetAutovacuumVacuumThreshold(5000)
		sp.SetAutovacuumVacuumScaleFactor(0.1)
		sp.SetFillfactor(90)
		return sp
	}

	tests := []struct {
		name     string
		desc     Descriptor
		expected string
	}{
		{
			name: "No storage parameters",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			expected: "",
		},
		{
			name: "Vacuum threshold only",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeWithVacuumThreshold(10000),
			},
			expected: `WITH (
  autovacuum_vacuum_threshold = 10000
)`,
		},
		{
			name: "Multiple options",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeMultiple(),
			},
			expected: `WITH (
  autovacuum_vacuum_threshold = 5000,
  autovacuum_vacuum_scale_factor = 0.1,
  fillfactor = 90
)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := storageParams2with(test.desc)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestStorageParams2alter(t *testing.T) {
	makeWithVacuumThreshold := func(threshold int32) *MessageOptions_StorageParameters {
		sp := &MessageOptions_StorageParameters{}
		sp.SetAutovacuumVacuumThreshold(threshold)
		return sp
	}

	makeMultiple := func() *MessageOptions_StorageParameters {
		sp := &MessageOptions_StorageParameters{}
		sp.SetAutovacuumVacuumThreshold(5000)
		sp.SetAutovacuumVacuumScaleFactor(0.1)
		sp.SetFillfactor(90)
		return sp
	}

	tests := []struct {
		name           string
		desc           Descriptor
		existingParams map[string]string
		expected       string
	}{
		{
			name: "No storage parameters in descriptor",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			existingParams: map[string]string{
				"fillfactor": "80",
			},
			expected: "",
		},
		{
			name: "No existing parameters, set new one",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeWithVacuumThreshold(10000),
			},
			existingParams: map[string]string{},
			expected: `ALTER TABLE "test_table"
SET (
  autovacuum_vacuum_threshold = 10000
)
`,
		},
		{
			name: "Existing parameter matches desired, no update needed",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeWithVacuumThreshold(10000),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_threshold": "10000",
			},
			expected: "",
		},
		{
			name: "Existing parameter differs, update needed",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeWithVacuumThreshold(10000),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_threshold": "5000",
			},
			expected: `ALTER TABLE "test_table"
SET (
  autovacuum_vacuum_threshold = 10000
)
`,
		},
		{
			name: "Multiple parameters, some match some differ",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeMultiple(),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_threshold":    "5000", // matches
				"autovacuum_vacuum_scale_factor": "0.2",  // differs (0.1 desired)
				"fillfactor":                     "80",   // differs (90 desired)
			},
			expected: `ALTER TABLE "test_table"
SET (
  autovacuum_vacuum_scale_factor = 0.1,
  fillfactor = 90
)
`,
		},
		{
			name: "Existing parameter not in descriptor, leave it alone",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams:  makeWithVacuumThreshold(10000),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_threshold": "10000",
				"fillfactor":                  "80", // Not in descriptor, should be left alone
			},
			expected: "",
		},
		{
			name: "Float parameter with precision difference",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams: func() *MessageOptions_StorageParameters {
					sp := &MessageOptions_StorageParameters{}
					sp.SetAutovacuumVacuumScaleFactor(0.1)
					return sp
				}(),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_scale_factor": "0.100000001", // Very close, should match
			},
			expected: "",
		},
		{
			name: "Float parameter with significant difference",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams: func() *MessageOptions_StorageParameters {
					sp := &MessageOptions_StorageParameters{}
					sp.SetAutovacuumVacuumScaleFactor(0.1)
					return sp
				}(),
			},
			existingParams: map[string]string{
				"autovacuum_vacuum_scale_factor": "0.2", // Different, should update
			},
			expected: `ALTER TABLE "test_table"
SET (
  autovacuum_vacuum_scale_factor = 0.1
)
`,
		},
		{
			name: "Boolean parameter",
			desc: &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: "test_table"},
				storageParams: func() *MessageOptions_StorageParameters {
					sp := &MessageOptions_StorageParameters{}
					sp.SetAutovacuumEnabled(false)
					return sp
				}(),
			},
			existingParams: map[string]string{
				"autovacuum_enabled": "true", // Different, should update
			},
			expected: `ALTER TABLE "test_table"
SET (
  autovacuum_enabled = false
)
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := storageParams2alter(test.desc, test.existingParams)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestPgQuoteIdent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase no special chars",
			input:    "pbidx_test_abc123",
			expected: "pbidx_test_abc123",
		},
		{
			name:     "contains dollar sign",
			input:    "pb$tenant_id",
			expected: `"pb$tenant_id"`,
		},
		{
			name:     "contains uppercase",
			input:    "pbidx_vector_index_MODEL_3DIMS_abc123",
			expected: `"pbidx_vector_index_MODEL_3DIMS_abc123"`,
		},
		{
			name:     "simple table name",
			input:    "pb_pet_models_animals_v1_8a3723d5",
			expected: "pb_pet_models_animals_v1_8a3723d5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := pgQuoteIdent(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestPgNormalizeExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "column IS NULL",
			input:    "pb$deleted_at IS NULL",
			expected: `"pb$deleted_at" IS NULL`,
		},
		{
			name:     "override expression with ops class",
			input:    "pb$min_hash bit_hamming_ops",
			expected: `"pb$min_hash" bit_hamming_ops`,
		},
		{
			name:     "override expression with vector ops",
			input:    "pb$model_embeddings_1 vector_cosine_ops",
			expected: `"pb$model_embeddings_1" vector_cosine_ops`,
		},
		{
			name:     "already quoted identifier",
			input:    `"pb$deleted_at" IS NULL`,
			expected: `"pb$deleted_at" IS NULL`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := pgNormalizeExpr(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestIndex2expectedDef(t *testing.T) {
	tests := []struct {
		name     string
		desc     Descriptor
		idx      *Index
		expected string
	}{
		{
			name: "BTREE index with multiple columns",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:    "pbidx_test1_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns: []string{"pb$tenant_id", "pb$pk", "pb$sk"},
			},
			expected: `CREATE INDEX pbidx_test1_abc123 ON public.pb_test_table_abc123 USING btree ("pb$tenant_id", "pb$pk", "pb$sk")`,
		},
		{
			name: "GIN index",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:    "pbidx_test2_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_GIN,
				Columns: []string{"pb$tenant_id", "pb$profile"},
			},
			expected: `CREATE INDEX pbidx_test2_abc123 ON public.pb_test_table_abc123 USING gin ("pb$tenant_id", "pb$profile")`,
		},
		{
			name: "BTREE_GIN index renders as gin",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:    "pbidx_fts_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
				Columns: []string{"pb$tenant_id", "pb$fts_data"},
			},
			expected: `CREATE INDEX pbidx_fts_abc123 ON public.pb_test_table_abc123 USING gin ("pb$tenant_id", "pb$fts_data")`,
		},
		{
			name: "UNIQUE BTREE index",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:     "pbidx_unique_abc123",
				Method:   MessageOptions_Index_INDEX_METHOD_BTREE,
				IsUnique: true,
				Columns:  []string{"pb$tenant_id", "pb$pksk"},
			},
			expected: `CREATE UNIQUE INDEX pbidx_unique_abc123 ON public.pb_test_table_abc123 USING btree ("pb$tenant_id", "pb$pksk")`,
		},
		{
			name: "BTREE index with WHERE predicate",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:           "pbidx_partial_abc123",
				Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:        []string{"pb$id", "pb$tenant_id"},
				WherePredicate: `pb$deleted_at IS NULL`,
			},
			expected: `CREATE INDEX pbidx_partial_abc123 ON public.pb_test_table_abc123 USING btree ("pb$id", "pb$tenant_id") WHERE ("pb$deleted_at" IS NULL)`,
		},
		{
			name: "HNSW index with override expression",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:               "pbidx_hnsw_abc123",
				Method:             MessageOptions_Index_INDEX_METHOD_HNSW_COSINE,
				Columns:            []string{"pb$min_hash"},
				OverrideExpression: "pb$min_hash bit_hamming_ops",
			},
			expected: `CREATE INDEX pbidx_hnsw_abc123 ON public.pb_test_table_abc123 USING hnsw ("pb$min_hash" bit_hamming_ops)`,
		},
		{
			name: "HNSW index with vector_cosine_ops",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:               "pbidx_vector_abc123",
				Method:             MessageOptions_Index_INDEX_METHOD_HNSW_COSINE,
				Columns:            []string{"pb$model_embeddings_1"},
				OverrideExpression: "pb$model_embeddings_1 vector_cosine_ops",
			},
			expected: `CREATE INDEX pbidx_vector_abc123 ON public.pb_test_table_abc123 USING hnsw ("pb$model_embeddings_1" vector_cosine_ops)`,
		},
		{
			name: "Index with uppercase in name",
			desc: &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx: &Index{
				Name:               "pbidx_vector_index_MODEL_3DIMS_abc123",
				Method:             MessageOptions_Index_INDEX_METHOD_HNSW_COSINE,
				Columns:            []string{"pb$model_embeddings_1"},
				OverrideExpression: "pb$model_embeddings_1 vector_cosine_ops",
			},
			expected: `CREATE INDEX "pbidx_vector_index_MODEL_3DIMS_abc123" ON public.pb_test_table_abc123 USING hnsw ("pb$model_embeddings_1" vector_cosine_ops)`,
		},
		{
			name: "Partitioned table uses ON ONLY",
			desc: &mockDescriptor{tableName: "pb_pasta_models_food_v1_29fd1107", isPartitioned: true},
			idx: &Index{
				Name:    "pbidx_pksk_split2_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns: []string{"pb$tenant_id", "pb$pk", "pb$sk"},
			},
			expected: `CREATE INDEX pbidx_pksk_split2_abc123 ON ONLY public.pb_pasta_models_food_v1_29fd1107 USING btree ("pb$tenant_id", "pb$pk", "pb$sk")`,
		},
		{
			name: "Date partitioned table uses ON ONLY",
			desc: &mockDescriptor{tableName: "pb_garlic_abc123", isPartitionedByCreatedAt: true},
			idx: &Index{
				Name:    "pbidx_fts_data_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
				Columns: []string{"pb$tenant_id", "pb$fts_data"},
			},
			expected: `CREATE INDEX pbidx_fts_data_abc123 ON ONLY public.pb_garlic_abc123 USING gin ("pb$tenant_id", "pb$fts_data")`,
		},
		{
			name: "KSUID partitioned table uses ON ONLY",
			desc: &mockDescriptor{tableName: "pb_cheese_abc123", partitionedByKsuidFieldName: "event_id"},
			idx: &Index{
				Name:    "pbidx_fts_data_abc123",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE_GIN,
				Columns: []string{"pb$tenant_id", "pb$fts_data"},
			},
			expected: `CREATE INDEX pbidx_fts_data_abc123 ON ONLY public.pb_cheese_abc123 USING gin ("pb$tenant_id", "pb$fts_data")`,
		},
		{
			name:     "Dropped index returns empty",
			desc:     &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx:      &Index{Name: "pbidx_dropped", IsDropped: true},
			expected: "",
		},
		{
			name:     "Primary index returns empty",
			desc:     &mockDescriptor{tableName: "pb_test_table_abc123"},
			idx:      &Index{Name: "pbidx_primary", IsPrimary: true},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := index2expectedDef(test.desc, test.idx)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}
