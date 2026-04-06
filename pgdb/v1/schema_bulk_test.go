package v1

import (
	"testing"
)

func TestMigrationsWithCatalog_NewTable(t *testing.T) {
	snap := &CatalogSnapshot{
		columns:       make(map[string]map[string]struct{}),
		indexes:       make(map[string]map[string]struct{}),
		stats:         make(map[string]map[string]struct{}),
		storageParams: make(map[string]map[string]string),
	}

	// Table doesn't exist in snapshot → should return CreateSchema output
	cols := snap.columnsForTable("nonexistent_table")
	if len(cols) != 0 {
		t.Errorf("expected empty columns for nonexistent table, got %d", len(cols))
	}
}

func TestMigrationsWithCatalog_ExistingTable(t *testing.T) {
	snap := &CatalogSnapshot{
		columns: map[string]map[string]struct{}{
			"test_table": {
				"col1": {},
				"col2": {},
			},
		},
		indexes: map[string]map[string]struct{}{
			"test_table": {
				"idx1": {},
			},
		},
		stats:         make(map[string]map[string]struct{}),
		storageParams: make(map[string]map[string]string),
	}

	cols := snap.columnsForTable("test_table")
	if len(cols) != 2 {
		t.Errorf("expected 2 columns, got %d", len(cols))
	}

	indexes := snap.indexesForTable("test_table")
	if len(indexes) != 1 {
		t.Errorf("expected 1 index, got %d", len(indexes))
	}

	// Non-existent table should return empty maps (not nil)
	indexes2 := snap.indexesForTable("other_table")
	if indexes2 == nil {
		t.Error("expected non-nil map for missing table indexes")
	}
	if len(indexes2) != 0 {
		t.Errorf("expected 0 indexes for missing table, got %d", len(indexes2))
	}
}

func TestCatalogSnapshot_StorageParams(t *testing.T) {
	snap := &CatalogSnapshot{
		columns: make(map[string]map[string]struct{}),
		indexes: make(map[string]map[string]struct{}),
		stats:   make(map[string]map[string]struct{}),
		storageParams: map[string]map[string]string{
			"test_table": {
				"fillfactor":                  "80",
				"autovacuum_vacuum_threshold": "10000",
			},
		},
	}

	params := snap.storageParamsForTable("test_table")
	if len(params) != 2 {
		t.Errorf("expected 2 storage params, got %d", len(params))
	}
	if params["fillfactor"] != "80" {
		t.Errorf("expected fillfactor=80, got %s", params["fillfactor"])
	}

	// Non-existent table
	params2 := snap.storageParamsForTable("other_table")
	if params2 == nil {
		t.Error("expected non-nil map for missing table storage params")
	}
	if len(params2) != 0 {
		t.Errorf("expected 0 params for missing table, got %d", len(params2))
	}
}
