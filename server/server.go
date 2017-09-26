package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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

type Server struct{}

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
		_, message, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}

		log.WithField("message", string(message)).Info("read")
	}
}
