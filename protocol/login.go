package protocol

import (
	"fmt"
	"log"
)

type login struct{}

func NewLogin() Protocol {
	return &login{}
}

func (l *login) ReceiveMessage(msg *NetworkMessage) error {
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
