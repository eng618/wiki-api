package main

import (
	"github.com/ENG618/wiki-api/server"
)

func main() {
	// Create and run server
	s := server.Server{}
	s.Run(":5000")
}
