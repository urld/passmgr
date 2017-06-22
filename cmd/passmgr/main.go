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
	filterCmd
	quitCmd
	noCmd
)

func main() {
	cmd, app := parseCmd()

	// setup:
	if isFile(app.filename) {
		app.Init()
	} else {
		app.InitEmpty()
	}
	app.Import()
	loop(app, cmd)
}

func parseCmd() (command, termApp) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".passmgr_store")

	// cmd parsing:
	add := flag.Bool("add", false, "store new credentials")
	del := flag.Bool("del", false, "delete stored credentials")
	filename := flag.String("file", defaultFilename, "specify the passmgr store")
	appTTL := flag.Int("appTTL", 120, "time in seconds after which the application quits if there is no user interaction")
	clipboardTTL := flag.Int("clipboardTTL", 15, "time in seconds after which the clipboard is reset")
	importFilename := flag.String("import", "", "file to import credentials from")
	flag.Parse()

	cmd := noCmd
	if *add {
		cmd = addCmd
	} else if *del {
		cmd = delCmd
	}
	return cmd, termApp{filename: calcFilename(*filename), clipboardTTL: *clipboardTTL, appTTL: *appTTL, importFilename: *importFilename}
}

func loop(app termApp, cmd command) {
	app.PrintTable()
	success := false
	for {
		timer := time.AfterFunc(time.Duration(app.appTTL)*time.Second, func() {
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
		case filterCmd:
			success = app.Filter()
		case quitCmd:
			return
		default:
			cmd = askCommand()
		}
		if success {
			app.PrintTable()
			cmd = askCommand()
		}
		timer.Stop()
	}
}
