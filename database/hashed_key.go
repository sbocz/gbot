package database

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"math/rand"
)

const (
	HashedKeySize int = sha256.Size
)

type HashedKey []byte

func NewRandomHashedKey() HashedKey {
	token := make([]byte, HashedKeySize)
	rand.Read(token)
	return token
}

func NewHashedKey(value interface{}) HashedKey {
	sum := sha256.Sum256(getBytes(value))
	return sum[:]
}

func getBytes(key interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return make([]byte, 0)
	}
	return buf.Bytes()
}
