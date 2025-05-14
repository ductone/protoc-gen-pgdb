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
	return nil
}

func (m *mockDescriptor) PKSKV2Field() *Column {
	return nil
}

func (m *mockDescriptor) DataField() *Column {
	return nil
}

func (m *mockDescriptor) SearchField() *Column {
	return nil
}

func (m *mockDescriptor) VersioningField() *Column {
	return nil
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
	return nil
}

func (m *mockDescriptor) Statistics(opts ...StatisticOptionsFunc) []*Statistic {
	return nil
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

func TestCol2add(t *testing.T) {
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
			result := col2add(test.desc, test.col)
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

func TestCol2alter(t *testing.T) {
	tests := []struct {
		name     string
		desc     Descriptor
		current  *Column
		wanted   *Column
		expected string
	}{
		{
			name: "Change nullability to NOT NULL",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:     "col1",
				Type:     "TEXT",
				Nullable: false,
			},
			wanted: &Column{
				Name:     "col1",
				Type:     "TEXT",
				Nullable: true,
			},
			expected: `ALTER TABLE "test_table"
ALTER COLUMN "col1"
SET NOT NULL`,
		},
		{
			name: "Change nullability to NULL",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:     "col2",
				Type:     "INTEGER",
				Nullable: true,
			},
			wanted: &Column{
				Name:     "col2",
				Type:     "INTEGER",
				Nullable: false,
			},
			expected: `ALTER TABLE "test_table"
ALTER COLUMN "col2"
DROP NOT NULL`,
		},
		{
			name: "Change collation",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:      "col3",
				Type:      "TEXT",
				Nullable:  true,
				Collation: "en_US",
			},
			wanted: &Column{
				Name:      "col3",
				Type:      "TEXT",
				Nullable:  true,
				Collation: "C",
			},
			expected: `ALTER TABLE "test_table"
ALTER COLUMN "col3"
SET DATA TYPE TEXT COLLATE "C"`,
		},
		{
			name: "Multiple changes",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:      "col6",
				Type:      "TEXT",
				Nullable:  true,
				Default:   "'default'",
				Collation: "en_US",
			},
			wanted: &Column{
				Name:      "col6",
				Type:      "TEXT",
				Nullable:  false,
				Default:   "",
				Collation: "C",
			},
			expected: `ALTER TABLE "test_table"
ALTER COLUMN "col6"
DROP NOT NULL, ALTER COLUMN "col6"
SET DATA TYPE TEXT COLLATE "C"`,
		},
		{
			name: "No changes",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:     "col7",
				Type:     "TEXT",
				Nullable: false,
				Default:  "'default'",
			},
			wanted: &Column{
				Name:     "col7",
				Type:     "TEXT",
				Nullable: false,
				Default:  "'default'",
			},
			expected: `ALTER TABLE "test_table"
`,
		},
		{
			name: "Skip collation change with override expression",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			current: &Column{
				Name:               "col8",
				Type:               "TEXT",
				Nullable:           true,
				Collation:          "en_US",
				OverrideExpression: "JSONB NOT NULL DEFAULT '{}'::jsonb",
			},
			wanted: &Column{
				Name:               "col8",
				Type:               "TEXT",
				Nullable:           true,
				Collation:          "C",
				OverrideExpression: "JSONB NOT NULL DEFAULT '{}'::jsonb",
			},
			expected: `ALTER TABLE "test_table"
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := col2alter(test.desc, test.current, test.wanted)
			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
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
				Name:    "idx_test_btree",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns: []string{"col1", "col2"},
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test_btree"
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
			name: "Create GIN index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:    "idx_test_gin",
				Method:  MessageOptions_Index_INDEX_METHOD_GIN,
				Columns: []string{"col3"},
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test_gin"
ON
  "test_table"
USING
  GIN
(
  "col3"
)
`,
		},
		{
			name: "Create UNIQUE index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:     "idx_test_unique",
				IsUnique: true,
				Method:   MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:  []string{"col4"},
			},
			expected: `CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test_unique"
ON
  "test_table"
USING
  BTREE
(
  "col4"
)
`,
		},
		{
			name: "Create index with WHERE predicate",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:           "idx_test_where",
				Method:         MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns:        []string{"col5"},
				WherePredicate: "col5 IS NOT NULL",
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test_where"
ON
  "test_table"
USING
  BTREE
(
  "col5"
)
WHERE col5 IS NOT NULL
`,
		},
		{
			name: "Create index with override expression",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:               "idx_test_override",
				Method:             MessageOptions_Index_INDEX_METHOD_GIN,
				OverrideExpression: "to_tsvector('english', col6)",
			},
			expected: `CREATE INDEX CONCURRENTLY IF NOT EXISTS
  "idx_test_override"
ON
  "test_table"
USING
  GIN
(
to_tsvector('english', col6)
)
`,
		},
		{
			name: "Drop index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:      "idx_test_drop",
				IsDropped: true,
			},
			expected: `DROP INDEX CONCURRENTLY IF EXISTS "idx_test_drop"`,
		},
		{
			name: "Drop unique index",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			idx: &Index{
				Name:      "idx_test_drop_unique",
				IsUnique:  true,
				IsDropped: true,
			},
			expected: `DROP INDEX IF EXISTS "idx_test_drop_unique"`,
		},
		{
			name: "Create index on partitioned table",
			desc: &mockDescriptor{
				tableName:     "test_table",
				isPartitioned: true,
			},
			idx: &Index{
				Name:    "idx_test_partitioned",
				Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
				Columns: []string{"col7"},
			},
			expected: `CREATE INDEX IF NOT EXISTS
  "idx_test_partitioned"
ON
  "test_table"
USING
  BTREE
(
  "col7"
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
				Name:    "stat_test_no_kinds",
				Columns: []string{"col1", "col2"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_test_no_kinds" ON "col1","col2" FROM "test_table"
`,
		},
		{
			name: "Create statistics with one kind",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:    "stat_test_ndistinct",
				Kinds:   []MessageOptions_Stat_StatsKind{MessageOptions_Stat_STATS_KIND_NDISTINCT},
				Columns: []string{"col3"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_test_ndistinct"(ndistinct) ON "col3" FROM "test_table"
`,
		},
		{
			name: "Create statistics with multiple kinds",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:    "stat_test_multiple",
				Kinds:   []MessageOptions_Stat_StatsKind{MessageOptions_Stat_STATS_KIND_NDISTINCT, MessageOptions_Stat_STATS_KIND_DEPENDENCIES, MessageOptions_Stat_STATS_KIND_MCV},
				Columns: []string{"col4", "col5"},
			},
			expected: `CREATE STATISTICS IF NOT EXISTS "stat_test_multiple"(ndistinct,dependencies,mcv) ON "col4","col5" FROM "test_table"
`,
		},
		{
			name: "Drop statistics",
			desc: &mockDescriptor{
				tableName: "test_table",
			},
			st: &Statistic{
				Name:      "stat_test_drop",
				IsDropped: true,
			},
			expected: `DROP STATISTICS IF EXISTS "stat_test_drop"`,
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
