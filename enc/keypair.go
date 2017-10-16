package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type KeyPair struct {
	Public  []byte
	Private []byte
}

func GenerateKeyPair() KeyPair {
	// Private key generation

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	err = priv.Validate()
	if err != nil {
		panic(err)
	}

	privDer := x509.MarshalPKCS1PrivateKey(priv)
	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}

	// Public key generation

	pub := priv.PublicKey

	pubDer, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		panic(err)
	}

	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}

	return KeyPair{
		Private: pem.EncodeToMemory(&privBlk),
		Public:  pem.EncodeToMemory(&pubBlk),
	}
}
