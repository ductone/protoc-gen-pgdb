package pgdb

type pbDataConvert struct {
	VarName string
}

func (pbdc *pbDataConvert) CodeForValue() (string, error) {
	return templateExecToString("field_pbdata.tmpl", pbdc)
}

func (pbdc *pbDataConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}
