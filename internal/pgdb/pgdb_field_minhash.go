package pgdb

type pbMinHashConvert struct {
	VarName       string
	GoName        string
	ByteArrayName string
	ByteArraySize int32
}

func (pbdc *pbMinHashConvert) CodeForValue() (string, error) {
	// Make template
	return templateExecToString("field_minhash.tmpl", pbdc)
}

func (pbdc *pbMinHashConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}

func (pbdc *pbMinHashConvert) VarForAppend() (string, error) {
	return "", nil
}

func (pbdc *pbMinHashConvert) GoType() (string, error) {
	return "[]byte", nil
}

func (pbdc *pbMinHashConvert) EnumForValue() (string, error) {
	return "", nil
}
