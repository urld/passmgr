// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filestore implements a secure passmgr.Store.
// The contents are stored encrypted in a single file.
package filestore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urld/passmgr"
)

// fileStore implements a Store whichs contents can be stored encrypted
// in a file.
type fileStore struct {
	aesCipher
	filename string
	salt     []byte
	subjects []passmgr.Subject
}

func (s *fileStore) List() []passmgr.Subject {
	l := make([]passmgr.Subject, len(s.subjects))
	for i, c := range s.subjects {
		l[i] = passmgr.Subject{URL: c.URL, User: c.User}
	}
	return l
}

func (s *fileStore) Load(subject passmgr.Subject) (passmgr.Subject, bool) {
	for _, sc := range s.subjects {
		if sc.URL == subject.URL && sc.User == subject.User {
			return sc, true
		}
	}
	return subject, false
}

func (s *fileStore) Store(newSubject passmgr.Subject) {
	for i, c := range s.subjects {
		if c.URL == newSubject.URL && c.User == newSubject.User {
			s.subjects[i] = newSubject
			return
		}
	}
	s.subjects = append(s.subjects, newSubject)
}

func (s *fileStore) Delete(subject passmgr.Subject) bool {
	for i, c := range s.subjects {
		if c.URL == subject.URL && c.User == subject.User {
			s.subjects = append(s.subjects[:i], s.subjects[i+1:]...)
			return true
		}
	}
	return false
}

func (s *fileStore) persist() error {
	content, err := marshalStore(s.subjects)
	if err != nil {
		return err
	}
	content, err = s.aesCipher.Encrypt(content)
	if err != nil {
		return err
	}
	content = append(s.salt, content...)
	content = append(magicnumberV1, content...)

	// write to temp file first and rename to actual target file
	// to ensure write operation to be atomic, so the old file is not
	// lost on errors during write.
	tmpFileName := s.filename + fmt.Sprint(time.Now().Unix())
	err = os.WriteFile(tmpFileName, content, os.FileMode(0600))
	if err != nil {
		_ = os.Remove(tmpFileName)
		return err
	}
	return os.Rename(tmpFileName, s.filename)
}

// magic number to detect the version of the file format.
var magicnumberV1 = []byte{0x70, 0x61, 0x73, 0x73, 0x6d, 0x67, 0x72, 0x01} // passmgr1
const (
	saltLenV1  = 32
	nonceLenV1 = 12
)

// Read reads and decrypts the encrypted contents of a Store from the
// specified file.
//
// The file starts with a magic number to identify the file format version.
// Currently only the format 'passmgr1' is supported.
//
//	 0      7 8                             39 40        51 52           n
//	+--------+--------------------------------+------------+-----    -----+
//	|passmgr1|              salt              |   nonce    |  ciphertext  |
//	+--------+--------------------------------+------------+-----    -----+
//
//	passmgr1:   70 61 73 73 6d 67 72 01
//	salt:       32 byte salt for scrypt key derivation
//	nonce:      12 byte random nonce for aes-gcm
//	ciphertext: aes256-gcm encrypted json encoded content of Store
func Read(filename, passphrase string) (passmgr.Store, error) {
	content, err := readSecretFile(filename)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(content, magicnumberV1) {
		return nil, fmt.Errorf("unknown file type")
	}
	magicnumberLen := len(magicnumberV1)

	content = content[magicnumberLen:]
	salt := content[:saltLenV1]
	key, err := deriveKey([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}
	cipher, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	if len(content) == saltLenV1 {
		return &fileStore{aesCipher: cipher, filename: filename, salt: salt}, nil
	}

	ciphertext := content[saltLenV1:]
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	subjects, err := unmarshalStore(plaintext)
	s := &fileStore{aesCipher: cipher, filename: filename, salt: salt, subjects: subjects}
	return s, err
}

// Write persists a Store to the file it was read from.
// This function will panic if store was not created by ReadFileStore.
func Write(store passmgr.Store) error {
	fStore := store.(*fileStore)
	return fStore.persist()
}

// ChangeKey generates a new salt and key for a Store.
// This function will panic if store was not created by ReadFileStore.
func ChangeKey(store passmgr.Store, newPassphrase string) error {
	fStore := store.(*fileStore)

	salt, err := genSalt()
	if err != nil {
		return err
	}
	fStore.salt = salt

	key, err := deriveKey([]byte(newPassphrase), fStore.salt)
	if err != nil {
		return err
	}
	cipher, err := newGCM(key)
	if err != nil {
		return err
	}

	fStore.aesCipher = cipher
	return nil
}

func readSecretFile(filename string) ([]byte, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		salt, err := genSalt()
		return append(magicnumberV1, salt...), err
	}
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filename)
}

func unmarshalStore(content []byte) ([]passmgr.Subject, error) {
	var subjects []passmgr.Subject
	err := json.Unmarshal(content, &subjects)
	return subjects, err
}

func marshalStore(subjects []passmgr.Subject) ([]byte, error) {
	return json.Marshal(&subjects)
}
