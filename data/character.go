package data

import (
	"github.com/graarh/golang-socketio"

	uuid "github.com/satori/go.uuid"
)

type Character struct {
	ID   uuid.UUID
	Pos  Point `json:"point"`
	Conn *gosocketio.Channel
	Send chan []byte
}

func NewCharacter(conn *gosocketio.Channel) *Character {
	return &Character{
		Pos: Point{0, 0},
	}
}
