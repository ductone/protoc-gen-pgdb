package v1

type StatisticOption struct {
	StatsPrefix  string
	ColumnPrefix string
	IsNested     bool
}

type StatisticOptionsFunc func(option *StatisticOption)

func StatsOptionStatsPrefix(prefix string) StatisticOptionsFunc {
	return func(option *StatisticOption) {
		option.StatsPrefix = prefix
	}
}
func StatsOptionColumnPrefix(prefix string) StatisticOptionsFunc {
	return func(option *StatisticOption) {
		option.ColumnPrefix = prefix
	}
}

func StatsOptionIsNested(b bool) StatisticOptionsFunc {
	return func(option *StatisticOption) {
		option.IsNested = b
	}
}

func NewStatisticOption(opts []StatisticOptionsFunc) *StatisticOption {
	option := &StatisticOption{
		StatsPrefix:  "pbstats_",
		ColumnPrefix: columnPrefix,
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func (r *StatisticOption) StatsName(in string) string {
	return r.StatsPrefix + in
}

func (r *StatisticOption) ColumnName(in string) string {
	// Tenant IDs are not reflected into ancestor messages
	//  indices from nested messages that reference tenant_ids need to be retargeted to the top message
	if in == "tenant_id" {
		return columnPrefix + "tenant_id"
	}
	return r.ColumnPrefix + in
}

func (r *StatisticOption) Nested(prefix string) []StatisticOptionsFunc {
	return []StatisticOptionsFunc{
		StatsOptionStatsPrefix(r.StatsPrefix + prefix),
		StatsOptionColumnPrefix(r.ColumnPrefix + prefix),
		StatsOptionIsNested(true),
	}
}
