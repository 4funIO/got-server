package protocol

import (
	"log"
	"time"
)

type logger struct {
	next Protocol
}

func WithLogger(p Protocol) Protocol {
	return &logger{p}
}

func (l *logger) ReceiveMessage(msg *NetworkMessage) (err error) {
	start := time.Now()
	defer func() {
		log.Printf(">> Response time: %v; Err: %v\n", time.Now().Sub(start), err)
	}()
	return l.next.ReceiveMessage(msg)
}
