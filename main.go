package main

import (
	"github.com/vloek/realm/server"
)

func main() {
	realm := server.NewServer()
	realm.Run()
}
