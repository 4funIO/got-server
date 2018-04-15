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

	rsa := protocol.NewRSA(
		"14299623962416399520070177382898895550795403345466153217470516082934737582776038882967213386204600674145392845853859217990626450972452084065728686565928113",
		"7630979195970404721891201847792002125535401292779123937207447574596692788513647179235335529307251350570728407373705564708871762033017096809910315212884101")

	for {
		log.Println("Waiting new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		var p protocol.Protocol
		if os.Args[3] == "login" {
			p = protocol.NewLogin(rsa)
		} else {
			p = protocol.NewGame()
		}

		go protocol.NewConnection(conn, os.Args[2]).Listen(p)
	}
}
