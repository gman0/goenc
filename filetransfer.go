package goenc

import (
	//"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	//"fmt"
	"github.com/gman0/goenc/enc"
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

		ciphertext, nonce, err := enc.SymmetricEncrypt(p.Sk, buf[:])
		if err != nil {
			return nil, err
		}

		err = p.ConnEnc.Encode(&ciphertext)
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
		buf  []byte
	)

	for sz < totalSz {
		err := p.ConnDec.Decode(&buf)
		if err != nil {
			return nil, err
		}

		plaintext, nonce, err := enc.SymmetricDecrypt(p.Sk, buf[:])
		if err != nil {
			return nil, err
		}

		seq++
		if binary.BigEndian.Uint32(plaintext[:4]) != seq {
			return nil, errors.New("Received wrong sequence number")
		}

		var n int64 = bufSize
		if sz+bufSize > totalSz {
			n = totalSz - sz
		}

		_, err = f.Write(plaintext[4 : n+4])
		if err != nil {
			return nil, err
		}

		sz += n
		hmac.Write(nonce[:])
	}

	return hmac.Sum(nil), nil
}
