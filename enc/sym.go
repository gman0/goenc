package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const (
	SymKeyLength = 32
	NonceLength  = 12
)

func GenerateSymmetricKey() (*[SymKeyLength]byte, error) {
	key := new([SymKeyLength]byte)
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateNonce() (*[NonceLength]byte, error) {
	nonce := new([NonceLength]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func SymmetricEncrypt(secret *[SymKeyLength]byte, plaintext []byte) (ciphertext []byte, nonce *[NonceLength]byte, err error) {
	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce, err = GenerateNonce()
	if err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nonce[:], nonce[:], plaintext, nil)
	return ciphertext, nonce, nil
}

func SymmetricDecrypt(secret *[SymKeyLength]byte, ciphertext []byte) (plaintext []byte, nonce *[NonceLength]byte, err error) {
	if len(ciphertext) <= NonceLength {
		return nil, nil, errors.New("SymmetricDecrypt(): message too short")
	}

	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = new([NonceLength]byte)
	copy(nonce[:], ciphertext)

	plaintext, err = gcm.Open(nil, nonce[:], ciphertext[NonceLength:], nil)
	if err != nil {
		return nil, nil, err
	}

	return plaintext, nonce, nil
}
