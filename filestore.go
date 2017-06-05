// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// fileStore implements a Store whichs contents can be stored encrypted
// in a file.
type fileStore struct {
	filename string
	salt     []byte
	key      []byte
	subjects []Subject
}

func (s *fileStore) List() []Subject {
	l := make([]Subject, len(s.subjects))
	for i, c := range s.subjects {
		l[i] = Subject{URL: c.URL, User: c.User}
	}
	return l
}

func (s *fileStore) Load(subject Subject) (Subject, bool) {
	for _, sc := range s.subjects {
		if sc.URL == subject.URL && sc.User == subject.User {
			return sc, true
		}
	}
	return subject, false
}

func (s *fileStore) Store(newSubject Subject) {
	for i, c := range s.subjects {
		if c.URL == newSubject.URL && c.User == newSubject.User {
			s.subjects[i] = newSubject
			return
		}
	}
	s.subjects = append(s.subjects, newSubject)
}

func (s *fileStore) Delete(subject Subject) bool {
	for i, c := range s.subjects {
		if c.URL == subject.URL && c.User == subject.User {
			s.subjects = append(s.subjects[:i], s.subjects[i+1:]...)
			return true
		}
	}
	return false
}

func (s *fileStore) Persist() error {
	content, err := marshalStore(s.subjects)
	if err != nil {
		return err
	}
	content, err = encrypt(s.key, content)
	if err != nil {
		return err
	}
	content = append(s.salt, content...)
	content = append(magicnumberV1, content...)
	return ioutil.WriteFile(s.filename, content, os.FileMode(0600))
}

// magic number to detect the version of the file format.
var magicnumberV1 = []byte{0x70, 0x61, 0x73, 0x73, 0x6d, 0x67, 0x72, 0x01} // passmgr1
const (
	saltLenV1  = 32
	nonceLenV1 = 12
)

// ReadFileStore reads and decrypts the encrypted contents of a Store from the
// specified file.
//
// The file starts with a magic number to identify the file format version.
// Currently only the format 'passmgr1' is supported.
//
//    0      7 8                             39 40        51 52           n
//   +--------+--------------------------------+------------+-----    -----+
//   |passmgr1|              salt              |   nonce    |  ciphertext  |
//   +--------+--------------------------------+------------+-----    -----+
//
//   passmgr1:   70 61 73 73 6d 67 72 01
//   salt:       32 byte salt for scrypt key derivation
//   nonce:      12 byte random nonce for aes-gcm
//   ciphertext: aes256-gcm encrypted json encoded content of Store
func ReadFileStore(filename, passphrase string) (Store, error) {
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

	if len(content) == saltLenV1 {
		return &fileStore{filename: filename, salt: salt, key: key}, nil
	}

	ciphertext := content[saltLenV1:]
	plaintext, err := decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	subjects, err := unmarshalStore(plaintext)
	s := &fileStore{filename: filename, salt: salt, key: key, subjects: subjects}
	return s, err
}

// WriteFileStore persists a Store to the file it was read from.
// This function will panic if  store was not created by ReadFileStore.
func WriteFileStore(store Store) error {
	fStore := store.(*fileStore)
	return fStore.Persist()
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
	return ioutil.ReadFile(filename)
}

func unmarshalStore(content []byte) ([]Subject, error) {
	var subjects []Subject
	err := json.Unmarshal(content, &subjects)
	return subjects, err
}

func marshalStore(subjects []Subject) ([]byte, error) {
	return json.Marshal(&subjects)
}
