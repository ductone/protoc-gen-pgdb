package v1

import "github.com/doug-martin/goqu/v9/exp"

func MarshalNestedRecord(msg DBReflectMessage, opts ...RecordOptionsFunc) (exp.Record, error) {
	ro := NewRecordOptions(opts)
	recs, err := msg.DBReflect().Record(opts...)
	if err != nil {
		return nil, err
	}
	rv := exp.Record{}
	for k, v := range recs {
		cname := ro.ColumnName(k)
		rv[cname] = v
	}
	return rv, nil

}
