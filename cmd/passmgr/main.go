// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os/user"
	"path/filepath"
)

type command int

const (
	addCmd command = iota
	getCmd
	delCmd
)

func main() {
	cmd, filename := parseCmd()

	// setup:
	app := NewPassmgrTerm(filename)
	if isFile(filename) {
		app.Init()
	} else {
		app.InitEmpty()
	}

	// main loop:
	app.PrintTable()
	success := false
	for {
		switch cmd {
		case getCmd:
			success = app.Get()
		case addCmd:
			success = app.Add()
			if success {
				app.PrintTable()
			}
		case delCmd:
			success = app.Delete()
			if success {
				app.PrintTable()
			}
		default:
			panic("illegal command")
		}

		if success {
			break
		}
	}
}

func parseCmd() (command, string) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".passmgr_store")

	// cmd parsing:
	add := flag.Bool("add", false, "store a new set of credentials")
	del := flag.Bool("del", false, "delete a set of credentials")
	filename := flag.String("file", defaultFilename, "specify the passmgr store file")
	flag.Parse()

	cmd := getCmd
	if *add {
		cmd = addCmd
	} else if *del {
		cmd = delCmd
	}
	return cmd, calcFilename(*filename)
}
