package gocket

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

var (
	ErrConnectionIsNil = fmt.Errorf("socket creation error: %s", "connection is nil")
	ErrServerIsNil     = fmt.Errorf("socket creation error: %s", "server is nil")
)

type Socket struct {
	id         uuid.UUID
	room       *Room
	conn       *websocket.Conn
	server     *Server
	events     map[string]EventFunc
	sendBuffer chan []byte
}

func NewSocket(conn *websocket.Conn, server *Server) (*Socket, error) {
	if conn == nil {
		return nil, ErrConnectionIsNil
	}

	if server == nil {
		return nil, ErrServerIsNil
	}

	socket := new(Socket)

	socket.id = uuid.New()
	socket.conn = conn
	socket.server = server
	socket.events = make(map[string]EventFunc)
	socket.sendBuffer = make(chan []byte)

	return socket, nil
}

func (s *Socket) On(event string, f EventFunc) {
	s.events[event] = f
}

func (s *Socket) read() {
	defer func() {
		s.close()
	}()

	for {
		var event EventResponse

		err := s.conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
				break
			}
			log.Println(err)
			break
		}

		if event.IsEmit() {
			eventFunc, ok := s.events[event.Name]
			if ok {
				eventFunc(EventData(event.Data))
			}
		}
	}
}

func (s *Socket) send() {
	defer func() {
		s.close()
	}()

	for {
		message, ok := <-s.sendBuffer
		if !ok {
			_ = s.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := s.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Socket) Join(name string) {
	if s.room.name == name {
		return
	}
	s.server.joinSocketToRoom(s, name)
}

func (s *Socket) Emit(name string, data interface{}) {
	b, err := marshalEventRequest(name, data)
	if err != nil {
		log.Println(err)
		return
	}

	s.sendBuffer <- b
}

func (s *Socket) To(name string) *Emitter {
	room := s.server.GetRoom(name)
	if room == nil {
		return newEmptyEmitter()
	}

	return NewEmitter(EmitTypeExceptSender, s, room.sockets)
}

func (s *Socket) close() {
	_ = s.conn.Close()
}

func (s *Socket) GetID() uuid.UUID {
	return s.id
}
