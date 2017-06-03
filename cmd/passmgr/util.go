// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)

func quitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func calcFilename(filename string) string {
	if isDir(filename) {
		return filepath.Join(filename, ".passmgr_store")
	}
	return filename
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

func setClipboard(str string) {
	err := clipboard.WriteAll(str)
	if err != nil {
		quitErr(err)
	}
}

func resetClipboard() {
	err := clipboard.WriteAll("")
	if err != nil {
		quitErr(err)
	}
}
