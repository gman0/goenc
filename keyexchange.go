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

func GenerateAndSendSymKey(p *p2p.Peer) error {
	// Generate symmetric key
	symKey, err := enc.GenerateSymmetricKey()
	if err != nil {
		return err
	}

	// Encrypt the symmetric key
	ciphertext, err := p.Public.Encrypt(symKey[:])
	if err != nil {
		return err
	}

	if err := p.ConnEnc.Encode(&ciphertext); err != nil {
		return err
	}

	p.Sk = symKey
	return nil
}

func RecvAndDecryptSymKey(p *p2p.Peer, selfKp *enc.KeyPair) error {
	var ciphertext []byte
	if err := p.ConnDec.Decode(&ciphertext); err != nil {
		return err
	}

	if key, err := selfKp.Private.Decrypt(ciphertext); err != nil {
		return err
	} else {
		p.Sk = new([enc.SymKeyLength]byte)
		copy(p.Sk[:], key)

		return nil
	}
}
