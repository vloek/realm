package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vloek/realm/data"
	"github.com/vloek/realm/server"
)

func main() {
	c := data.Character{}
	c.Pos = data.Point{X: 1.0, Y: 0.0}

	server := server.Server{}
	server.Run()

	log.Info(c)
}
