package enc

import (
	"crypto/sha256"
	"io"
	"os"
)

func FileChecksum(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
