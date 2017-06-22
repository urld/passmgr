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
  -change-key
    	change the master passphrase
  -clipboardTTL int
    	time in seconds after which the clipboard is reset (default 15)
  -del
    	delete stored credentials
  -file string
    	specify the passmgr store (default "/home/david/.passmgr_store")
  -import string
    	file to import credentials from


In its default mode (no arguments), passmgr allows to select stored passphrases
which are then copied to the clipboard for a limited amount of time in order
to be pasted into a passphrase field. After this time, the clipboard is erased.

All credentials are stored AES256-GCM encrypted in a single file which by default
is located in the users home directory.
The encryption key for this file is derived from a master passphrase using scrypt.

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

Import Example:
  $ passmgr -import dump.json
  [passmgr] master passphrase for /home/david/.passmgr_store:

  n)   User    URL
  1)   david   example.com
  2)   david   github.com
  3)   david   google.com
  4)   david   facebook.com

  Do you which to save the imported changes? [Y/n]

Where dump.json looks like this:
  [
    {"User":"david", "URL":"github.com", "Secrets":{"passphrase":"secret2"}},
    {"User":"david", "URL":"google.com", "Secrets":{"passphrase":"secret3"}},
    {"User":"david", "URL":"facebook.com", "Secrets":{"passphrase":"secret4"}},
  ]

Please make sure you delete the json file after it is imported, since it contains
all your secrets in plaintext.



*/
package main
