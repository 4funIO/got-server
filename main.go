package main

import (
	"log"
	"net"
	"os"

	"github.com/4funIO/got-server/protocol"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("Usage: ./program [proxy from] [proxy to] [protocol (login, game)]")
		os.Exit(1)
	}
	if os.Args[3] != "login" && os.Args[3] != "game" {
		log.Println("Protocol param must be login or game")
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	for {
		log.Println("Waiting new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		var p protocol.Protocol
		if os.Args[3] == "login" {
			p = protocol.NewLogin()
		} else {
			p = protocol.NewGame()
		}

		go protocol.NewConnection(conn, os.Args[2]).Listen(p)
	}
}
