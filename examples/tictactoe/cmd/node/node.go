package main

import (
	"log"

	"github.com/llbarbosas/smash-tictactoe"
)

func main() {
	node, err := smash.NewNode()

	if err != nil {
		log.Fatal(err)
	}

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
