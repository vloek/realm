package data

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Character struct {
	ID   uuid.UUID
	Pos  Point
	Conn *websocket.Conn
	Send chan []byte
}
