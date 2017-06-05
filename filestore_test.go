// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package passmgr

import (
	"fmt"
	"io"
	"os"
	"testing"
)

const fileUnderTest = "testdata/.store_under_test"
const masterPassphrase = "123"
const passphraseKey = "passphrase"

func getTestFile(name string) string {
	in, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	testFile, err := os.Create(fileUnderTest)
	defer testFile.Close()
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(testFile, in)

	if err != nil {
		panic(err)
	}
	return testFile.Name()
}

func assertEqual(t *testing.T, a, b interface{}, message string) {
	if a == b {
		return
	}
	msg := fmt.Sprintf("%s: %v != %v", message, a, b)
	t.Error(msg)
}

func TestNewStoreIsEmpty(t *testing.T) {
	defer os.Remove(fileUnderTest)
	store, err := ReadFileStore("notExistFile", masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	subjects := store.List()
	if len(subjects) != 0 {
		t.Error("new store is not empty")
	}
}

func TestStoreAddNewSubject(t *testing.T) {
	defer os.Remove(fileUnderTest)
	store, err := ReadFileStore("notExistFile", masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	user := "test"
	url := "example.com"
	secrets := make(map[string]string)
	secrets[passphraseKey] = "secret"

	subj := Subject{User: user, URL: url, Secrets: secrets}
	store.Store(subj)

	subjects := store.List()
	if len(subjects) != 1 {
		t.Error("new subject not added")
	}
	if subjects[0].User != user {
		t.Error("added user != retrieved user")
	}
	if subjects[0].URL != url {
		t.Error("added url != retrieved url")
	}
}

func TestLoadMultipleUsersSingleUrl(t *testing.T) {
	filename := getTestFile("testdata/multipleUsers_singleUrl")
	defer os.Remove(filename)
	store, err := ReadFileStore(filename, masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	subjects := store.List()
	assertEqual(t, 2, len(subjects), "2 subjects should be present")

	subj1, ok := store.Load(Subject{User: "user1", URL: "example.com"})
	if !ok {
		t.Error("successful load should return true")
	}

	assertEqual(t, "user1", subj1.User, "User")
	assertEqual(t, "user1", subjects[0].User, "User")
	assertEqual(t, "example.com", subj1.URL, "Url")
	assertEqual(t, "example.com", subjects[0].URL, "Url")

	assertEqual(t, 0, len(subjects[0].Secrets), "len(Secrets)")
	assertEqual(t, 1, len(subj1.Secrets), "len(Secrets)")
	assertEqual(t, "secret1", subj1.Secrets[passphraseKey], "Secrets")

	subj2, ok := store.Load(Subject{User: "user2", URL: "example.com"})
	if !ok {
		t.Error("successful load should return true")
	}

	assertEqual(t, "user2", subj2.User, "User")
	assertEqual(t, "user2", subjects[1].User, "User")
	assertEqual(t, "example.com", subj2.URL, "Url")
	assertEqual(t, "example.com", subjects[1].URL, "Url")

	assertEqual(t, 0, len(subjects[1].Secrets), "len(Secrets)")
	assertEqual(t, 1, len(subj2.Secrets), "len(Secrets)")
	assertEqual(t, "secret2", subj2.Secrets[passphraseKey], "Secrets")
}

func TestWriteIsReadable(t *testing.T) {
	filename := getTestFile("testdata/multipleUsers_singleUrl")
	defer os.Remove(filename)
	store, err := ReadFileStore(filename, masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	subjects := store.List()
	assertEqual(t, 2, len(subjects), "2 subjects should be present")

	user := "test"
	url := "example.com"
	secrets := make(map[string]string)
	secrets[passphraseKey] = "secret"
	subj := Subject{User: user, URL: url, Secrets: secrets}
	store.Store(subj)

	err = WriteFileStore(store)
	if err != nil {
		t.Fatal(err)
	}

	store2, err := ReadFileStore(filename, masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}
	subjects2 := store2.List()
	assertEqual(t, 3, len(subjects2), "3 subjects should be present")
}

func TestDeleteSubject(t *testing.T) {
	filename := getTestFile("testdata/multipleUsers_singleUrl")
	defer os.Remove(filename)
	store, err := ReadFileStore(filename, masterPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	subjects := store.List()
	assertEqual(t, 2, len(subjects), "2 subjects should be present")
	delSubj, ok := store.Load(Subject{User: "user2", URL: "example.com"})
	if !ok {
		t.Error("successful load should return true")
	}
	if !store.Delete(delSubj) {
		t.Error("successful delete should return true")
	}

	subjects = store.List()
	assertEqual(t, 1, len(subjects), "only 1 subject should be present")
	_, ok = store.Load(Subject{User: "user1", URL: "example.com"})
	if !ok {
		t.Error("successful load should return true")
	}
	_, ok = store.Load(Subject{User: "user2", URL: "example.com"})
	if ok {
		t.Error("failed load should return false")
	}
}

func TestInvalidMagicNumber(t *testing.T) {
	filename := getTestFile("testdata/invalid")
	defer os.Remove(filename)
	_, err := ReadFileStore(filename, masterPassphrase)
	if err == nil {
		t.Error("ReadFileStore should return error on invalid file")
	}
}

func TestInvalidMasterPassphrase(t *testing.T) {
	filename := getTestFile("testdata/invalid")
	defer os.Remove(filename)
	_, err := ReadFileStore(filename, "invalidMasterPassphrase")
	if err == nil {
		t.Error("ReadFileStore should return error on invalid master passphrase")
	}
}
