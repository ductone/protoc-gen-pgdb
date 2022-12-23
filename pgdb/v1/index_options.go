package v1

type IndexOption struct {
	IndexPrefix  string
	ColumnPrefix string
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

func NewIndexOptions(opts []IndexOptionsFunc) *IndexOption {
	option := &IndexOption{
		IndexPrefix:  "pbidx_",
		ColumnPrefix: "pb$",
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
	return r.ColumnPrefix + in
}

func (r *IndexOption) Nested(prefix string) []IndexOptionsFunc {
	return []IndexOptionsFunc{
		IndexOptionIndexPrefix(r.IndexPrefix + prefix),
		IndexOptionColumnPrefix(r.ColumnPrefix + prefix),
	}
}
