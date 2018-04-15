package protocol

import (
	"fmt"
	"log"
)

type login struct {
	rsa RSA
}

func NewLogin(rsa RSA) Protocol {
	return &login{rsa}
}

func (l *login) ReceiveMessage(msg *NetworkMessage) error {
	// TODO(oliverkra): hard coded, must be fixed... missing some stuf of parsePacket (connection.cpp:177)
	msg.position = 7

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

	if !l.rsa.DecryptNetworkMessage(msg) {
		return ErrDisconnectUser{
			Message: fmt.Sprintf("Are u trying to hack me? :("),
			Version: version,
		}
	}

	key := make([]uint32, 4)
	key[0] = msg.GetUint32()
	key[1] = msg.GetUint32()
	key[2] = msg.GetUint32()
	key[3] = msg.GetUint32()

	log.Println(">>>>>>>KEY[0]: ", key[0])
	log.Println(">>>>>>>KEY[1]: ", key[1])
	log.Println(">>>>>>>KEY[2]: ", key[2])
	log.Println(">>>>>>>KEY[3]: ", key[3])

	return nil
}
