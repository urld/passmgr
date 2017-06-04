// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passmgr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

func genSalt() ([]byte, error) {

	salt := make([]byte, saltLenV1)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil

}

// deriveKey derives a key for AES-256, using scrypt.
func deriveKey(passwd []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(passwd, salt, 32768, 8, 4, 32)
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCMWithNonceSize(block, nonceLenV1)
}

func encrypt(key []byte, plaintext []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	// setup nonce:
	nonceSize := gcm.NonceSize()
	ciphertext := make([]byte, nonceSize+len(plaintext)+gcm.Overhead())
	nonce := ciphertext[:nonceSize]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt:
	gcm.Seal(ciphertext[:nonceSize], nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	// extract nonce:
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce := ciphertext[:nonceSize]

	// decrypt:
	return gcm.Open(ciphertext[:0], nonce, ciphertext[nonceSize:], nil)
}
