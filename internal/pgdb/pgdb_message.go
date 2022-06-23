package pgdb

import (
	"fmt"
	"io"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type messageTemplateContext struct {
	ReceiverType            string
	MessageType             string
	DescriptorType          string
	Fields                  []*fieldContext
	WantRecordStringBuilder bool
}

func renderMessage(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDB_v1 = true
	ix.GoquExp = true
	ix.ProtobufProto = true
	ix.Strings = true
	c := &messageTemplateContext{
		ReceiverType:            "*" + m.Name().UpperCamelCase().String(),
		MessageType:             getMessageType(m),
		DescriptorType:          getDescriptorType(m),
		Fields:                  getMessageFields(ctx, m, ix),
		WantRecordStringBuilder: true, // unconditionally used by pk/sk builder
	}
	return templates["message.tmpl"].Execute(w, c)
}

func getMessageType(m pgs.Message) string {
	return "pgdbMessage" + m.Name().UpperCamelCase().String()
}

type varNamer struct {
	prefix string
	offset int
}

func (fn *varNamer) Next() *varNamer {
	return &varNamer{offset: fn.offset + 1, prefix: fn.prefix}
}

func (fn *varNamer) String() string {
	return fmt.Sprintf("%s%d", fn.prefix, fn.offset)
}

func getMessageFields(ctx pgsgo.Context, m pgs.Message, ix *importTracker) []*fieldContext {
	fields := m.Fields()
	rv := make([]*fieldContext, 0, len(fields)+lenCommonFields)
	cf, err := getCommonFields(ctx, m)
	if err != nil {
		panic(err)
	}
	rv = append(rv, cf...)
	vn := &varNamer{prefix: "v", offset: 0}
	for _, field := range fields {
		vn = vn.Next()
		rv = append(rv, getField(ctx, field, vn, ix))
	}
	return rv
}
