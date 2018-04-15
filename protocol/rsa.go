package protocol

import (
	"log"
	"math/big"
)

type RSA interface {
	DecryptNetworkMessage(msg *NetworkMessage) bool
}

func NewRSA(pString, qString string) RSA {
	p, valid := big.NewInt(0).SetString(pString, 10)
	if !valid {
		log.Panicf("Invalid number: %s", pString)
	}
	q, valid := big.NewInt(0).SetString(qString, 10)
	if !valid {
		log.Panicf("Invalid number: %s", qString)
	}

	e := big.NewInt(65537)
	n := big.NewInt(0).Mul(p, q)
	p1 := big.NewInt(0).Sub(p, big.NewInt(1))
	q1 := big.NewInt(0).Sub(q, big.NewInt(1))
	pq1 := big.NewInt(0).Mul(p1, q1)
	d := big.NewInt(0).ModInverse(e, pq1)

	return &rsa{n, d}
}

type rsa struct {
	n, d *big.Int
}

func (r *rsa) DecryptNetworkMessage(msg *NetworkMessage) bool {
	if msg.length-msg.position < 128 {
		return false
	}

	block := msg.buffer[msg.position : msg.position+128]
	block = r.decrypt(block)
	msg.replaceAtPosition(block)

	return msg.GetByte() == 0
}

func (d *rsa) decrypt(msg []byte) []byte {
	c := big.NewInt(0).SetBytes(msg[:128])
	m := big.NewInt(0).Exp(c, d.d, d.n)
	rs := m.Bytes()
	diff := 128 - len(rs)
	for ; diff > 0; diff-- {
		rs = append([]byte{0}, rs...)
	}
	return rs
}
