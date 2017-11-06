package goenc

import (
	"fmt"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
)

/*
type AwaitingRequest struct {
	Req    *Request
	Answer chan bool
}

var (
	nextReqId = 0
	reqs  = make(map[int]*Request)
)
*/

// Called from server side
func HandleClientConnection(p *p2p.Peer, kp *enc.KeyPair) error {
	// Send my public key
	if err := SendPubKey(p, kp); err != nil {
		return err
	}

	// Receive client's public key
	if err := RecvPubKey(p); err != nil {
		return err
	}

	if err := RecvAndDecryptSymKey(p, kp); err != nil {
		return err
	}

	req, err := RecvAndDecryptRequest(p)
	if err != nil {
		return err
	}
	fmt.Println(req)

	hmac, err := RecvFile("/tmp/recv", req.Size, p)
	if err != nil {
		return err
	}

	var sig []byte
	if err = p.Recv(&sig); err != nil {
		return err
	}

	if err = p.Public.VerifySignature(hmac, sig); err != nil {
		return err
	}

	return nil
}

// Called from client side
func HandleServerConnection(filepath string, p *p2p.Peer, kp *enc.KeyPair) error {
	// Receive server's public key
	if err := RecvPubKey(p); err != nil {
		return err
	}

	// Send my public key
	if err := SendPubKey(p, kp); err != nil {
		return err
	}

	// Send the encrypted symmetric key
	if err := GenerateAndSendSymKey(p); err != nil {
		return err
	}

	req, err := NewRequest(filepath)
	if err != nil {
		return err
	}

	fmt.Println(req)
	if err = req.EncryptAndSend(p); err != nil {
		return err
	}

	hmac, err := SendFile(filepath, p)
	if err != nil {
		return err
	}

	sig, err := kp.Private.Sign(hmac)
	if err != nil {
		return err
	}

	if err = p.Send(sig); err != nil {
		return err
	}

	return nil
}
