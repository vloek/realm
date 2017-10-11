package data

import (
	uuid "github.com/satori/go.uuid"

	"github.com/googollee/go-socket.io"
)

type Character struct {
	ID   uuid.UUID `json:"id"`
	Pos  Point     `json:"point"`
	Conn socketio.Socket
	Send chan []byte
}

type CharacterBroadcast struct {
	ID  uuid.UUID `json:"id"`
	Pos Point     `json:"point"`
}

func NewCharacter(conn socketio.Socket) *Character {
	return &Character{
		ID:   uuid.NewV4(),
		Conn: conn,
		Pos:  Point{0, 0},
	}
}

func (c *Character) SerializeToBroadcast() *CharacterBroadcast {
	return &CharacterBroadcast{
		ID:  c.ID,
		Pos: c.Pos,
	}
}
