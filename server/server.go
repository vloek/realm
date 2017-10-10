package server

import (
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	uuid "github.com/satori/go.uuid"
	"github.com/vloek/realm/data"
	"github.com/vloek/realm/data/messages"
)

type Message struct {
	Hello string
}

type Server struct {
	characters map[uuid.UUID]*data.Character
	conn       *socketio.Server
}

func NewServer() *Server {
	server, _ := socketio.NewServer(nil)

	return &Server{
		characters: make(map[uuid.UUID]*data.Character),
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
	// Connect
	s.conn.On("connection", s.connection)

	// Disconnect
	s.conn.On("disconnection", s.disconnection)

	//handle custom event
	s.conn.On("login", s.loginDo)

	return s
}

func (s *Server) connection(c socketio.Socket) {
	log.Println("New client connected")
}

func (s *Server) disconnection(c socketio.Socket) {
	c.Leave("game")

	for _, ch := range s.characters {
		if ch.Conn == c {
			log.Println("Found!")
			delete(s.characters, ch.ID)
		}
	}

	log.Println("Disconnected")
}

func (s *Server) loginDo(c socketio.Socket, lm messages.LoginMessage) string {
	log.Println("Login..")

	if lm.IsValid() {
		c.Join("game")

		char := data.NewCharacter(c)

		s.characters[char.ID] = char

		c.BroadcastTo("game", "join", char.SerializeToBroadcast())
	}

	return "OK"
}

func (s *Server) watchCharacters() {
	tick := time.Tick(5000 * time.Millisecond)
	for {
		select {
		case <-tick:
			for k, ch := range s.characters {
				log.Println("Ch:", k, ch)
			}
		}
	}
}
