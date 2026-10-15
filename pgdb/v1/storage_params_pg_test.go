package v1

import (
	"context"
	"strconv"
	"testing"

	"github.com/ductone/protoc-gen-pgdb/internal/pgtest"
	"github.com/stretchr/testify/require"
)

// TestStorageParamsRoundTrip creates a table with each rendered WITH clause on a
// real server, reads back pg_class.reloptions, and requires the drift ALTER to
// be empty. Unit tests cannot catch a misspelled reloption name, or a value the
// server stores in a form the drift check does not recognise — either would
// make every Migrations run re-emit the same ALTER forever.
func TestStorageParamsRoundTrip(t *testing.T) {
	pg, err := pgtest.Start()
	require.NoError(t, err)
	defer pg.Stop()

	ctx := context.Background()
	var serverVersion int
	require.NoError(t, pg.DB.QueryRow(ctx, "SELECT current_setting('server_version_num')::int").Scan(&serverVersion))

	tests := []struct {
		name       string
		minVersion int
		build      func(sp *MessageOptions_StorageParameters)
	}{
		{
			name: "insert-driven",
			build: func(sp *MessageOptions_StorageParameters) {
				sp.SetAutovacuumVacuumInsertThreshold(1000000)
				sp.SetAutovacuumVacuumInsertScaleFactor(0.005)
			},
		},
		{
			name: "vacuum_index_cleanup off",
			build: func(sp *MessageOptions_StorageParameters) {
				sp.SetVacuumIndexCleanup(MessageOptions_StorageParameters_VACUUM_INDEX_CLEANUP_OFF)
			},
		},
		{
			name: "vacuum_index_cleanup auto",
			build: func(sp *MessageOptions_StorageParameters) {
				sp.SetVacuumIndexCleanup(MessageOptions_StorageParameters_VACUUM_INDEX_CLEANUP_AUTO)
			},
		},
		{
			name:       "PG18 dead-tuple cap and eager freeze rate",
			minVersion: 180000,
			build: func(sp *MessageOptions_StorageParameters) {
				sp.SetAutovacuumVacuumMaxThreshold(5000000)
				sp.SetVacuumMaxEagerFreezeFailureRate(0.05)
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if serverVersion < test.minVersion {
				t.Skipf("server_version_num %d < %d", serverVersion, test.minVersion)
			}
			sp := &MessageOptions_StorageParameters{}
			test.build(sp)
			tableName := "storage_params_round_trip_" + strconv.Itoa(i)
			desc := &mockDescriptorWithStorageParams{
				mockDescriptor: mockDescriptor{tableName: tableName},
				storageParams:  sp,
			}

			_, err := pg.DB.Exec(ctx, `CREATE TABLE "`+tableName+`" (id int)`+partitionStorageParams(desc))
			require.NoError(t, err)

			var reloptions []string
			require.NoError(t, pg.DB.QueryRow(ctx,
				"SELECT reloptions FROM pg_class WHERE relname = $1", tableName).Scan(&reloptions))
			existing := make(map[string]string)
			parseReloptions(reloptions, existing)

			require.Empty(t, storageParams2alterTable(desc, tableName, existing),
				"declared parameters must read back as converged; stored reloptions: %v", reloptions)
		})
	}
}
