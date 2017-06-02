// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Passmgr is a simple password manager which allows to securely store
passphrases and retrieve them via commandline.

Retrieved passphrases are copied to the clipboard for a limited amount of
time, in order to be pasted in a passphrase field by the user. After this
time, the clipboard is erased.

All credentials are stored encrypted with AES256-GCM in a file named
'.passmgr_store' which is located in the users home directory by default.
The key to encrypt/decrypt the file is derived from the master passphrase
using scrypt.
*/
package main
