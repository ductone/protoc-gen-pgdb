package pgdb

type pbVectorConvert struct {
	VarName        string
	EnumName       string
	GoName         string
	FloatArrayName string
	EnumModelValue string
	VectorType     string // "vector" or "halfvec"
}

func (pbdc *pbVectorConvert) CodeForValue() (string, error) {
	// Make template
	return templateExecToString("field_vector.tmpl", pbdc)
}

func (pbdc *pbVectorConvert) VarForValue() (string, error) {
	return pbdc.VarName, nil
}

func (pbdc *pbVectorConvert) VarForAppend() (string, error) {
	return "", nil
}

func (pbdc *pbVectorConvert) GoType() (string, error) {
	return "[]float", nil
}

func (pbdc *pbVectorConvert) EnumForValue() (string, error) {
	return "", nil
}
