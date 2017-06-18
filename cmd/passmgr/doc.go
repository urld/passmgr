// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Passmgr is a simple password manager which allows to securely store
passphrases and retrieve them via commandline.

Usage of passmgr:
  -add
    	store new credentials
  -appTTL int
    	time in seconds after which the application quits if there is no user interaction (default 120)
  -clipboardTTL int
    	time in seconds after which the clipboard is reset (default 15)
  -del
    	delete stored credentials
  -file string
    	specify the passmgr store (default "/home/david/.passmgr_store")


In its default mode (no arguments), passmgr allows to select stored passphrases
which are then copied to the clipboard for a limited amount of time in order
to be pasted into a passphrase field. After this time, the clipboard is erased.

Select Example:
  $ passmgr
  [passmgr] master passphrase for /home/david/.passmgr_store:

  n)   User                URL
  1)   urld                github.com
  2)   david@example.com   facebook.com
  3)   david               twitter.com
  4)   other@example.com   google.com

  Choose a command [(S)elect/(f)ilter/(a)dd/(d)elete/(q)uit]
  Select: 1

  Passphrase copied to clipboard!
  Clipboard will be erased in 15 seconds.

  ...............

  Passphrase erased from clipboard.

Filter Example:
  # ...

  Choose a command [(S)elect/(f)ilter/(a)dd/(d)elete/(q)uit] f
  Filter: david

  n)   User                URL
  2)   david@example.com   facebook.com
  3)   david               twitter.com


All credentials are stored AES256-GCM encrypted in a single file which by default
is located in the users home directory.
The encryption key for this file is derived from a master passphrase using scrypt.

*/
package main
