// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/atotto/clipboard"
)

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
