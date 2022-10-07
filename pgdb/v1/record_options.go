package v1

type RecordOption struct {
	Prefix         string
	DataFieldsOnly bool
}

type RecordOptionsFunc func(option *RecordOption)

func RecordOptionPrefix(prefix string) RecordOptionsFunc {
	return func(option *RecordOption) {
		option.Prefix = prefix
	}
}

func RecordOptionDataFieldsOnly() RecordOptionsFunc {
	return func(option *RecordOption) {
		option.DataFieldsOnly = true
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
