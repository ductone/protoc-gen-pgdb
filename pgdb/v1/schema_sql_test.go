package v1

import (
	"bytes"
	"testing"
)

// mockDescriptor is a mock implementation of the Descriptor interface for testing
type mockDescriptor struct {
	tableName                  string
	isPartitioned              bool
	isPartitionedByCreatedAt   bool
	partitionedByKsuidFieldName string
	partitionDateRange         MessageOptions_PartitionedByDateRange
}

func (m *mockDescriptor) TableName() string {
	return m.tableName
}

func (m *mockDescriptor) Fields(opts ...DescriptorFieldOptionFunc) []*Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) PKSKField() *Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) PKSKV2Field() *Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) DataField() *Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) SearchField() *Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) VersioningField() *Column {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) TenantField() *Column {
	return nil // Not needed for these tests
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
	return nil // Not needed for these tests
}

func (m *mockDescriptor) IndexPrimaryKey(opts ...IndexOptionsFunc) *Index {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) Statistics(opts ...StatisticOptionsFunc) []*Statistic {
	return nil // Not needed for these tests
}

func (m *mockDescriptor) GetPartitionDateRange() MessageOptions_PartitionedByDateRange {
	return m.partitionDateRange
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

func TestIndex2sql(t *testing.T) {
	tests := []struct {
		name     string
		desc     Descriptor
		idx      *Index
		expected string
	}{
		{
			name: "Create BTREE index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:     "idx_test",
				IsUnique: false,
				Method:   MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:  []string{"col1", "col2"},
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test"
ON
  "test_table"
USING
  BTREE
(
  "col1", 
  "col2"
)
`,
		},
		{
			name: "Create unique GIN index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:     "idx_unique_test",
				IsUnique: true,
				Method:   MessageOptions_Index_INDEX_METHOD_GIN,
				Columns:  []string{"col1"},
			},
			expected: `CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_unique_test"
ON
  "test_table"
USING
  GIN
(
  "col1"
)
`,
		},
		{
			name: "Create index with where predicate",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:           "idx_where_test",
				Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:        []string{"col1"},
				WherePredicate: "col1 IS NOT NULL",
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_where_test"
ON
  "test_table"
USING
  BTREE
(
  "col1"
)
WHERE col1 IS NOT NULL
`,
		},
		{
			name: "Create index with override expression",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:               "idx_expr_test",
				Method:             MessageOptions_Index_INDEX_METHOD_BTREE,
				OverrideExpression: "lower(col1)",
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_expr_test"
ON
  "test_table"
USING
  BTREE
(
lower(col1)
)
`,
		},
		{
			name: "Drop index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:      "idx_drop_test",
				IsDropped: true,
			},
			expected: `DROP INDEX CONCURRENTLY IF EXISTS "idx_drop_test"`,
		},
		{
			name: "Drop unique index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:      "idx_unique_drop_test",
				IsUnique:  true,
				IsDropped: true,
			},
			expected: `DROP INDEX IF EXISTS "idx_unique_drop_test"`,
		},
		{
			name: "Create index on partitioned table",
			desc: &mockDescriptor{
				tableName:     "test_table",
				isPartitioned: true,
			},
			idx: &Index{
				Name:     "idx_partitioned_test",
				Method:   MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:  []string{"col1"},
			},
			expected: `CREATE INDEX IF NOT EXISTS
  "idx_partitioned_test"
ON
  "test_table"
USING
  BTREE
(
  "col1"
)
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := index2sql(test.desc, test.idx)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestStatistics2sql(t *testing.T) {
	tests := []struct {
		name     string
		desc     Descriptor
		st       *Statistic
		expected string
	}{
		{
			name: "Create statistics with no kinds",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:    "stat_test",
				Columns: []string{"col1", "col2"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_test" ON "col1","col2" FROM "test_table"
`,
		},
		{
			name: "Create statistics with one kind",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:    "stat_ndistinct",
				Kinds:   []MessageOptions_Stat_StatsKind{MessageOptions_Stat_STATS_KIND_NDISTINCT},
				Columns: []string{"col1"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_ndistinct"(ndistinct) ON "col1" FROM "test_table"
`,
		},
		{
			name: "Create statistics with multiple kinds",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:    "stat_multiple",
				Kinds:   []MessageOptions_Stat_StatsKind{MessageOptions_Stat_STATS_KIND_NDISTINCT, MessageOptions_Stat_STATS_KIND_DEPENDENCIES, MessageOptions_Stat_STATS_KIND_MCV},
				Columns: []string{"col1", "col2", "col3"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_multiple"(ndistinct,dependencies,mcv) ON "col1","col2","col3" FROM "test_table"
`,
		},
		{
			name: "Drop statistics",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:      "stat_drop",
				IsDropped: true,
			},
			expected: `DROP STATISTICS IF EXISTS "stat_drop"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := statistics2sql(test.desc, test.st)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
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