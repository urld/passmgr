// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passmgr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

// Cipher provides methods to encrypt and decrypt arbitrary content.
// The returned byte slice of each operation is guaranteed to be a valid
// input for the opposite operation.
type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type aesGcm struct {
	cipher.AEAD
	nonce []byte
}

func newGCM(key []byte) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCMWithNonceSize(block, nonceLenV1)
	if err != nil {
		return nil, err
	}

	cipher := &aesGcm{AEAD: aead, nonce: nil}
	return cipher, nil
}

func (c *aesGcm) Encrypt(plaintext []byte) ([]byte, error) {
	// setup nonce:
	nonceSize := c.AEAD.NonceSize()
	if c.nonce == nil {
		// initial random nonce generation
		c.nonce = make([]byte, nonceSize)
		if _, err := io.ReadFull(rand.Reader, c.nonce); err != nil {
			return nil, err
		}
	} else {
		c.incrementNonce()
	}

	ciphertext := make([]byte, nonceSize+len(plaintext)+c.AEAD.Overhead())
	copy(ciphertext[:nonceSize], c.nonce)

	// encrypt:
	c.AEAD.Seal(ciphertext[:nonceSize], c.nonce, plaintext, nil)
	return ciphertext, nil
}

func (c *aesGcm) Decrypt(ciphertext []byte) ([]byte, error) {
	// extract nonce:
	nonceSize := c.AEAD.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	c.nonce = ciphertext[:nonceSize]

	// decrypt:
	return c.AEAD.Open(ciphertext[:0], c.nonce, ciphertext[nonceSize:], nil)
}

func (c *aesGcm) incrementNonce() {
	counter := binary.BigEndian.Uint64(c.nonce[4:])
	counter++
	binary.BigEndian.PutUint64(c.nonce[4:], counter)
}

// genSalt generates a salt, which can be used for the deriveKey fuction.
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
