package v1

type RecordOption struct {
	Prefix   string
	IsNested bool
	Nulled   bool
}

type RecordOptionsFunc func(option *RecordOption)

func RecordOptionPrefix(prefix string) RecordOptionsFunc {
	return func(option *RecordOption) {
		option.Prefix = prefix
	}
}

func RecordOptionIsNested(b bool) RecordOptionsFunc {
	return func(option *RecordOption) {
		option.IsNested = b
	}
}

func RecordOptionNulled(b bool) RecordOptionsFunc {
	return func(option *RecordOption) {
		option.Nulled = b
	}
}

func NewRecordOptions(opts []RecordOptionsFunc) *RecordOption {
	option := &RecordOption{}
	for _, opt := range opts {
		opt(option)
	}
	if option.Prefix == "" {
		option.Prefix = columnPrefix
	}
	return option
}

func (r *RecordOption) ColumnName(in string) string {
	return r.Prefix + in
}

func (r *RecordOption) Nested(prefix string) []RecordOptionsFunc {
	return []RecordOptionsFunc{
		RecordOptionPrefix(r.Prefix + prefix),
		RecordOptionIsNested(true),
		RecordOptionNulled(r.Nulled),
	}
}
