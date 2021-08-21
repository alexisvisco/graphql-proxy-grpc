package virtual

import (
	changecase "github.com/ku/go-change-case"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
)

type Message struct {
	Descriptor *Descriptor `json:"-"`
	Name       Name

	Fields []*Field

	IsInput bool
}

func (msg Message) String() string {
	return msg.Name.String()
}

func (msg Message) IsEmpty() bool {
	return len(msg.Fields) == 0
}

func (pkg *Package) getOrCreateMessage(message *protogen.Message, descriptor *Descriptor) (*Message, error) {
	// since we are using proto 3 => skip this entry
	// map entry was designed for proto 2 because map wasn't supported
	if message.Desc.IsMapEntry() {
		return nil, nil
	}

	n := Name{
		Identifier: message.GoIdent,
		GqlName:    changecase.Pascal(message.GoIdent.GoName),
	}

	if msg := pkg.getMessage(n); msg != nil {
		return msg, nil
	}

	vm := &Message{
		Descriptor: descriptor,
		Name:       n,
		Fields:     nil,
	}

	pkg.addMessage(vm)
	log.Debug().Msgf("register message %s in package %s", message.Desc.Name(), message.Desc.ParentFile().Package())

	for _, enum := range message.Enums {
		pkg.getOrCreateEnum(enum)
	}

	for _, msg := range message.Messages {
		_, err := pkg.getOrCreateMessage(msg, descriptor)
		if err != nil {
			return nil, err
		}
	}

	for _, field := range message.Fields {
		f, err := createField(field, descriptor)
		if err != nil {
			return nil, err
		}
		vm.Fields = append(vm.Fields, f)
	}

	return vm, nil
}

func (msg *Message) markAsInput() {
	msg.IsInput = true
	for _, field := range msg.Fields {
		if field.Type.Message != nil {
			field.Type.Message.IsInput = true
		} else if field.Type.Kv != nil {
			field.Type.Kv.IsInput = true
		}
	}
}
