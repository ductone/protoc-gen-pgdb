package pgdb

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
	dynamopb "github.com/pquerna/protoc-gen-dynamo/dynamo"
)

type dynamoKeyType int64

const (
	DKT_PK    dynamoKeyType = 1
	DKT_SK    dynamoKeyType = 2
	DKT_GS1PK dynamoKeyType = 3
	DKT_GS1SK dynamoKeyType = 4
)

type dynamoKeyDataConvert struct {
	VarName string
	Message pgs.Message
	KeyType dynamoKeyType
	Parts   []string
}

func (dkdc *dynamoKeyDataConvert) CodeForValue() (string, error) {
	dynExt := dynamopb.DynamoMessageOptions{}
	ok, err := dkdc.Message.Extension(dynamopb.E_Msg, &dynExt)
	if err != nil {
		return "", err
	}
	if ok && dynExt.Disabled {
		panic(fmt.Errorf("pgdb: dynamoKeyDataConvert failed for: %v(%v): dynamo extension must not be disabled if pgdb is enabled",
			dkdc.Message.FullyQualifiedName(), dkdc.Message.Descriptor().GetName()))
	}
	dkdc.Message.Fields()
	if len(dynExt.Key) == 0 {
		panic(fmt.Errorf("pgdb: dynamoKeyDataConvert failed for: %v(%v): dynamo extension must contain keys pgdb is enabled",
			dkdc.Message.FullyQualifiedName(), dkdc.Message.Descriptor().GetName()))
	}
	dkdc.Parts = []string{}
	keyFields := []string{}
	switch dkdc.KeyType {
	case DKT_PK:
		// prefix
		dkdc.Parts = append(dkdc.Parts,
			fmt.Sprintf(`"%s_%s"`,
				dkdc.Message.Package().ProtoName().LowerSnakeCase().String(), dkdc.Message.Name().LowerSnakeCase().String()),
		)
		keyFields = dynExt.Key[0].PkFields
	case DKT_SK:
		if dynExt.Key[0].SkConst != "" {
			dkdc.Parts = append(dkdc.Parts,
				fmt.Sprintf(`"%s"`,
					dynExt.Key[0].SkConst),
			)
		} else {
			keyFields = dynExt.Key[0].SkFields
		}
	default:
		panic(fmt.Errorf("pgdb: dynamoKeyDataConvert failed for: %v(%v): invalid key type",
			dkdc.Message.FullyQualifiedName(), dkdc.Message.Descriptor().GetName()))
	}

	for _, fieldName := range keyFields {
		field := fieldByName(dkdc.Message, fieldName)
		formatter, err := typeToString(field.Type().ProtoType(), "m.self."+field.Name().UpperCamelCase().String())
		if err != nil {
			panic(err)
		}
		dkdc.Parts = append(dkdc.Parts,
			formatter,
		)
	}
	return templateExecToString("field_dynamo_pksk.tmpl", dkdc)
}

func (dkdc *dynamoKeyDataConvert) VarForValue() (string, error) {
	return dkdc.VarName, nil
}

func fieldByName(msg pgs.Message, name string) pgs.Field {
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == name {
			return f
		}
	}
	panic(fmt.Sprintf("Failed to find field %s on %s", name, msg.FullyQualifiedName()))
}
