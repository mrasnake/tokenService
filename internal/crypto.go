package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"time"
)

func Encrypter(key, nonce, secret []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, secret, nil)
	return ciphertext, nil
}

func Decrypter(key, nonce []byte, secret string) ([]byte, error) {
	ciphertext, _ := hex.DecodeString(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("unable to create cypher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCM: %w", err)
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func keyGen() []byte {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *mrand.Rand = mrand.New(
		mrand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}
