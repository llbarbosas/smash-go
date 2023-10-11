package main

import (
	"flag"
	"log"

	"github.com/llbarbosas/smash-go"
)

func main() {
	linkAddr := flag.String("link", "", "link node to a remote node")
	remoteBusAPIAddr := flag.String("remotebus", ":6060", "set remote bus API addr")
	managmentAPIAddr := flag.String("managment", ":6061", "set managment API addr")
	flag.Parse()

	node, err := smash.NewNode(smash.NodeConfig{
		BusAddr:          *remoteBusAPIAddr,
		ManagmentAPIAddr: *managmentAPIAddr,
	})

	if err != nil {
		log.Fatal(err)
	}

	if *linkAddr != "" {
		if _, err := node.Link(*linkAddr); err != nil {
			log.Fatal(err)
		}
	}

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
