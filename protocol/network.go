package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func (msg *NetworkMessage) GetByte() uint8 {
	if !msg.canRead(1) {
		return 0
	}
	v := msg.buffer[msg.position]
	msg.position++
	return uint8(v)
}

func (msg *NetworkMessage) GetUint16() uint16 {
	var rs uint16
	size := int(unsafe.Sizeof(rs))
	if msg.canRead(size) {
		rs = binary.LittleEndian.Uint16(msg.buffer[msg.position : msg.position+size])
		msg.position += size
	}
	return rs
}

func (msg *NetworkMessage) GetUint32() uint32 {
	var rs uint32
	size := int(unsafe.Sizeof(rs))
	if msg.canRead(size) {
		rs = binary.LittleEndian.Uint32(msg.buffer[msg.position : msg.position+size])
		msg.position += size
	}
	return rs
}

func (msg *NetworkMessage) GetString() string {
	size := int(msg.GetUint16())
	if !msg.canRead(size) {
		return ""
	}
	rs := string(msg.buffer[msg.position : msg.position+size])
	msg.position += size
	return rs
}

func (msg *NetworkMessage) GetCurrentBlock() []byte {
	nullPosition := bytes.IndexByte(msg.buffer[msg.position:], 0)
	return msg.buffer[msg.position-1 : msg.position+nullPosition]
}

func (msg *NetworkMessage) replaceAtPosition(block []byte) {
	for i, v := range block {
		msg.buffer[msg.position+i] = v
	}
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

func debugPacket(packet []byte) {
	for index, m := 0, []byte{}; index < len(packet); index++ {
		if packet[index] == 0 {
			if len(m) == 0 {
				fmt.Printf("[%03d - %03d]: %v\n", index-len(m), index, m)
			} else {
				fmt.Printf("[%03d - %03d]: %v\n", index-len(m), index-1, m)
			}
			m = []byte{}
			continue
		}
		m = append(m, packet[index])
	}
}
