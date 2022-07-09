package main

import (
	"github.com/ensiouel/gocket"
	"log"
	"net/http"
)

type Test struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	server := gocket.NewServer()

	server.OnConnecting(func(socket *gocket.Socket) {
		socket.On("test", func(d gocket.EventData) {
			var t Test

			if err := d.Unmarshal(&t); err != nil {
				log.Println(err)
				return
			}

			socket.Join("test")

			socket.To("test").Emit("h", "3")
		})

		log.Println(gocket.EventConnecting, socket)
	})

	server.OnDisconnecting(func(socket *gocket.Socket) {})

	http.Handle("/ws", server)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
