package enc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

const (
	hashFn = crypto.SHA512
)

func (key *RSAPrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha512.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, key.PrivateKey, hashFn, digest)
}

func (key *RSAPublicKey) VerifySignature(message []byte, sig []byte) error {
	hash := sha512.New()
	hash.Write(message)
	digest := hash.Sum(nil)

	return rsa.VerifyPKCS1v15(key.PublicKey, hashFn, digest, sig)
}
