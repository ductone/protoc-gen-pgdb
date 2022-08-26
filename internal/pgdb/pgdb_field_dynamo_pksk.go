package pgdb

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	dynamopb "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
)

type dynamoKeyType int64

const (
	DynamoKeyTypePartition dynamoKeyType = 1
	DynamoKeyTypeSort      dynamoKeyType = 2
)

type dynamoKeyDataConvert struct {
	ctx     pgsgo.Context
	VarName string
	Message pgs.Message
	KeyType dynamoKeyType
	Parts   []string
}

func (tidc *dynamoKeyDataConvert) GoType() (string, error) {
	return fieldConvertString, nil
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
	case DynamoKeyTypePartition:
		// prefix
		dkdc.Parts = append(dkdc.Parts,
			fmt.Sprintf(`"%s_%s"`,
				dkdc.Message.Package().ProtoName().LowerSnakeCase().String(), dkdc.Message.Name().LowerSnakeCase().String()),
		)
		keyFields = dynExt.Key[0].PkFields
	case DynamoKeyTypeSort:
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
		formatter, err := typeToString(field.Type().ProtoType(), "m.self."+dkdc.ctx.Name(field).String())
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
