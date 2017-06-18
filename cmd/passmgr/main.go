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
	"time"
)

type command int

const (
	addCmd command = iota
	getCmd
	delCmd
	quitCmd
	noCmd
)

func main() {
	cmd, filename := parseCmd()

	// setup:
	app := newTermApp(filename)
	if isFile(filename) {
		app.Init()
	} else {
		app.InitEmpty()
	}
	loop(app, cmd)
}

func parseCmd() (command, string) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".passmgr_store")

	// cmd parsing:
	add := flag.Bool("add", false, "store new credentials")
	del := flag.Bool("del", false, "delete stored credentials")
	filename := flag.String("file", defaultFilename, "specify the passmgr store")
	flag.Parse()

	cmd := noCmd
	if *add {
		cmd = addCmd
	} else if *del {
		cmd = delCmd
	}
	return cmd, calcFilename(*filename)
}

func loop(app termApp, cmd command) {
	app.PrintTable()
	success := true
	for {
		timer := time.AfterFunc(1*time.Minute, func() {
			fmt.Println("\n[passmgr] Exited due to inactivity.")
			os.Exit(1)
		})
		switch cmd {
		case getCmd:
			success = app.Get()
		case addCmd:
			success = app.Add()
		case delCmd:
			success = app.Delete()
		case quitCmd:
			return
		}
		if success {
			app.PrintTable()
			cmd = askCommand()
		}
		timer.Stop()
	}
}
