package gocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

type Room struct {
	id        uuid.UUID
	name      string
	joining   chan *Socket
	leaving   chan *Socket
	sockets   map[*Socket]bool
	events    map[string]func(*Socket)
	running   bool
	broadcast chan []byte
}

func NewRoom(name string) *Room {
	room := new(Room)

	room.id = uuid.New()
	room.name = name
	room.joining = make(chan *Socket)
	room.leaving = make(chan *Socket)
	room.sockets = make(map[*Socket]bool)
	room.events = make(map[string]func(*Socket))
	room.broadcast = make(chan []byte)

	return room
}

func (r *Room) run() {
	defer func() {
		close(r.joining)
		close(r.leaving)
		close(r.broadcast)
	}()

	r.running = true

	for {
		select {
		case socket := <-r.joining:
			log.Printf("socket %d joining to room %s\n", socket.id.ID(), r.name)
			socket.room = r
			r.sockets[socket] = true
		case socket := <-r.leaving:
			log.Printf("socket %d leaving to room %s\n", socket.id.ID(), r.name)
			delete(r.sockets, socket)
		case message := <-r.broadcast:
			for socket := range r.sockets {
				select {
				case socket.sendBuffer <- message:
				default:
					close(socket.sendBuffer)
					delete(r.sockets, socket)
				}
			}
		}
	}
}

func (r *Room) Emit(name string, data interface{}) {
	request := EventRequest{
		Type: EventTypeEmit,
		Data: data,
		Name: name,
	}

	b, err := json.Marshal(&request)
	if err != nil {
		log.Println(err)
		return
	}

	r.broadcast <- b
}

func (r *Room) Join(s *Socket) {
	if s == nil || r.running == false {
		return
	}

	if s.room != nil {
		s.room.Leave(s)
	}

	r.joining <- s
}

func (r *Room) Leave(s *Socket) {
	if s == nil || r.running == false {
		return
	}

	r.leaving <- s
}

func (r *Room) Sockets() map[*Socket]bool {
	return r.sockets
}
