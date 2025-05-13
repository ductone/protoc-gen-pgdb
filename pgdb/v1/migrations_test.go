package v1

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	"github.com/stretchr/testify/require"
)

type mockDBReflectMessage struct {
	dbr *mockDBReflect
}

func (m *mockDBReflectMessage) DBReflect() Message {
	return m.dbr
}

type mockDBReflect struct {
	desc *migrationsTestDescriptor
}

func (m *mockDBReflect) Descriptor() Descriptor {
	return m.desc
}

func (m *mockDBReflect) Record(opts ...RecordOptionsFunc) (exp.Record, error) {
	return nil, nil
}

func (m *mockDBReflect) SearchData(opts ...RecordOptionsFunc) []*SearchContent {
	return nil
}

type migrationsTestDescriptor struct {
	tableName                   string
	isPartitioned               bool
	isPartitionedByCreatedAt    bool
	partitionedByKsuidFieldName string
	partitionDateRange          MessageOptions_PartitionedByDateRange
	fields                      []*Column
	indexes                     []*Index
	statistics                  []*Statistic
}

func (m *migrationsTestDescriptor) TableName() string {
	return m.tableName
}

func (m *migrationsTestDescriptor) Fields(opts ...DescriptorFieldOptionFunc) []*Column {
	return m.fields
}

func (m *migrationsTestDescriptor) PKSKField() *Column {
	return nil
}

func (m *migrationsTestDescriptor) PKSKV2Field() *Column {
	return nil
}

func (m *migrationsTestDescriptor) DataField() *Column {
	return nil
}

func (m *migrationsTestDescriptor) SearchField() *Column {
	return nil
}

func (m *migrationsTestDescriptor) VersioningField() *Column {
	return nil
}

func (m *migrationsTestDescriptor) TenantField() *Column {
	return nil
}

func (m *migrationsTestDescriptor) IsPartitioned() bool {
	return m.isPartitioned
}

func (m *migrationsTestDescriptor) IsPartitionedByCreatedAt() bool {
	return m.isPartitionedByCreatedAt
}

func (m *migrationsTestDescriptor) GetPartitionedByKsuidFieldName() string {
	return m.partitionedByKsuidFieldName
}

func (m *migrationsTestDescriptor) Indexes(opts ...IndexOptionsFunc) []*Index {
	return m.indexes
}

func (m *migrationsTestDescriptor) IndexPrimaryKey(opts ...IndexOptionsFunc) *Index {
	for _, idx := range m.indexes {
		if idx.IsPrimary {
			return idx
		}
	}
	return nil
}

func (m *migrationsTestDescriptor) Statistics(opts ...StatisticOptionsFunc) []*Statistic {
	return m.statistics
}

func (m *migrationsTestDescriptor) GetPartitionDateRange() MessageOptions_PartitionedByDateRange {
	return m.partitionDateRange
}

func TestMigrationsTableDoesNotExist(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	msg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_table",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_table",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, msg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for table creation")

	require.Contains(t, migrations[0], "CREATE TABLE IF NOT EXISTS")
	require.Contains(t, migrations[0], "test_migrations_table")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, msg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after table creation")
}

func TestMigrationsColumnNeedsToBeAdded(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	initialMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_column",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_column",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	initialSchema, err := CreateSchema(initialMsg)
	require.NoError(t, err)
	for _, stmt := range initialSchema {
		_, err := pg.DB.Exec(ctx, stmt)
		require.NoError(t, err, "Failed to execute initial schema creation: %s", stmt)
	}

	updatedMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_column",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: true,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_column",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for adding column")

	require.Contains(t, migrations[0], "ALTER TABLE")
	require.Contains(t, migrations[0], "ADD COLUMN")
	require.Contains(t, migrations[0], "pb$description")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after adding column")
}

func TestMigrationsColumnNeedsToBeUpdated(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	initialMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_update_column",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: true,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_update_column",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	initialSchema, err := CreateSchema(initialMsg)
	require.NoError(t, err)
	for _, stmt := range initialSchema {
		_, err := pg.DB.Exec(ctx, stmt)
		require.NoError(t, err, "Failed to execute initial schema creation: %s", stmt)
	}

	updatedMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_update_column",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false, // Changed from true to false
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_update_column",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for updating column")

	require.Contains(t, migrations[0], "ALTER TABLE")
	require.Contains(t, migrations[0], "ALTER COLUMN")
	require.Contains(t, migrations[0], "pb$description")
	require.Contains(t, migrations[0], "SET NOT NULL")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after updating column")
}

func TestMigrationsIndexNeedsToBeAdded(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	initialMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_index",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_index",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	initialSchema, err := CreateSchema(initialMsg)
	require.NoError(t, err)
	for _, stmt := range initialSchema {
		_, err := pg.DB.Exec(ctx, stmt)
		require.NoError(t, err, "Failed to execute initial schema creation: %s", stmt)
	}

	updatedMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_index",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_index",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
					{
						Name:    "idx_test_migrations_description",
						Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
						Columns: []string{"pb$description"},
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for adding index")

	require.Contains(t, migrations[0], "CREATE INDEX")
	require.Contains(t, migrations[0], "idx_test_migrations_description")
	require.Contains(t, migrations[0], "pb$description")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after adding index")
}

func TestMigrationsIndexNeedsToBeDropped(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	initialMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_drop_index",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_drop_index",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
					{
						Name:    "idx_test_migrations_description",
						Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
						Columns: []string{"pb$description"},
					},
				},
			},
		},
	}

	initialSchema, err := CreateSchema(initialMsg)
	require.NoError(t, err)
	for _, stmt := range initialSchema {
		_, err := pg.DB.Exec(ctx, stmt)
		require.NoError(t, err, "Failed to execute initial schema creation: %s", stmt)
	}

	updatedMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_drop_index",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_drop_index",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
					{
						Name:      "idx_test_migrations_description",
						Method:    MessageOptions_Index_INDEX_METHOD_BTREE,
						Columns:   []string{"pb$description"},
						IsDropped: true,
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for dropping index")

	require.Contains(t, migrations[0], "DROP INDEX")
	require.Contains(t, migrations[0], "idx_test_migrations_description")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after dropping index")
}

func TestMigrationsStatisticsNeedToBeUpdated(t *testing.T) {
	ctx := context.Background()
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	initialMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_stats",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_stats",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
			},
		},
	}

	initialSchema, err := CreateSchema(initialMsg)
	require.NoError(t, err)
	for _, stmt := range initialSchema {
		_, err := pg.DB.Exec(ctx, stmt)
		require.NoError(t, err, "Failed to execute initial schema creation: %s", stmt)
	}

	updatedMsg := &mockDBReflectMessage{
		dbr: &mockDBReflect{
			desc: &migrationsTestDescriptor{
				tableName: "test_migrations_add_stats",
				fields: []*Column{
					{
						Name:     "pb$id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$tenant_id",
						Type:     "TEXT",
						Nullable: false,
					},
					{
						Name:     "pb$description",
						Type:     "TEXT",
						Nullable: false,
					},
				},
				indexes: []*Index{
					{
						Name:      "pk_test_migrations_add_stats",
						IsPrimary: true,
						Columns:   []string{"pb$id", "pb$tenant_id"},
					},
				},
				statistics: []*Statistic{
					{
						Name:    "stat_test_migrations",
						Columns: []string{"pb$tenant_id", "pb$description"},
						Kinds:   []MessageOptions_Stat_StatsKind{MessageOptions_Stat_STATS_KIND_NDISTINCT},
					},
				},
			},
		},
	}

	migrations, err := Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "Expected migration statements for adding statistics")

	require.Contains(t, migrations[0], "CREATE STATISTICS")
	require.Contains(t, migrations[0], "stat_test_migrations")
	require.Contains(t, migrations[0], "pb$tenant_id")
	require.Contains(t, migrations[0], "pb$description")

	for _, migration := range migrations {
		_, err := pg.DB.Exec(ctx, migration)
		require.NoError(t, err, "Failed to execute migration: %s", migration)
	}

	migrations, err = Migrations(ctx, pg.DB, updatedMsg)
	require.NoError(t, err)
	require.Empty(t, migrations, "Expected no migrations after adding statistics")
}
