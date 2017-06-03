// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passmgr

import (
	"encoding/json"
	"errors"
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
		l[i] = Subject{Description: c.Description, User: c.User}
	}
	return l
}

func (s *fileStore) Load(subject Subject) (Subject, bool) {
	for _, sc := range s.subjects {
		if sc.Description == subject.Description && sc.User == subject.User {
			return sc, true
		}
	}
	return subject, false
}

func (s *fileStore) Store(newSubject Subject) {
	for i, c := range s.subjects {
		if c.Description == newSubject.Description && c.User == newSubject.User {
			s.subjects[i] = newSubject
			return
		}
	}
	s.subjects = append(s.subjects, newSubject)
}

func (s *fileStore) Delete(subject Subject) bool {
	for i, c := range s.subjects {
		if c.Description == subject.Description && c.User == subject.User {
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
	return ioutil.WriteFile(s.filename, content, os.FileMode(0600))
}

// ReadFileStore reads and decrypts the encrypted contents of a Store from the
// specified file.
// The file consists of a fixed size salt, followed by a ciphertext of
// arbitrary length. The salt is required to derive a encryption key from
// the master passphrase.
func ReadFileStore(filename, passphrase string) (Store, error) {
	content, err := readSecretFile(filename)
	if err != nil {
		return nil, err
	}

	salt := content[:saltLen]
	key, err := deriveKey([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}

	if len(content) == saltLen {
		return &fileStore{filename: filename, salt: salt, key: key}, nil
	}

	ciphertext := content[saltLen:]
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
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return genSalt()
	}
	if err != nil {
		return nil, err
	}
	if info.Mode().Perm() != 0600 {
		return nil, errors.New("passmgr store file must have permissions set to 0600")
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
