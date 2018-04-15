package protocol

import (
	"fmt"
	"log"
)

type game struct {
	firstMessageReceived bool
}

func NewGame() Protocol {
	return &game{false}
}

func (l *game) ReceiveMessage(msg *NetworkMessage) error {
	if !l.firstMessageReceived {
		l.firstMessageReceived = true
		return l.receiveFirstMessage(msg)
	}

	return nil
}

func (l *game) receiveFirstMessage(msg *NetworkMessage) error {
	// skip client OS
	msg.SkipBytes(2)

	// client version
	version := msg.GetUint16()
	log.Printf("> New client connection using version: %d", version)

	/*
	 * Skipped bytes:
	 * 4 bytes: protocolVersion
	 * 12 bytes: dat, spr, pic signatures (4 bytes each)
	 * 1 byte: 0
	 */
	if version >= 971 {
		msg.SkipBytes(17)
	} else {
		msg.SkipBytes(12)
	}

	if version <= ClientVersionMin {
		return ErrDisconnectUser{
			Message: fmt.Sprintf("Only clients with protocol %s allowed!", ClientVersionStr),
			Version: version,
		}
	}

	return nil
}
