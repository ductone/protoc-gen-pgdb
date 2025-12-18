package v1

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/ductone/protoc-gen-pgdb/internal/slice"
)

func index2sql(desc Descriptor, idx *Index) string {
	buf := &bytes.Buffer{}

	if idx.IsDropped {
		_, _ = buf.WriteString("DROP INDEX")
		// WARNING: unique indexes cannot be dropped
		// concurrently.  Maybe unsafe?
		if !idx.IsUnique && !desc.IsPartitioned() {
			_, _ = buf.WriteString(" CONCURRENTLY")
		}
		_, _ = buf.WriteString(" IF EXISTS ")
		pgWriteString(buf, idx.Name)
		return buf.String()
	}

	_, _ = buf.WriteString("CREATE")
	if idx.IsUnique {
		_, _ = buf.WriteString(" UNIQUE")
	}
	_, _ = buf.WriteString(" INDEX")
	// Note cannot drop or add indexes concurrently on the master partition tables
	if !desc.IsPartitioned() && !desc.IsPartitionedByCreatedAt() && desc.GetPartitionedByKsuidFieldName() == "" {
		_, _ = buf.WriteString(" CONCURRENTLY")
	}
	_, _ = buf.WriteString(" IF NOT EXISTS\n  ")
	pgWriteString(buf, idx.Name)
	_, _ = buf.WriteString("\nON\n  ")
	pgWriteString(buf, desc.TableName())
	_, _ = buf.WriteString("\nUSING\n  ")
	switch idx.Method {
	case MessageOptions_Index_INDEX_METHOD_UNSPECIFIED:
		panic("MessageOptions_Index_INDEX_METHOD_UNSPECIFIED found on " + idx.Name)
	case MessageOptions_Index_INDEX_METHOD_BTREE:
		_, _ = buf.WriteString("BTREE")
	case MessageOptions_Index_INDEX_METHOD_GIN:
		_, _ = buf.WriteString("GIN")
	case MessageOptions_Index_INDEX_METHOD_BTREE_GIN:
		// btree gin just means we can index
		// col types in a multi-col index that aren't
		// noramlly supporte dy gin, eg, varchar,
		// but its not actually a new index type!
		_, _ = buf.WriteString("GIN")
	case MessageOptions_Index_INDEX_METHOD_HNSW_COSINE:
		_, _ = buf.WriteString("HNSW")
	}
	_, _ = buf.WriteString("\n(\n")
	if idx.OverrideExpression != "" {
		_, _ = buf.WriteString(idx.OverrideExpression)
	} else {
		_, _ = buf.WriteString(strings.Join(slice.Convert(idx.Columns, func(in string) string {
			return `  "` + in + `"`
		}), ", \n"))
	}
	_, _ = buf.WriteString("\n)\n")
	if idx.WherePredicate != "" {
		_, _ = buf.WriteString("WHERE ")
		_, _ = buf.WriteString(idx.WherePredicate)
		_, _ = buf.WriteString("\n")
	}
	return buf.String()
}

func statistics2sql(desc Descriptor, st *Statistic) string {
	buf := &bytes.Buffer{}

	if st.IsDropped {
		_, _ = buf.WriteString("DROP STATISTICS")
		_, _ = buf.WriteString(" IF EXISTS ")
		pgWriteString(buf, st.Name)
		return buf.String()
	}

	_, _ = buf.WriteString("CREATE STATISTICS")
	_, _ = buf.WriteString(" IF NOT EXISTS ")
	pgWriteString(buf, st.Name)
	kinds := st.Kinds
	if len(kinds) != 0 {
		_, _ = buf.WriteString("(")
		_, _ = buf.WriteString(strings.Join(slice.Convert(kinds, func(in MessageOptions_Stat_StatsKind) string {
			switch in {
			case MessageOptions_Stat_STATS_KIND_NDISTINCT:
				return "ndistinct"
			case MessageOptions_Stat_STATS_KIND_DEPENDENCIES:
				return "dependencies"
			case MessageOptions_Stat_STATS_KIND_MCV:
				return "mcv"
			default:
				panic("MessageOptions_Stat_STATS_KIND_UNSPECIFIED found on " + st.Name)
			}
		}), ","))

		_, _ = buf.WriteString(")")
	}
	_, _ = buf.WriteString(" ON ")
	_, _ = buf.WriteString(strings.Join(slice.Convert(st.Columns, func(in string) string {
		return `"` + in + `"`
	}), ","))
	_, _ = buf.WriteString(" FROM ")
	pgWriteString(buf, desc.TableName())
	_, _ = buf.WriteString("\n")
	return buf.String()
}

func pgWriteString(buf *bytes.Buffer, input string) {
	_, _ = buf.WriteString(`"`)
	// TODO(pquerna): not completely correct escaping
	_, _ = buf.WriteString(input)
	_, _ = buf.WriteString(`"`)
}

// autovacuum2with generates a WITH clause for PostgreSQL storage parameters.
// Returns empty string if no autovacuum options are configured.
func autovacuum2with(desc Descriptor) string {
	av := desc.GetAutovacuum()
	if av == nil {
		return ""
	}

	params := make([]string, 0)

	if av.HasVacuumThreshold() {
		params = append(params, "autovacuum_vacuum_threshold = "+formatInt32(av.GetVacuumThreshold()))
	}
	if av.HasVacuumScaleFactor() {
		params = append(params, "autovacuum_vacuum_scale_factor = "+formatFloat32(av.GetVacuumScaleFactor()))
	}
	if av.HasAnalyzeThreshold() {
		params = append(params, "autovacuum_analyze_threshold = "+formatInt32(av.GetAnalyzeThreshold()))
	}
	if av.HasAnalyzeScaleFactor() {
		params = append(params, "autovacuum_analyze_scale_factor = "+formatFloat32(av.GetAnalyzeScaleFactor()))
	}
	if av.HasVacuumCostDelay() {
		params = append(params, "autovacuum_vacuum_cost_delay = "+formatInt32(av.GetVacuumCostDelay()))
	}
	if av.HasVacuumCostLimit() {
		params = append(params, "autovacuum_vacuum_cost_limit = "+formatInt32(av.GetVacuumCostLimit()))
	}
	if av.HasFreezeMinAge() {
		params = append(params, "autovacuum_freeze_min_age = "+formatInt64(av.GetFreezeMinAge()))
	}
	if av.HasFreezeMaxAge() {
		params = append(params, "autovacuum_freeze_max_age = "+formatInt64(av.GetFreezeMaxAge()))
	}
	if av.HasFreezeTableAge() {
		params = append(params, "autovacuum_freeze_table_age = "+formatInt64(av.GetFreezeTableAge()))
	}
	if av.HasFillfactor() {
		params = append(params, "fillfactor = "+formatInt32(av.GetFillfactor()))
	}
	if av.HasToastTupleTarget() {
		params = append(params, "toast_tuple_target = "+formatInt32(av.GetToastTupleTarget()))
	}
	if av.HasEnabled() {
		params = append(params, "autovacuum_enabled = "+formatBool(av.GetEnabled()))
	}

	if len(params) == 0 {
		return ""
	}

	return "WITH (\n  " + strings.Join(params, ",\n  ") + "\n)"
}

func formatInt32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func formatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatFloat32(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func formatBool(v bool) string {
	return strconv.FormatBool(v)
}

func col2alter(desc Descriptor, col *Column) string {
	buf := &bytes.Buffer{}

	_, _ = buf.WriteString("ALTER TABLE ")
	pgWriteString(buf, desc.TableName())
	_, _ = buf.WriteString("\n")
	_, _ = buf.WriteString("ADD COLUMN IF NOT EXISTS")
	_, _ = buf.WriteString("\n")
	_, _ = buf.WriteString(col2spec(col))
	return buf.String()
}

func col2spec(col *Column) string {
	sbuf := &bytes.Buffer{}
	_, _ = sbuf.WriteString("  ")
	pgWriteString(sbuf, col.Name)
	_, _ = sbuf.WriteString(" ")
	if col.OverrideExpression != "" {
		_, _ = sbuf.WriteString(col.OverrideExpression)
	} else {
		_, _ = sbuf.WriteString(col.Type)
		if !col.Nullable {
			_, _ = sbuf.WriteString(" NOT NULL")
		}
		if col.Default != "" {
			_, _ = sbuf.WriteString(" DEFAULT " + col.Default)
		}
		if col.Collation != "" {
			_, _ = sbuf.WriteString(" COLLATE ")
			pgWriteString(sbuf, col.Collation)
		}
	}
	return sbuf.String()
}

func ksuidColOverrideExpression(col *Column) string {
	sbuf := &bytes.Buffer{}
	_, _ = sbuf.WriteString(col.Type)
	if !col.Nullable {
		_, _ = sbuf.WriteString(" NOT NULL")
	}
	if col.Default != "" {
		_, _ = sbuf.WriteString(" DEFAULT " + col.Default)
	}
	_, _ = sbuf.WriteString(" COLLATE \"C\"")
	return sbuf.String()
}
