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
		fmt.Fprintln(os.Stderr, "The passmgr store does not exist yet. Add some passphrases first.")
		fmt.Fprintln(os.Stderr, "See passmgr -h for help.")
		os.Exit(1)
	}
	masterPassphrase := askSecret("[passmgr] master passphrase for %s: ", app.filename)

	store, err := passmgr.ReadFileStore(app.filename, masterPassphrase)
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
	store, err := passmgr.ReadFileStore(app.filename, masterPassphrase)
	if err != nil {
		quitErr(err)
	}
	app.store = store
	app.subjects = store.List()

	err = passmgr.WriteFileStore(store)
	if err != nil {
		quitErr(err)
	}
}

func (app *termApp) PrintTable() {
	app.subjects = app.store.List()
	if len(app.subjects) == 0 {
		fmt.Println("\n-- store is empty --\n")
		return
	}

	fmt.Println("")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.StripEscape)
	fmt.Fprintf(w, "n)\t%s\t%s\n", "User", "Description")
	for i, c := range app.subjects {
		fmt.Fprintf(w, "%d)\t%s\t%s\n", i+1, c.User, c.Description)
	}
	_ = w.Flush()
	fmt.Println("")
}

const passphraseKey = "passphrase"

func (app *termApp) Add() bool {
	var subject passmgr.Subject
	subject.Description = ask("Description: ")
	subject.User = ask("User: ")
	subject.Secrets = make(map[string]string)
	subject.Secrets[passphraseKey] = askSecret("Passphrase: ")

	app.store.Store(subject)
	err := passmgr.WriteFileStore(app.store)
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

	fmt.Println("\nPassphrase copied to clipboard!\nClipboard will be erased in 6 seconds.\n")
	setClipboard(passphrase)
	for i := 1; i <= 6; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
	resetClipboard()
	fmt.Println("\n\nPassphrase erased from clipboard.")
	return true
}

func (app *termApp) Delete() bool {
	if len(app.subjects) == 0 {
		return true
	}

	idx, ok := app.readSelection("Select: ")
	if !ok {
		return false
	}
	fmt.Printf("All secrets of '%s | %s' will be deleted.\n", app.subjects[idx-1].User, app.subjects[idx-1].Description)
	if !askConfirm("Do you want to continue ?") {
		return true
	}

	app.store.Delete(app.subjects[idx-1])
	err := passmgr.WriteFileStore(app.store)
	if err != nil {
		quitErr(err)
	}
	return true
}

func (app *termApp) readSelection(prompt string) (int, bool) {
	idx, err := strconv.Atoi(ask("Select: "))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Please type a valid number.")
		return 0, false
	}
	if idx < 1 || idx > len(app.subjects) {
		fmt.Fprintf(os.Stderr, "Please select a number within this range: %d..%d\n", 1, len(app.subjects))
		return 0, false
	}
	return idx, true
}

func askConfirm(prompt string) bool {
	switch strings.ToLower(ask(prompt + " [Y/n] ")) {
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
	return text[:len(text)-1]
}

func askSecret(prompt string, a ...interface{}) string {
	prompt = fmt.Sprintf(prompt, a...)
	secret, err := speakeasy.Ask(prompt)
	if err != nil {
		quitErr(err)
	}
	return secret
}
