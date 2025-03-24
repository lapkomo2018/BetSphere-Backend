package hash

import (
	"crypto/sha256"
	"fmt"
)

type Hasher interface {
	Hash(password string) string
	Compare(hashedPassword, password string) error
}

type Hash struct {
	salt []byte
}

func NewHasher(salt string) *Hash {
	return &Hash{
		salt: []byte(salt),
	}
}

func (h *Hash) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum(h.salt))
}

func (h *Hash) Compare(hashedPassword, password string) error {
	if hashedPassword != h.Hash(password) {
		return fmt.Errorf("invalid password")
	}

	return nil
}
