package stego

import (
	"hash/fnv"
	"math/rand"

	"golang.org/x/crypto/argon2"
)

func predictableRandom(key []byte) *rand.Rand {
	hasher := fnv.New64()
	hasher.Write(key)
	return rand.New(rand.NewSource(int64(hasher.Sum64())))
}

func deriveKey(key string, salt string) []byte {
	const (
		time    = 2
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32 // for aes 256 bit
	)
	return argon2.IDKey([]byte(key), []byte(salt), time, memory, threads, keyLen)
}
