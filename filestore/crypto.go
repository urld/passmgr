// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filestore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

// aesCipher provides methods to encrypt and decrypt arbitrary content.
// The returned byte slice of each operation is guaranteed to be a valid
// input for the opposite operation.
type aesCipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type aesGcm struct {
	cipher.AEAD
	nonce []byte
}

func newGCM(key []byte) (aesCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCMWithNonceSize(block, nonceLenV1)
	if err != nil {
		return nil, err
	}

	return &aesGcm{AEAD: aead, nonce: nil}, nil
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
		if err := c.incrementNonce(); err != nil {
			return nil, err
		}
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
	if c.nonce == nil {
		c.nonce = make([]byte, nonceSize)
	}
	copy(c.nonce, ciphertext[:nonceSize])

	// decrypt:
	return c.AEAD.Open(ciphertext[:0], c.nonce, ciphertext[nonceSize:], nil)
}

func (c *aesGcm) incrementNonce() error {
	// increment the counter part of the nonce:
	counter := binary.BigEndian.Uint64(c.nonce[4:])
	counter++
	binary.BigEndian.PutUint64(c.nonce[4:], counter)
	// Change the random part of the nonce too, to avoid nonce reuse.
	// Nonce reuse could occur if a previous version of the store file
	// is checked out or restored using version control software like git
	// or other backup mechanisms.
	_, err := io.ReadFull(rand.Reader, c.nonce[:4])
	return err

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
