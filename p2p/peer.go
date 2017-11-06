package p2p

import (
	"encoding/gob"
	"github.com/gman0/goenc/enc"
	"net"
)

type Peer struct {
	Public  enc.RSAPublicKey
	Sk      *[enc.SymKeyLength]byte
	Conn    net.Conn
	ConnEnc *gob.Encoder
	ConnDec *gob.Decoder
}

func (p *Peer) Send(val interface{}) error {
	return p.ConnEnc.Encode(val)
}

func (p *Peer) Recv(val interface{}) error {
	return p.ConnDec.Decode(val)
}
