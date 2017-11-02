package goenc

import (
	"errors"
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

/*
func (r *Request) Send(p *p2p.Peer) error {
	return gob.NewEncoder(p.Conn).Encode(r)
}

func RecvRequest(p *p2p.Peer) (*Request, error) {
	dec := gob.NewDecoder(p.Conn)
	//
}
*/

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
