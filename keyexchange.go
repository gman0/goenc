package goenc

import (
	"encoding/gob"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
)

func SendPubKey(p *p2p.Peer, selfKp *enc.KeyPair) error {
	return gob.NewEncoder(p.Conn).Encode(selfKp.Public)
}

func RecvPubKey(p *p2p.Peer) error {
	dec := gob.NewDecoder(p.Conn)
	return dec.Decode(&p.Public)
}

func GenerateAndSendSymKey(p *p2p.Peer, kp *enc.KeyPair) error {
	// Generate shared key
	sk, err := enc.GenerateSymmetricKey()
	if err != nil {
		return err
	}

	if err = AsymSendByteSlice(p, kp, sk[:]); err != nil {
		return err
	}

	p.Sk = sk

	return nil
}

func RecvAndDecryptSymKey(p *p2p.Peer, kp *enc.KeyPair) error {
	sk, err := AsymRecvByteSlice(p, kp)
	if err != nil {
		return err
	}

	p.Sk = new([enc.SymKeyLength]byte)
	copy(p.Sk[:], sk)

	return nil
}
