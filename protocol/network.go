package protocol

import (
	"encoding/binary"
	"unsafe"
)

const networkMessageMaxSize = 24590
const headerLength = 2
const checksumLength = 4
const xteaMultiple = 8
const maxBodyLength = networkMessageMaxSize - headerLength - checksumLength - xteaMultiple
const maxProtocolBodyLength = maxBodyLength - 10
const initialBufferPosition = 8

// NetworkMessage .
type NetworkMessage struct {
	buffer   []byte
	position int
	length   int
}

func netNetworkMessage(buffer []byte) *NetworkMessage {
	return &NetworkMessage{
		buffer:   buffer,
		position: initialBufferPosition,
		length:   len(buffer),
	}
}

func (msg *NetworkMessage) SkipBytes(count int) {
	msg.position += count
}

func (msg *NetworkMessage) GetUint16() uint16 {
	var rs uint16
	size := int(unsafe.Sizeof(rs))
	if !msg.canRead(size) {
		return 0
	}
	return binary.BigEndian.Uint16(msg.buffer[msg.position : msg.position+size])
}

func (msg NetworkMessage) canAdd(size int) bool {
	return size+msg.position < maxBodyLength
}

func (msg NetworkMessage) canRead(size int) bool {
	if msg.position+size > msg.length+8 || size >= networkMessageMaxSize-msg.position {
		return false
	}
	return true
}
