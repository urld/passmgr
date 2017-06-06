// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package main

import (
	"fmt"
	"io"
)

func Println(a ...interface{}) {
	fmt.Print(a...)
	fmt.Print("\r\n")
}

func Fprintln(w io.Writer, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	fmt.Fprint(w, s)
	fmt.Fprint(w, "\r\n")
}
