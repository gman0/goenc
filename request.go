package goenc

import (
	"errors"
	"github.com/gman0/goenc/p2p"
	"net/http"
	"os"
	"path/filepath"
)

type Request struct {
	Size int64
	Type string
	Name string
}

func NewRequest(file string) (*Request, error) {
	st, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	if !st.Mode().IsRegular() {
		return nil, errors.New("Source path is not a file")
	}

	mime, err := guessMimeType(file)
	if err != nil {
		return nil, err
	}

	return &Request{
		Size: st.Size(),
		Type: mime,
		Name: filepath.Base(file),
	}, nil
}

func (r *Request) EncryptAndSend(p *p2p.Peer) error {
	_, err := SymSend(p, r)
	return err
}

func RecvAndDecryptRequest(p *p2p.Peer) (*Request, error) {
	r := &Request{}
	if _, err := SymRecv(p, r); err != nil {
		return nil, err
	}

	return r, nil
}

func guessMimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}
