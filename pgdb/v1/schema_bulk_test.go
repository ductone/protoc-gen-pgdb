package v1

import (
	"testing"
)

func TestCatalogSnapshot_ColumnsForTable_Missing(t *testing.T) {
	snap := &CatalogSnapshot{
		columns:       make(map[string]map[string]struct{}),
		indexes:       make(map[string]map[string]struct{}),
		stats:         make(map[string]map[string]struct{}),
		storageParams: make(map[string]map[string]string),
	}

	cols := snap.columnsForTable("nonexistent_table")
	if cols != nil {
		t.Errorf("expected nil for nonexistent table, got %v", cols)
	}
	// len(nil map) == 0, which triggers CreateSchema path in MigrationsWithCatalog
	if len(cols) != 0 {
		t.Errorf("expected len 0 for nil columns, got %d", len(cols))
	}
}

func TestCatalogSnapshot_Accessors(t *testing.T) {
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

	// Missing tables return nil (safe for len and lookups)
	if snap.indexesForTable("other") != nil {
		t.Error("expected nil for missing table indexes")
	}
	if snap.statsForTable("other") != nil {
		t.Error("expected nil for missing table stats")
	}
	if snap.storageParamsForTable("other") != nil {
		t.Error("expected nil for missing table storage params")
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
}

func TestMigrationsWithCatalog_NilSnapshot(t *testing.T) {
	_, err := MigrationsWithCatalog(nil, nil, 0)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}
