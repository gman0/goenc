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

func GenerateAndSendSymKey(p *p2p.Peer) ([]byte, error) {
	// Generate symmetric key
	symKey, err := enc.GenerateSymmetricKey()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key
	encSymKey, err := p.Public.Encrypt(symKey)
	if err != nil {
		return nil, err
	}

	return symKey, gob.NewEncoder(p.Conn).Encode(&encSymKey)
}

func RecvAndDecryptSymKey(p *p2p.Peer, selfKp *enc.KeyPair) ([]byte, error) {
	var encKey []byte
	dec := gob.NewDecoder(p.Conn)
	if err := dec.Decode(&encKey); err != nil {
		return nil, err
	}

	if key, err := selfKp.Private.Decrypt(encKey); err != nil {
		return nil, err
	} else {
		return key, nil
	}
}
