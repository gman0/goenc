package goenc

import (
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"

	"fmt"
)

// Called from server side
func HandleClientConnection(p *p2p.Peer, kp *enc.KeyPair) error {
	// Receive client's public key
	if err := RecvPubKey(p); err != nil {
		return err
	}

	// Send my public key
	if err := SendPubKey(p, kp); err != nil {
		return err
	}

	// Send the encrypted symmetric key
	symKey, err := GenerateAndSendSymKey(p)
	if err != nil {
		return err
	}

	fmt.Println(symKey)

	return nil
}

// Called from client side
func HandleServerConnection(filepath string, p *p2p.Peer, kp *enc.KeyPair) error {
	// Send my public key
	if err := SendPubKey(p, kp); err != nil {
		return err
	}

	// Receive server's public key
	if err := RecvPubKey(p); err != nil {
		return err
	}

	symKey, err := RecvAndDecryptSymKey(p, kp)
	if err != nil {
		return err
	}

	fmt.Println(symKey)

	return nil
}
