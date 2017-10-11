package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	uuid "github.com/satori/go.uuid"
	"github.com/vloek/realm/data"
)

type Message struct {
	Hello string
}

type Server struct {
	characters map[uuid.UUID]*data.Character `json:"characters"`
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

	fmt.Println("Run on 8032 port")
	http.ListenAndServe(":8032", serveMux)
}

func (s *Server) Initialize() *Server {
	// Connect
	s.conn.On("connection", s.connection)

	// Disconnect
	s.conn.On("disconnection", s.disconnection)

	//handle custom event
	s.conn.On("newplayer", s.loginDo)

	// click
	s.conn.On("click", s.clicked)

	return s
}

func (s *Server) clicked(c socketio.Socket, pos map[string]float64) {
	log.Printf("%s", pos)
	fmt.Println("CLick to ", pos["x"], pos["y"], "\n")
	if char := s.findCharacterBySocket(c); char != nil {
		char.Pos = data.Point{X: pos["x"], Y: pos["y"]}
		s.conn.BroadcastTo("game", "move", char.SerializeToBroadcast())
	}
}

func (s *Server) connection(c socketio.Socket) {
	log.Println("New client connected")
}

func (s *Server) disconnection(c socketio.Socket) {
	c.Leave("game")

	char := s.findCharacterBySocket(c)
	if char == nil {
		return
	}
	delete(s.characters, char.ID)

	log.Println("Disconnected")
}

func (s *Server) findCharacterBySocket(c socketio.Socket) *data.Character {
	for _, ch := range s.characters {
		if ch.Conn == c {
			return ch
		}
	}
	return nil
}

func (s *Server) loginDo(c socketio.Socket, x map[string]interface{}) string {
	log.Println("Login..", x)

	c.Join("game")

	char := data.NewCharacter(c)

	random := rand.New(rand.NewSource(2))
	xx := float64(random.Intn(120))
	yy := float64(random.Intn(200))
	pos := data.Point{X: xx, Y: yy}
	char.Pos = pos

	s.characters[char.ID] = char

	var allplayers []*data.CharacterBroadcast
	for _, ch := range s.characters {
		allplayers = append(allplayers, ch.SerializeToBroadcast())
	}

	// c.BroadcastTo("game", "join", char.SerializeToBroadcast())
	s.conn.BroadcastTo("game", "newplayer", char.SerializeToBroadcast())

	c.Emit("allplayers", allplayers)

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
