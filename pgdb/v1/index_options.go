package v1

type IndexOption struct {
	IndexPrefix  string
	ColumnPrefix string
	IsNested     bool
}

type IndexOptionsFunc func(option *IndexOption)

func IndexOptionIndexPrefix(prefix string) IndexOptionsFunc {
	return func(option *IndexOption) {
		option.IndexPrefix = prefix
	}
}
func IndexOptionColumnPrefix(prefix string) IndexOptionsFunc {
	return func(option *IndexOption) {
		option.ColumnPrefix = prefix
	}
}

func IndexOptionIsNested(b bool) IndexOptionsFunc {
	return func(option *IndexOption) {
		option.IsNested = b
	}
}

const columnPrefix = "pb$"

func NewIndexOptions(opts []IndexOptionsFunc) *IndexOption {
	option := &IndexOption{
		IndexPrefix:  "pbidx_",
		ColumnPrefix: columnPrefix,
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func (r *IndexOption) IndexName(in string) string {
	return r.IndexPrefix + in
}

func (r *IndexOption) ColumnName(in string) string {
	// Tenant IDs are not reflected into ancestor messages
	//  indices from nested messages that reference tenant_ids need to be retargeted to the top message
	if in == "tenant_id" {
		return columnPrefix + "tenant_id"
	}
	return r.ColumnPrefix + in
}

func (r *IndexOption) Nested(prefix string) []IndexOptionsFunc {
	return []IndexOptionsFunc{
		IndexOptionIndexPrefix(r.IndexPrefix + prefix),
		IndexOptionColumnPrefix(r.ColumnPrefix + prefix),
		IndexOptionIsNested(true),
	}
}
