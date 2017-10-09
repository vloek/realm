package server

import (
	"log"
	"net/http"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/vloek/realm/data"
)

type Server struct {
	characters []*data.Character
	conn       *http.ServeMux
}

func NewServer() *Server {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to room
		c.Join("game")
	})

	//handle custom event
	server.On("login", func(c *gosocketio.Channel, character data.Character) string {
		//send event to all in room
		c.BroadcastTo("game", "message", character)
		return "OK"
	})

	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)

	return &Server{
		characters: make([]*data.Character, 1024),
		conn:       serveMux,
	}
}

func (s *Server) Run() {
	http.ListenAndServe(":8032", s.conn)
}
