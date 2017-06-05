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

	cmd := getCmd
	if *add {
		cmd = addCmd
	} else if *del {
		cmd = delCmd
	}
	return cmd, calcFilename(*filename)
}

func loop(app termApp, cmd command) {
	app.PrintTable()
	success := false
	for {
		timer := time.AfterFunc(1*time.Minute, func() {
			fmt.Println("\n[passmgr] Exited due to inactivity.")
			os.Exit(1)
		})
		switch cmd {
		case getCmd:
			success = app.Get()
			if success {
				if askConfirm("Quit?") {
					return
				}
				app.PrintTable()
			}
		case addCmd:
			success = app.Add()
			if success {
				app.PrintTable()
				if askConfirm("Quit?") {
					return
				}
			}
		case delCmd:
			success = app.Delete()
			if success {
				app.PrintTable()
				if askConfirm("Quit?") {
					return
				}
			}
		default:
			panic("illegal command")
		}
		timer.Stop()
	}
}
