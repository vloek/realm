package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vloek/realm/data"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Server struct {
	// Registered clients.
	characters map[uuid.UUID]*data.Character

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *data.Character

	// Unregister requests from clients.
	unregister chan *data.Character
}

func NewServer() *Server {
	return &Server{
		broadcast:  make(chan []byte),
		register:   make(chan *data.Character),
		unregister: make(chan *data.Character),
		characters: make(map[uuid.UUID]*data.Character),
	}
}

func (s *Server) Run() {
	log.Info("Running..")
	s.listen("localhost:8090")
}

func (s *Server) listen(addr string) {
	http.HandleFunc("/ws", s.wsHandler)
	http.ListenAndServe(addr, nil)
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	for {
		select {
		case character := <-s.register:
			s.characters[uuid.NewV4()] = character
		case cTmp := <-s.unregister:
			if character, ok := s.characters[cTmp.ID]; ok {
				delete(s.characters, cTmp.ID)
				close(character.Send)
			}
		case _ = <-s.broadcast:
			for _ = range s.characters {
				// sync messages on characters
			}
		}
	}
}
