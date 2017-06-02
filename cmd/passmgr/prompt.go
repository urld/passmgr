// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bgentry/speakeasy"
)

func ask(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		quitErr(err)
	}
	return text[:len(text)-1]
}

func askSecret(prompt string) string {
	secret, err := speakeasy.Ask(prompt)
	if err != nil {
		quitErr(err)
	}
	return secret
}
