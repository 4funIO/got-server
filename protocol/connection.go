package protocol

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

type Connection struct {
	conn         net.Conn
	proxyAddress string
}

func NewConnection(conn net.Conn, proxyAddress string) *Connection {
	return &Connection{
		conn:         conn,
		proxyAddress: proxyAddress,
	}
}

func (c *Connection) Listen(p Protocol) {
	proxy, err := net.Dial("tcp", c.proxyAddress)
	if err != nil {
		log.Panicf("Can't dial proxied server on address %s: %v", c.proxyAddress, err)
	}

	g, _ := errgroup.WithContext(context.Background())

	// Receive from client, send to server
	g.Go(func() error {
		for {
			msg := make([]byte, networkMessageMaxSize)
			size, err := c.conn.Read(msg)
			if err != nil {
				if err == io.EOF {
					log.Println("[CLIENT-TO-SERVER] Recebido EOF")
					return nil
				}
				return fmt.Errorf("[CLIENT-TO-SERVER] Reading from client: %v", err)
			}
			msg = msg[:size]

			// TODO(oliverkra): must translate some errors, like: ErrDisconnectUser
			if err := p.ReceiveMessage(netNetworkMessage(msg)); err != nil {
				return err
			}

			if _, err := proxy.Write(msg); err != nil {
				return fmt.Errorf("[CLIENT-TO-SERVER] Sending to server: %v", err)
			}
			log.Printf("[CLIENT-TO-SERVER] Packets: %d", size)
		}
	})

	// Receive from server, send to client
	g.Go(func() error {
		for {
			msg := make([]byte, networkMessageMaxSize)
			size, err := proxy.Read(msg)
			if err != nil {
				if err == io.EOF {
					log.Println("[CLIENT-TO-SERVER] Recebido EOF")
					return nil
				}
				return fmt.Errorf("[SERVER-TO-CLIENT] Reading from server: %v", err)
			}

			size, err = c.conn.Write(msg[:size])
			if err != nil {
				return fmt.Errorf("[SERVER-TO-CLIENT] Sending to client: %v", err)
			}
			log.Printf("[SERVER-TO-CLIENT] Packets: %d", size)
		}
	})

	if err := g.Wait(); err != nil {
		log.Panic(err)
	}
	log.Println("CLIENT DISCONNECTED")
}
