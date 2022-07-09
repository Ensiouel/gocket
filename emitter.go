package gocket

import (
	"encoding/json"
	"log"
)

const (
	EmitTypeBroadcast = EmitType(iota)
	EmitTypeBroadcastExceptSender
)

type EmitType int

type Emitter struct {
	emitType EmitType
	sender   *Socket
	sockets  map[*Socket]bool
}

func newEmptyEmitter() *Emitter {
	return &Emitter{}
}

func NewEmitter(t EmitType, sender *Socket, sockets map[*Socket]bool) *Emitter {
	emitter := newEmptyEmitter()

	emitter.emitType = t
	emitter.sender = sender
	emitter.sockets = sockets

	return emitter
}

func marshalEventRequest(name string, data interface{}) ([]byte, error) {
	request := EventRequest{
		Type: EventTypeEmit,
		Data: data,
		Name: name,
	}

	return json.Marshal(&request)
}

func (e *Emitter) Emit(name string, data interface{}) {
	b, err := marshalEventRequest(name, data)
	if err != nil {
		log.Println(err)
		return
	}

	for socket := range e.sockets {
		if e.emitType == EmitTypeBroadcastExceptSender && socket == e.sender {
			continue
		}

		socket.sendBuffer <- b
	}
}
