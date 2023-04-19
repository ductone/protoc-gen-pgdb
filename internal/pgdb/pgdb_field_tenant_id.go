package pgdb

import (
	"fmt"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type tenantIdDataConvert struct {
	ctx     pgsgo.Context
	VarName string
	Message pgs.Message
}

func (tidc *tenantIdDataConvert) GoType() (string, error) {
	return "string", nil
}

func (tidc *tenantIdDataConvert) CodeForValue() (string, error) {
	fieldName, err := getTenantIDField(tidc.Message)
	if err != nil {
		return "", err
	}
	fieldRef := fieldByName(tidc.Message, fieldName)
	// no-op, we just needed to calculate our variable name!
	return templateExecToString("proto_format_cast.tmpl", &formatContext{
		VarName:   tidc.VarName,
		InputName: "m.self." + tidc.ctx.Name(fieldRef).String(),
		CastType:  "string",
		IsArray:   false,
	})
}

func (tidc *tenantIdDataConvert) VarForValue() (string, error) {
	return tidc.VarName, nil
}

func (tidc *tenantIdDataConvert) VarForAppend() (string, error) {
	return "", nil
}

func (tidc *tenantIdDataConvert) EnumForValue() (string, error) {
	return "", nil
}

func getTenantIDField(msg pgs.Message) (string, error) {
	fieldName := "tenant_id"
	ext := pgdb_v1.MessageOptions{}
	_, err := msg.Extension(pgdb_v1.E_Msg, &ext)
	if err != nil {
		return "", fmt.Errorf("pgdb: getTenantIDField: failed to extract Message extension from '%s': %w", msg.FullyQualifiedName(), err)
	}
	if ext.TenantIdField != "" {
		fieldName = ext.TenantIdField
	}
	// panics if tenant id not found
	_ = fieldByName(msg, fieldName)
	return fieldName, nil
}
