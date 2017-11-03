package goenc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/gman0/goenc/enc"
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
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(r)

	ciphertext, _, err := enc.SymmetricEncrypt(p.Sk, buf.Bytes())
	if err != nil {
		return err
	}

	return gob.NewEncoder(p.Conn).Encode(&ciphertext)
}

func RecvAndDecryptRequest(p *p2p.Peer) (*Request, error) {
	var ciphertext []byte
	if err := gob.NewDecoder(p.Conn).Decode(&ciphertext); err != nil {
		return nil, err
	}

	plaintext, _, err := enc.SymmetricDecrypt(p.Sk, ciphertext)
	if err != nil {
		return nil, err
	}

	var r Request
	buf := bytes.NewBuffer(plaintext)
	if err = gob.NewDecoder(buf).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
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
