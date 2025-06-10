package pgdb

type pbByteArrayConvert struct {
	VarName       string
	GoName        string
	ByteArrayName string
	ByteArraySize int32
}

func (pbdc *pbByteArrayConvert) CodeForValue() (string, error) {
	// Make template
	return templateExecToString("field_byte_array.tmpl", pbdc)
}

func (pbdc *pbByteArrayConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}

func (pbdc *pbByteArrayConvert) VarForAppend() (string, error) {
	return "", nil
}

func (pbdc *pbByteArrayConvert) GoType() (string, error) {
	return "[]byte", nil
}

func (pbdc *pbByteArrayConvert) EnumForValue() (string, error) {
	return "", nil
}
