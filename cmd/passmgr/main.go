// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/urld/passmgr"
)

type actionType int

const (
	actionGet actionType = iota
	actionDel
)

func main() {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".passmgr_store")

	add := flag.Bool("add", false, "store a new set of credentials")
	del := flag.Bool("del", false, "delete a set of credentials")
	filename := flag.String("file", defaultFilename, "specify the passmgr store file")
	flag.Parse()

	if *add {
		mainAdd(calcFilename(*filename))
	} else if *del {
		mainList(calcFilename(*filename), actionDel)
	} else {
		mainList(calcFilename(*filename), actionGet)
	}
}

const passphraseKey = "passphrase"

func mainAdd(filename string) {
	subject := passmgr.Subject{}
	subject.Description = ask("Description: ")
	subject.User = ask("User: ")
	subject.Secrets = make(map[string]string)
	subject.Secrets[passphraseKey] = askSecret("Passphrase: ")
	fmt.Println("OK")

	fmt.Println("Master Passphrase to modify passmgr storage is required...")
	masterPassphrase := askSecret("Master Passphrase: ")
	fileStore, err := passmgr.ReadFileStore(filename, masterPassphrase)
	if err != nil {
		quitErr(err)
	}

	fileStore.Store(subject)
	err = passmgr.WriteFileStore(fileStore)
	if err != nil {
		quitErr(err)
	}

}

func mainList(filename string, action actionType) {
	if !isFile(filename) {
		fmt.Fprintln(os.Stderr, "The passmgr store does not exist yet. You need to add some passphrases first.\nSee passmgr -h for help.")
		os.Exit(1)
	}
	fmt.Println("Master Passphrase to access passmgr storage is required...")
	masterPassphrase := askSecret("Master Passphrase: ")
	fileStore, err := passmgr.ReadFileStore(filename, masterPassphrase)
	if err != nil {
		quitErr(err)
	}

	// print list of known entries
	subjects := fileStore.List()
	if len(subjects) == 0 {
		fmt.Fprintln(os.Stderr, "The passmgr store is empty. You need to add some passphrases first.\nSee passmgr -h for help.")
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.StripEscape)
	fmt.Fprintf(w, "n)\t%s\t%s\n", "User", "Description")
	for i, c := range subjects {
		fmt.Fprintf(w, "%d)\t%s\t%s\n", i+1, c.User, c.Description)
	}
	_ = w.Flush()

	var actionPrompt string
	switch action {
	case actionGet:
		actionPrompt = "Select: "
	case actionDel:
		actionPrompt = "Delete: "
	default:
		panic("unknown action")
	}

	// try to select entry
	for {
		idx, err := strconv.Atoi(ask(actionPrompt))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Please insert a valid number.")
			continue
		}
		if idx < 1 || idx > len(subjects) {
			fmt.Fprintf(os.Stderr, "Please insert a number within this range: %d-%d\n", 1, len(subjects))
			continue
		}
		switch action {
		case actionGet:
			subject, _ := fileStore.Load(subjects[idx-1])
			passphrase, ok := subject.Secrets[passphraseKey]
			if !ok {
				continue
			}
			fmt.Println("passphrase copied to clipboard...")
			setClipboard(passphrase)
			time.Sleep(5 * time.Second)
			resetClipboard()
			fmt.Println("passphrase erased from clipboard")
		case actionDel:
			fileStore.Delete(subjects[idx-1])
			err = passmgr.WriteFileStore(fileStore)
			if err != nil {
				quitErr(err)
			}
		default:
			panic("unknown action")
		}

		break
	}

}

func isDir(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func calcFilename(filename string) string {
	if isDir(filename) {
		return filepath.Join(filename, ".passmgr_store")
	}
	return filename
}

func quitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
