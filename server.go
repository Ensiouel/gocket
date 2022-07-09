package gocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	EventConnecting    = "connecting"
	EventDisconnecting = "disconnecting"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	room  *Room
	rooms map[string]*Room
}

func NewServer() *Server {
	server := new(Server)

	server.rooms = make(map[string]*Room)

	server.room = NewRoom("main")
	go server.room.run()

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	socket, err := NewSocket(conn, s)
	if err != nil {
		log.Println(err)
		return
	}

	s.room.Join(socket)

	go socket.read()
	go socket.send()

	if event, ok := s.room.events[EventConnecting]; ok {
		event(socket)
	}
}

func (s *Server) GetRoom(name string) *Room {
	return s.rooms[name]
}

func (s *Server) joinSocketToRoom(socket *Socket, name string) {
	room := s.GetRoom(name)

	if room == nil {
		room = NewRoom(name)
		s.rooms[name] = room

		go room.run()
	}

	room.Join(socket)
}

func (s *Server) OnConnecting(event func(*Socket)) {
	s.room.events[EventConnecting] = event
}

func (s *Server) OnDisconnecting(event func(*Socket)) {
	s.room.events[EventDisconnecting] = event
}
