package pgdb

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type messageTemplateContext struct {
	ReceiverType   string
	MessageType    string
	DescriptorType string
	Fields         []*fieldContext
}

func renderMessage(ctx pgsgo.Context, w io.Writer, in pgs.File, m pgs.Message, ix *importTracker) error {
	ix.PGDB_v1 = true
	ix.GoquExp = true
	ix.ProtobufProto = true
	c := &messageTemplateContext{
		ReceiverType:   "*" + m.Name().UpperCamelCase().String(),
		MessageType:    getMessageType(m),
		DescriptorType: getDescriptorType(m),
		Fields:         getMessageFields(m),
	}
	return templates["message.tmpl"].Execute(w, c)
}

func getMessageType(m pgs.Message) string {
	return "pgdbMessage" + m.Name().UpperCamelCase().String()
}

func getMessageFields(m pgs.Message) []*fieldContext {
	fields := m.Fields()
	rv := make([]*fieldContext, 0, len(fields))
	for _, field := range fields {
		rv = append(rv, getField(field))
	}
	return rv
}
