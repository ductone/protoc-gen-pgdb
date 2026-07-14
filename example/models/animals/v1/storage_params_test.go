package v1

import (
	"testing"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/stretchr/testify/require"
)

// Pins the generated GetStorageParameters output. The template used to emit an
// open-struct composite literal, which does not compile against the
// API_OPAQUE-generated pgdb runtime — this test exists to keep the example
// exercising the builder-form emission.
func TestPetStorageParameters(t *testing.T) {
	pet := (*Pet)(nil)
	desc := pet.DBReflect(pgdb_v1.DefaultDialect).Descriptor()
	sp := desc.GetStorageParameters()
	require.NotNil(t, sp)
	require.Equal(t, int32(1000), sp.GetAutovacuumVacuumThreshold())
	require.InDelta(t, 0.01, sp.GetAutovacuumVacuumScaleFactor(), 1e-9)
	require.Equal(t, int32(90), sp.GetFillfactor())
	require.False(t, sp.HasAutovacuumEnabled())
}
