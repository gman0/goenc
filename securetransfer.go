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

func AsymSendByteSlice(p *p2p.Peer, kp *enc.KeyPair, data []byte) error {
	ciphertext, err := p.Public.Encrypt(data)
	if err != nil {
		return err
	}

	sig, err := kp.Private.Sign(ciphertext)
	if err != nil {
		return err
	}

	if err = p.Send(ciphertext); err != nil {
		return err
	}

	if err = p.Send(sig); err != nil {
		return err
	}

	return nil
}

func AsymSend(p *p2p.Peer, kp *enc.KeyPair, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	return AsymSendByteSlice(p, kp, buf.Bytes())
}

func AsymRecvByteSlice(p *p2p.Peer, kp *enc.KeyPair) ([]byte, error) {
	var ciphertext []byte
	if err := p.Recv(&ciphertext); err != nil {
		return nil, err
	}

	var sig []byte
	if err := p.Recv(&sig); err != nil {
		return nil, err
	}

	if err := p.Public.VerifySignature(ciphertext, sig); err != nil {
		return nil, err
	}

	return kp.Private.Decrypt(ciphertext)
}

func AsymRecv(p *p2p.Peer, kp *enc.KeyPair, data interface{}) error {
	plaintext, err := AsymRecvByteSlice(p, kp)
	if err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(plaintext)).Decode(data)
}
