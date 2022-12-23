package v1

type RecordOption struct {
	Prefix string
}

type RecordOptionsFunc func(option *RecordOption)

func RecordOptionPrefix(prefix string) RecordOptionsFunc {
	return func(option *RecordOption) {
		option.Prefix = prefix
	}
}

func NewRecordOptions(opts []RecordOptionsFunc) *RecordOption {
	option := &RecordOption{
		Prefix: "pb$",
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func (r *RecordOption) ColumnName(in string) string {
	return r.Prefix + in
}

func (r *RecordOption) Chain(prefix string) []RecordOptionsFunc {
	return []RecordOptionsFunc{
		RecordOptionPrefix(r.Prefix + prefix),
	}
}
