package data

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

type Character struct {
	ID   uuid.UUID
	Pos  Point
	Conn *websocket.Conn
	Send chan []byte
}
