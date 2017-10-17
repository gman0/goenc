package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func (key *RSAPublicKey) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), rand.Reader,
		key.PublicKey, plaintext, nil)
}

func (key *RSAPrivateKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), rand.Reader,
		key.PrivateKey, ciphertext, nil)
}
