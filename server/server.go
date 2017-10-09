package server

import (
	"log"
	"net/http"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/vloek/realm/data"
	"github.com/vloek/realm/data/messages"
)

type Server struct {
	characters []*data.Character
	conn       *gosocketio.Server
}

func NewServer() *Server {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	return &Server{
		characters: make([]*data.Character, 0),
		conn:       server,
	}
}

func (s *Server) Run() {
	go s.watchCharacters()
	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/", s.conn)

	http.ListenAndServe(":8032", serveMux)
}

func (s *Server) Initialize() *Server {
	s.conn.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to room
		c.Join("game")
	})

	//handle custom event
	s.conn.On("login", func(c *gosocketio.Channel, lm messages.LoginMessage) string {
		//send event to all in room

		if lm.IsValid() {
			log.Println("Char valid and add")
			char := data.NewCharacter(c)

			s.characters = append(s.characters, char)

			c.BroadcastTo("game", "character", char)
		}

		return "OK"
	})

	return s
}

func (s *Server) watchCharacters() {
	tick := time.Tick(5000 * time.Millisecond)
	for {
		select {
		case <-tick:
			for ch := range s.characters {
				log.Println("Ch: %s", ch)
			}
		}
	}
}
