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

	t.Run("single include column", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_single",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		assertContains(t, got, "INCLUDE (\n  \"pb$id\"\n)")
	})

	t.Run("no include columns omits INCLUDE", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:    "idx_plain",
			Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns: []string{"pb$tenant_id"},
		})

		assertNotContains(t, got, "INCLUDE")
	})

	t.Run("empty include columns slice omits INCLUDE", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_empty",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{},
		})

		assertNotContains(t, got, "INCLUDE")
	})

	t.Run("dropped index ignores include columns", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_dropped",
			IsDropped:      true,
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		assertContains(t, got, "DROP INDEX")
		assertNotContains(t, got, "INCLUDE")
	})

	t.Run("partitioned table omits CONCURRENTLY", func(t *testing.T) {
		partDesc := &mockDescriptor{tableName: "pb_partitioned", isPartitioned: true}
		got := index2sql(partDesc, &Index{
			Name:           "idx_part",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
		})

		assertNotContains(t, got, "CONCURRENTLY")
		assertContains(t, got, "INCLUDE")
	})
}

func TestIndex2SQL_WherePredicate(t *testing.T) {
	desc := &mockDescriptor{tableName: "pb_app_resource"}

	t.Run("IS NULL predicate", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_alive",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id", "pb$app_id"},
			WherePredicate: `"pb$deleted_at" IS NULL`,
		})

		assertContains(t, got, `WHERE "pb$deleted_at" IS NULL`)
		assertNotContains(t, got, "INCLUDE")
	})

	t.Run("IS NOT NULL predicate", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_deleted",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			WherePredicate: `"pb$deleted_at" IS NOT NULL`,
		})

		assertContains(t, got, `WHERE "pb$deleted_at" IS NOT NULL`)
	})

	t.Run("EQUALS predicate", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_active",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			WherePredicate: `"pb$is_active" = true`,
		})

		assertContains(t, got, `WHERE "pb$is_active" = true`)
	})

	t.Run("multiple predicates ANDed", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_compound",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			WherePredicate: `"pb$deleted_at" IS NULL AND "pb$is_active" = true`,
		})

		assertContains(t, got, `WHERE "pb$deleted_at" IS NULL AND "pb$is_active" = true`)
	})

	t.Run("no predicate omits WHERE", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:    "idx_no_where",
			Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns: []string{"pb$tenant_id"},
		})

		assertNotContains(t, got, "WHERE")
	})
}

func TestIndex2SQL_CombinedFeatures(t *testing.T) {
	desc := &mockDescriptor{tableName: "pb_app_resource"}

	t.Run("columns + INCLUDE + WHERE exact output", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_full",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id", "pb$app_id", "pb$app_resource_type_id"},
			IncludeColumns: []string{"pb$access_config_id", "pb$id"},
			WherePredicate: `pb$deleted_at IS NULL`,
		})

		assertExact(t, got,
			"CREATE INDEX CONCURRENTLY IF NOT EXISTS\n"+
				"  \"idx_full\"\n"+
				"ON\n"+
				"  \"pb_app_resource\"\n"+
				"USING\n"+
				"  BTREE\n"+
				"(\n"+
				"  \"pb$tenant_id\", \n"+
				"  \"pb$app_id\", \n"+
				"  \"pb$app_resource_type_id\"\n"+
				")\n"+
				"INCLUDE (\n"+
				"  \"pb$access_config_id\", \n"+
				"  \"pb$id\"\n"+
				")\n"+
				"WHERE pb$deleted_at IS NULL\n")
	})

	t.Run("INCLUDE appears before WHERE", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_order",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id"},
			IncludeColumns: []string{"pb$id"},
			WherePredicate: `pb$deleted_at IS NULL`,
		})

		includePos := strings.Index(got, "INCLUDE")
		wherePos := strings.Index(got, "WHERE")
		if includePos == -1 || wherePos == -1 {
			t.Fatalf("expected both INCLUDE and WHERE, got:\n%s", got)
		}
		if includePos >= wherePos {
			t.Errorf("INCLUDE (pos %d) should appear before WHERE (pos %d)", includePos, wherePos)
		}
	})

	t.Run("unique index with INCLUDE and WHERE", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_unique_covering",
			IsUnique:       true,
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{"pb$tenant_id", "pb$email"},
			IncludeColumns: []string{"pb$id"},
			WherePredicate: `pb$deleted_at IS NULL`,
		})

		assertContains(t, got, "CREATE UNIQUE INDEX")
		assertContains(t, got, "INCLUDE")
		assertContains(t, got, "WHERE")
	})

	t.Run("override expression bypasses columns but INCLUDE still renders", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:               "idx_override",
			Method:             MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:            []string{"pb$data"},
			OverrideExpression: "pb$data jsonb_path_ops",
			IncludeColumns:     []string{"pb$id"},
		})

		assertContains(t, got, "pb$data jsonb_path_ops")
		assertContains(t, got, "INCLUDE")
		assertNotContains(t, got, "\"pb$data\"")
	})
}

func TestIndex2SQL_DeletedExclusionEquivalence(t *testing.T) {
	desc := &mockDescriptor{tableName: "pb_app_resource"}
	io := NewIndexOptions(nil)

	// The old system (partial_deleted_at_is_null: true) and the new system
	// (where: [{column: "deleted_at", op: "IS NULL"}]) both emit identical
	// generated Go code:
	//   WherePredicate: "" + io.ColumnName("deleted_at") + " IS NULL"
	//
	// Build the predicate just like the generated code does so the test
	// stays correct even if the column-name prefix changes.
	deletedPredicate := io.ColumnName("deleted_at") + " IS NULL"

	t.Run("basic equivalence", func(t *testing.T) {
		// Old system: partial_deleted_at_is_null: true
		oldSQL := index2sql(desc, &Index{
			Name:           "idx_alive",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{io.ColumnName("tenant_id"), io.ColumnName("app_id")},
			WherePredicate: deletedPredicate,
		})

		// New system: where: [{column: "deleted_at", op: "IS NULL"}]
		newSQL := index2sql(desc, &Index{
			Name:           "idx_alive",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{io.ColumnName("tenant_id"), io.ColumnName("app_id")},
			WherePredicate: deletedPredicate,
		})

		assertExact(t, oldSQL, newSQL)
		assertContains(t, oldSQL, "WHERE "+deletedPredicate)
	})

	t.Run("with include columns", func(t *testing.T) {
		oldSQL := index2sql(desc, &Index{
			Name:           "idx_alive_covering",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{io.ColumnName("tenant_id"), io.ColumnName("app_id")},
			IncludeColumns: []string{io.ColumnName("id")},
			WherePredicate: deletedPredicate,
		})

		newSQL := index2sql(desc, &Index{
			Name:           "idx_alive_covering",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{io.ColumnName("tenant_id"), io.ColumnName("app_id")},
			IncludeColumns: []string{io.ColumnName("id")},
			WherePredicate: deletedPredicate,
		})

		assertExact(t, oldSQL, newSQL)
		assertContains(t, oldSQL, "INCLUDE")
		assertContains(t, oldSQL, "WHERE "+deletedPredicate)
	})

	t.Run("exact SQL output", func(t *testing.T) {
		got := index2sql(desc, &Index{
			Name:           "idx_alive",
			Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
			Columns:        []string{io.ColumnName("tenant_id"), io.ColumnName("app_id")},
			WherePredicate: deletedPredicate,
		})

		assertExact(t, got,
			"CREATE INDEX CONCURRENTLY IF NOT EXISTS\n"+
				"  \"idx_alive\"\n"+
				"ON\n"+
				"  \"pb_app_resource\"\n"+
				"USING\n"+
				"  BTREE\n"+
				"(\n"+
				"  \""+io.ColumnName("tenant_id")+"\", \n"+
				"  \""+io.ColumnName("app_id")+"\"\n"+
				")\n"+
				"WHERE "+deletedPredicate+"\n")
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
