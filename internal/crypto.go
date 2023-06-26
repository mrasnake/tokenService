package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	mrand "math/rand"
	"time"
)

// Encrypter takes the key, auth nonce and secret and returns the encrypted secret as a byte array.
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

// Decrypter takes the key, auth nonce and secret and returns the decrypted secret as a byte array.
func Decrypter(key, nonce []byte, secret string) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("unable to create cypher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCM: %w", err)
	}

	return aesgcm.Open(nil, nonce, []byte(secret), nil)
}

// keyGen return a randomly generated key as a byte array.
func keyGen() []byte {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}