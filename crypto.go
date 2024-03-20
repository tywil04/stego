package stego

import (
	"crypto/aes"
	"crypto/cipher"
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

func encrypt(key, iv, additionalData, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGcm.Seal(nil, iv, plainText, additionalData), nil
}

func decrypt(key, iv, additionalData, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGcm.Open(nil, iv, cipherText, additionalData)
}
