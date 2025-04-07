package hash

import (
	"crypto/sha256"
	"fmt"
)

type Hasher interface {
	Hash(s string) string
	Compare(hashedString, s string) bool
}

type Hash struct {
	salt []byte
}

func NewHasher(salt string) *Hash {
	return &Hash{
		salt: []byte(salt),
	}
}

func (h *Hash) Hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))

	return fmt.Sprintf("%x", hash.Sum(h.salt))
}

func (h *Hash) Compare(hashedString, s string) bool {
	return hashedString == h.Hash(s)
}
