package v1

import (
	"bytes"
	"strings"
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

func TestIndex2SQL_IncludeColumns(t *testing.T) {
	desc := &mockDescriptor{tableName: "pb_app_resource"}

	t.Run("btree with include columns", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_covering",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id", "pb$app_id"},
			IncludeColumns: []string{"pb$access_config_id", "pb$id"},
		})

		assertExact(t, got,
			"CREATE INDEX CONCURRENTLY IF NOT EXISTS\n"+
				"  \"idx_covering\"\n"+
				"ON\n"+
				"  \"pb_app_resource\"\n"+
				"USING\n"+
				"  BTREE\n"+
				"(\n"+
				"  \"pb$tenant_id\", \n"+
				"  \"pb$app_id\"\n"+
				")\n"+
				"INCLUDE (\n"+
				"  \"pb$access_config_id\", \n"+
				"  \"pb$id\"\n"+
				")\n")
	})

	t.Run("btree with include columns and partial predicate", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_covering_partial",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id", "pb$app_id"},
			IncludeColumns: []string{"pb$id"},
			WherePredicate: `"pb$deleted_at" IS NULL`,
		})

		expected := "CREATE INDEX CONCURRENTLY IF NOT EXISTS\n" +
			"  \"idx_covering_partial\"\n" +
			"ON\n" +
			"  \"pb_app_resource\"\n" +
			"USING\n" +
			"  BTREE\n" +
			"(\n" +
			"  \"pb$tenant_id\", \n" +
			"  \"pb$app_id\"\n" +
			")\n" +
			"INCLUDE (\n" +
			"  \"pb$id\"\n" +
			")\n" +
			"WHERE \"pb$deleted_at\" IS NULL\n"
		if got != expected {
			t.Errorf("index2sql mismatch.\ngot:\n%s\nexpected:\n%s", got, expected)
		}
	})

	t.Run("btree without include columns omits INCLUDE clause", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:    "idx_no_include",
			Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns: []string{"pb$tenant_id", "pb$app_id"},
		})

		if strings.Contains(got, "INCLUDE") {
			t.Errorf("expected no INCLUDE clause, got:\n%s", got)
		}
	})

	t.Run("single include column", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_single_include",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		expected := "CREATE INDEX CONCURRENTLY IF NOT EXISTS\n" +
			"  \"idx_single_include\"\n" +
			"ON\n" +
			"  \"pb_app_resource\"\n" +
			"USING\n" +
			"  BTREE\n" +
			"(\n" +
			"  \"pb$tenant_id\"\n" +
			")\n" +
			"INCLUDE (\n" +
			"  \"pb$id\"\n" +
			")\n"
		if got != expected {
			t.Errorf("index2sql mismatch.\ngot:\n%s\nexpected:\n%s", got, expected)
		}
	})

	t.Run("empty include columns slice omits INCLUDE clause", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_empty_include",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{},
		})

		if strings.Contains(got, "INCLUDE") {
			t.Errorf("expected no INCLUDE clause for empty slice, got:\n%s", got)
		}
	})

	t.Run("partitioned table omits CONCURRENTLY", func(t *testing.T) {
		partDesc := &mockDescriptor{tableName: "pb_partitioned", isPartitioned: true}
		got := index2sql(partDesc, &Index{
			Name:           "idx_part_include",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		if strings.Contains(got, "CONCURRENTLY") {
			t.Errorf("expected no CONCURRENTLY for partitioned table, got:\n%s", got)
		}
		if !strings.Contains(got, "INCLUDE") {
			t.Errorf("expected INCLUDE clause, got:\n%s", got)
		}
	})

	t.Run("dropped index ignores include columns", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_dropped",
			IsDropped:      true,
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		if strings.Contains(got, "INCLUDE") {
			t.Errorf("expected no INCLUDE clause for dropped index, got:\n%s", got)
		}
		if !strings.Contains(got, "DROP INDEX") {
			t.Errorf("expected DROP INDEX, got:\n%s", got)
		}
	})
}

func assertExact(t *testing.T, got, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("index2sql mismatch.\ngot:\n%s\nexpected:\n%s", got, expected)
	}
}

func assertContains(t *testing.T, got, substr string) {
	t.Helper()
	if !strings.Contains(got, substr) {
		t.Errorf("expected output to contain %q, got:\n%s", substr, got)
	}
}

func assertNotContains(t *testing.T, got, substr string) {
	t.Helper()
	if strings.Contains(got, substr) {
		t.Errorf("expected output to NOT contain %q, got:\n%s", substr, got)
	}
}
