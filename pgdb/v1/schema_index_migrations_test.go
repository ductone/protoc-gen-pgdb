package v1

import (
	"strings"
	"testing"
)

// The exact sequence from ductone/c1 PR #24781: a replacement ordered index
// declared alongside a superseded prefix index marked dropped. Regardless of
// which order the two appear in desc.Indexes() (proto declaration order),
// indexMigrations must emit the CREATE before the DROP -- a lock_timeout skip
// on the CREATE must not leave the DROP to run and remove the only surviving
// access path.
func TestIndexMigrations_ReplacementCreatedBeforeOldDropped(t *testing.T) {
	oldPrefix := &Index{
		Name:      "pbidx_app_and_connector_id",
		Method:    MessageOptions_Index_INDEX_METHOD_BTREE,
		IsDropped: true,
	}
	newOrdered := &Index{
		Name:    "pbidx_app_connector_created_at",
		Method:  MessageOptions_Index_INDEX_METHOD_BTREE,
		Columns: []string{"tenant_id", "app_id", "connector_id", "created_at"},
	}
	existing := map[string]struct{}{"pbidx_app_and_connector_id": {}}

	tests := []struct {
		name    string
		indexes []*Index
	}{
		{name: "old declared first (the #24781 shape)", indexes: []*Index{oldPrefix, newOrdered}},
		{name: "new declared first", indexes: []*Index{newOrdered, oldPrefix}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desc := &mockDescriptor{tableName: "completed_connector_sync_lifecycle", indexes: test.indexes}
			got := indexMigrations(desc, existing)

			if len(got) != 2 {
				t.Fatalf("expected 2 statements (1 create, 1 drop), got %d:\n%v", len(got), got)
			}
			if !strings.HasPrefix(got[0], "CREATE INDEX") {
				t.Errorf("statement 0 must be the CREATE, got:\n%s", got[0])
			}
			if !strings.HasPrefix(got[1], "DROP INDEX") {
				t.Errorf("statement 1 must be the DROP, got:\n%s", got[1])
			}
			if !strings.Contains(got[1], "pbidx_app_and_connector_id") {
				t.Errorf("DROP must target the old index, got:\n%s", got[1])
			}
			if !strings.Contains(got[0], "pbidx_app_connector_created_at") {
				t.Errorf("CREATE must target the new index, got:\n%s", got[0])
			}
		})
	}
}

// Multiple creates and multiple drops interleaved in declaration order: every
// create precedes every drop, and relative order within each group is
// preserved (declaration order), so an unrelated pair of drops/creates on the
// same table doesn't get scrambled.
func TestIndexMigrations_MultipleCreatesAndDropsGroupedStable(t *testing.T) {
	dropA := &Index{Name: "drop_a", IsDropped: true}
	createA := &Index{Name: "create_a", Method: MessageOptions_Index_INDEX_METHOD_BTREE}
	dropB := &Index{Name: "drop_b", IsDropped: true}
	createB := &Index{Name: "create_b", Method: MessageOptions_Index_INDEX_METHOD_BTREE}

	desc := &mockDescriptor{
		tableName: "t",
		indexes:   []*Index{dropA, createA, dropB, createB},
	}
	existing := map[string]struct{}{"drop_a": {}, "drop_b": {}}

	got := indexMigrations(desc, existing)
	if len(got) != 4 {
		t.Fatalf("expected 4 statements, got %d:\n%v", len(got), got)
	}

	wantOrder := []string{"create_a", "create_b", "drop_a", "drop_b"}
	for i, name := range wantOrder {
		if !strings.Contains(got[i], `"`+name+`"`) {
			t.Errorf("statement %d: expected to reference %q, got:\n%s", i, name, got[i])
		}
	}
	for i := range 2 {
		if !strings.HasPrefix(got[i], "CREATE INDEX") {
			t.Errorf("statement %d must be a CREATE, got:\n%s", i, got[i])
		}
	}
	for i := 2; i < 4; i++ {
		if !strings.HasPrefix(got[i], "DROP INDEX") {
			t.Errorf("statement %d must be a DROP, got:\n%s", i, got[i])
		}
	}
}

// Existing behavior must survive the refactor into indexMigrations: a
// not-dropped index that already exists, and a dropped index that never
// existed, both produce no statement.
func TestIndexMigrations_NoOpCases(t *testing.T) {
	alreadyThere := &Index{Name: "already_there", Method: MessageOptions_Index_INDEX_METHOD_BTREE}
	neverExisted := &Index{Name: "never_existed", IsDropped: true}
	primary := &Index{Name: "pk", IsPrimary: true}

	desc := &mockDescriptor{tableName: "t", indexes: []*Index{alreadyThere, neverExisted, primary}}
	existing := map[string]struct{}{"already_there": {}}

	got := indexMigrations(desc, existing)
	if len(got) != 0 {
		t.Fatalf("expected no statements, got %d:\n%v", len(got), got)
	}
}

// A missing, not-dropped index still produces a single CREATE (unchanged
// single-index behavior).
func TestIndexMigrations_SingleMissingIndexCreates(t *testing.T) {
	missing := &Index{Name: "missing", Method: MessageOptions_Index_INDEX_METHOD_BTREE}
	desc := &mockDescriptor{tableName: "t", indexes: []*Index{missing}}

	got := indexMigrations(desc, map[string]struct{}{})
	if len(got) != 1 || !strings.HasPrefix(got[0], "CREATE INDEX") {
		t.Fatalf("expected a single CREATE INDEX statement, got:\n%v", got)
	}
}
