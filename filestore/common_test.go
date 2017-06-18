// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package filestore

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}, message string) {
	switch a.(type) {
	case []byte:
		if bytes.Equal(a.([]byte), b.([]byte)) {
			return
		}
	default:
		if a == b {
			return
		}
	}
	msg := fmt.Sprintf("%s: %v != %v", message, a, b)
	t.Error(msg)
}

func remove(t *testing.T, filename string) {
	err := os.Remove(filename)
	if err != nil {
		t.Fatal(err)
	}
}

func close(t *testing.T, c io.Closer) {
	err := c.Close()
	if err != nil {
		t.Fatal(err)
	}
}
