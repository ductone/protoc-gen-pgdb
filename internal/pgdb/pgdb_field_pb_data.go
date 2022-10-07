package pgdb

type pbDataConvert struct {
	VarName string
}

func (pbdc *pbDataConvert) CodeForValue() (string, error) {
	return templateExecToString("field_pb_data.tmpl", pbdc)
}

func (pbdc *pbDataConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}

func (tidc *pbDataConvert) VarForAppend() (string, error) {
	return "", nil
}

func (pbdc *pbDataConvert) GoType() (string, error) {
	return "[]byte", nil
}
