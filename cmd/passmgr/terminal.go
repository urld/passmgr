// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/urld/passmgr"
	"github.com/urld/passmgr/filestore"
)

// termApp provides means to interact with a passmgr store via terminal.
type termApp struct {
	filename string
	store    passmgr.Store
	subjects []passmgr.Subject
}

func newTermApp(filename string) termApp {
	return termApp{filename: filename}
}

func (app *termApp) Init() {
	if !isFile(app.filename) {
		Fprintln(os.Stderr, "The passmgr store does not exist yet. Add some passphrases first.")
		Fprintln(os.Stderr, "See passmgr -h for help.")
		os.Exit(1)
	}
	masterPassphrase := askSecret("[passmgr] master passphrase for %s: ", app.filename)

	store, err := filestore.Read(app.filename, masterPassphrase)
	if err != nil {
		quitErr(err)
	}
	app.store = store
	app.subjects = store.List()
}

func (app *termApp) InitEmpty() {
	masterPassphrase := askSecret("[passmgr] new master passphrase for %s: ", app.filename)
	if masterPassphrase != askSecret("[passmgr] retype master passphrase for %s: ", app.filename) {
		quitErr(fmt.Errorf("error: passphrases did not match"))
	}
	store, err := filestore.Read(app.filename, masterPassphrase)
	if err != nil {
		quitErr(err)
	}
	app.store = store
	app.subjects = store.List()

	err = filestore.Write(store)
	if err != nil {
		quitErr(err)
	}
}

func (app *termApp) PrintTable() {
	app.subjects = app.store.List()
	if len(app.subjects) == 0 {
		Println("")
		Println("-- store is empty --")
		Println("")
		return
	}

	Println("")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.StripEscape)
	Fprintln(w, "n)\t%s\t%s", "User", "URL")
	for i, c := range app.subjects {
		Fprintln(w, "%d)\t%s\t%s", i+1, c.User, c.URL)
	}
	_ = w.Flush()
	Println("")
}

const passphraseKey = "passphrase"

func (app *termApp) Add() bool {
	var subject passmgr.Subject
	Println("Enter the values for the new entry")
	subject.User = ask("\tUser: ")
	subject.URL = ask("\tURL: ")
	subject.Secrets = make(map[string]string)
	subject.Secrets[passphraseKey] = askSecret("\tPassphrase: ")

	app.store.Store(subject)
	err := filestore.Write(app.store)
	if err != nil {
		quitErr(err)
	}
	return true
}

func (app *termApp) Get() bool {
	if len(app.subjects) == 0 {
		return true
	}

	idx, ok := app.readSelection("Select: ")
	if !ok {
		return false
	}

	subject, _ := app.store.Load(app.subjects[idx-1])
	passphrase, ok := subject.Secrets[passphraseKey]
	if !ok {
		// ignore for now. may become relevant if support for
		// multiple secrets gets added.
		return false
	}

	Println("")
	Println("Passphrase copied to clipboard!")
	Println("Clipboard will be erased in 6 seconds.")
	Println("")
	setClipboard(passphrase)
	for i := 1; i <= 6; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
	resetClipboard()
	Println("")
	Println("")
	Println("Passphrase erased from clipboard.")
	return true
}

func (app *termApp) Delete() bool {
	if len(app.subjects) == 0 {
		return true
	}

	idx, ok := app.readSelection("Delete: ")
	if !ok {
		return false
	}
	subject := app.subjects[idx-1]
	if !askConfirm("Delete all secrets for '%s | %s?", subject.User, subject.URL) {
		return true
	}

	app.store.Delete(app.subjects[idx-1])
	err := filestore.Write(app.store)
	if err != nil {
		quitErr(err)
	}
	return true
}

func (app *termApp) readSelection(prompt string) (int, bool) {
	idx, err := strconv.Atoi(ask(prompt))
	if err != nil {
		Fprintln(os.Stderr, "Please type a valid number.")
		return 0, false
	}
	if idx < 1 || idx > len(app.subjects) {
		Fprintln(os.Stderr, "Please select a number within this range: %d..%d", 1, len(app.subjects))
		return 0, false
	}
	return idx, true
}

func askConfirm(prompt string, a ...interface{}) bool {
	switch strings.ToLower(ask(prompt+" [Y/n] ", a...)) {
	case "y":
		return true
	case "":
		return true
	case "n":
		return false
	default:
		return askConfirm(prompt)
	}
}

func ask(prompt string, a ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt, a...)
	text, err := reader.ReadString('\n')
	if err != nil {
		quitErr(err)
	}
	return strings.TrimRight(text, "\r\n")
}

func askSecret(prompt string, a ...interface{}) string {
	prompt = fmt.Sprintf(prompt, a...)
	secret, err := speakeasy.Ask(prompt)
	if err != nil {
		quitErr(err)
	}
	return secret
}
