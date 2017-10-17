package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

const (
	PrivateKeyType = "RSA PRIVATE KEY"
	PublicKeyType  = "PUBLIC KEY"
	KeyLength      = 2048
)

type RSAPrivateKey struct {
	*rsa.PrivateKey
}

func (key *RSAPrivateKey) ToPEM() []byte {
	der := x509.MarshalPKCS1PrivateKey(key.PrivateKey)
	block := pem.Block{
		Type:    PrivateKeyType,
		Headers: nil,
		Bytes:   der,
	}

	return pem.EncodeToMemory(&block)
}

func ParsePrivateKeyFromPEM(data []byte) (RSAPrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return RSAPrivateKey{}, errors.New("Error while decoding a PEM private key")
	}

	if block.Type != PrivateKeyType {
		return RSAPrivateKey{}, errors.New(fmt.Sprintf(
			"Cannot parse PEM private key of unsupported type %s", block.Type))
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return RSAPrivateKey{}, err
	} else {
		return RSAPrivateKey{key}, nil
	}
}

type RSAPublicKey struct {
	*rsa.PublicKey
}

func (key *RSAPublicKey) ToPEM() ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key.PublicKey)
	if err != nil {
		return nil, err
	}

	block := pem.Block{
		Type:    PublicKeyType,
		Headers: nil,
		Bytes:   der,
	}

	return pem.EncodeToMemory(&block), nil
}

func ParsePublicKeyFromPEM(data []byte) (RSAPublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return RSAPublicKey{}, errors.New("Error while decoding a PEM public key")
	}

	if block.Type != PublicKeyType {
		return RSAPublicKey{}, errors.New(fmt.Sprintf(
			"Cannot parse PEM public key of unsupported type %s", block.Type))
	}

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return RSAPublicKey{}, err
	} else {
		return RSAPublicKey{key.(*rsa.PublicKey)}, nil
	}
}

type KeyPair struct {
	Public  RSAPublicKey
	Private RSAPrivateKey
}

func GenerateKeyPair() KeyPair {
	priv, err := rsa.GenerateKey(rand.Reader, KeyLength)
	if err != nil {
		panic(err)
	}

	err = priv.Validate()
	if err != nil {
		panic(err)
	}

	return KeyPair{
		Private: RSAPrivateKey{priv},
		Public:  RSAPublicKey{&priv.PublicKey},
	}
}
