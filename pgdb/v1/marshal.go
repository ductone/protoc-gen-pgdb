package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

func MarshalNestedRecord(msg DBReflectMessage, opts ...RecordOptionsFunc) (exp.Record, error) {
	recs, err := msg.DBReflect().Record(opts...)
	if err != nil {
		return nil, err
	}
	return recs, nil
}
