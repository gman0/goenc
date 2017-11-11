package goenc

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/gman0/goenc/p2p"
	"io"
	"os"
)

const (
	// The resulting envelope buffer contains an uint32 counter along with bufSize data
	bufSize = 4 + 8192
)

func SendFile(filepath string, p *p2p.Peer) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var (
		buf  = new([bufSize]byte) // counter + data buffer
		hmac = sha256.New()       // total hash calculated from nonces
		seq  uint32               // sequence number
	)

	for {
		seq++
		binary.BigEndian.PutUint32(buf[:4], seq)
		_, err := f.Read(buf[4:])

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		nonce, err := SymSendByteSlice(p, buf[:])
		if err != nil {
			return nil, err
		}

		hmac.Write(nonce[:])
	}

	return hmac.Sum(nil), nil
}

func RecvFile(filepath string, totalSz int64, p *p2p.Peer) ([]byte, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	var (
		sz   int64
		hmac = sha256.New()
		seq  uint32
	)

	for sz < totalSz {
		plaintext, nonce, err := SymRecvByteSlice(p)
		if err != nil {
			return nil, err
		}

		seq++
		if binary.BigEndian.Uint32(plaintext[:4]) != seq {
			return nil, errors.New("Received wrong sequence number")
		}

		var n int64 = bufSize
		if sz+bufSize > totalSz {
			n = totalSz - sz + 4
		}

		_, err = f.Write(plaintext[4:n])
		if err != nil {
			return nil, err
		}

		sz += n - 4
		hmac.Write(nonce[:])
	}

	return hmac.Sum(nil), nil
}
