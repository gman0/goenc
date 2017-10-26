package p2p

import (
	"github.com/gman0/goenc/enc"
	"net"
)

type Peer struct {
	Public enc.RSAPublicKey
	Conn   net.Conn
}
