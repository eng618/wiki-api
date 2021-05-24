package main

import (
	"github.com/ENG618/wiki-api/server"
)

func main() {
	// Create and run server
	s := server.Server{}
	s.Initialize()
	s.Run(":5000")
}
