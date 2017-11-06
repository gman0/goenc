package goenc

import (
	"bytes"
	"encoding/gob"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
)

func SymSendByteSlice(p *p2p.Peer, data []byte) (*[enc.NonceLength]byte, error) {
	ciphertext, nonce, err := enc.SymmetricEncrypt(p.Sk, data)
	if err != nil {
		return nil, err
	}

	return nonce, p.Send(ciphertext)
}

func SymSend(p *p2p.Peer, data interface{}) (*[enc.NonceLength]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	return SymSendByteSlice(p, buf.Bytes())
}

func SymRecvByteSlice(p *p2p.Peer) ([]byte, *[enc.NonceLength]byte, error) {
	var ciphertext []byte
	if err := p.Recv(&ciphertext); err != nil {
		return nil, nil, err
	}

	return enc.SymmetricDecrypt(p.Sk, ciphertext)
}

func SymRecv(p *p2p.Peer, data interface{}) (*[enc.NonceLength]byte, error) {
	plaintext, nonce, err := SymRecvByteSlice(p)
	if err != nil {
		return nil, err
	}

	return nonce, gob.NewDecoder(bytes.NewBuffer(plaintext)).Decode(data)
}
