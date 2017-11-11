package goenc

import (
	"errors"
	"fmt"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
	"os"
	"sync"
)

type Response struct {
	Accept bool
	Dest   string
}

type AwaitingRequest struct {
	Req      *Request
	RespChan chan Response
	Resp     Response
}

func (ar *AwaitingRequest) WaitForResponse() bool {
	ar.Resp = <-ar.RespChan
	return ar.Resp.Accept
}

var (
	mtx       = sync.Mutex{}
	nextReqId = 0
	reqs      = make(map[int]*AwaitingRequest)
)

func addAwaitingRequest(r *Request) (*AwaitingRequest, int) {
	mtx.Lock()

	ar := &AwaitingRequest{
		Req:      r,
		RespChan: make(chan Response),
	}

	id := nextReqId
	nextReqId++
	reqs[id] = ar
	mtx.Unlock()

	return ar, id
}

func GetAwaitingRequest(id int) (*AwaitingRequest, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if r, ok := reqs[id]; !ok {
		return nil, errors.New("Request not found")
	} else {
		return r, nil
	}
}

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

	ar, id := addAwaitingRequest(req)
	PrintRequestInfo(id, req, p)

	if !ar.WaitForResponse() {
		fmt.Printf("> #%d declined, closing connection\n", id)
		return nil
	}

	hmac, err := RecvFile(ar.Resp.Dest, req.Size, p)
	if err != nil {
		return err
	}

	var sig []byte
	if err = p.Recv(&sig); err != nil {
		return err
	}

	if err = p.Public.VerifySignature(hmac, sig); err != nil {
		os.Remove(ar.Resp.Dest)
		return err
	}

	fmt.Printf("> #%d finished, closing connection\n", id)

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

	fmt.Printf("> connection with %s (%s) closed\n", p.Conn.RemoteAddr(), filepath)

	return nil
}
