package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func GenerateSymmetricKey() ([]byte, error) {
	key := new([32]byte)
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		return nil, err
	}

	return key[:], nil
}

func SymmetricEncrypt(secret, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func SymmetricDecrypt(secret, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("SymmetricDecrypt: ciphertext block size is too small")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream works in-place if the both arguments are the same
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
