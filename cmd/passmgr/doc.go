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
  2)   david               twitter.com
  3)   david               example.com
  4)   other@example.com   twitter.com

  Command: (S)elect, (f)ilter, (a)dd, (d)elete or (q)uit?
  Select: 1

  Passphrase copied to clipboard!
  Clipboard will be erased in 15 seconds.

  ...............

  Passphrase erased from clipboard.

Filter Example:
  # ...

  Command: (S)elect, (f)ilter, (a)dd, (d)elete or (q)uit? f
  Filter: twitterdavid

  n)   User    URL
  2)   david   twitter.com
  3)   david   example.com

The filter can be reset by leaving it empty.

Import Example:
  $ passmgr -import dump.json
  [passmgr] master passphrase for /home/david/.passmgr_store:

  n)   User                URL
  1)   urld                github.com
  2)   david               twitter.com
  3)   david               example.com
  4)   other@example.com   twitter.com
  5)   import1             github.com
  6)   import2             google.com
  7)   import3             facebook.com

  Do you wish to save the imported changes? [Y/n]

The dump.json has to look like this:
  [
    {"User":"import1", "URL":"github.com", "Secrets":{"passphrase":"secret2"}},
    {"User":"import2", "URL":"google.com", "Secrets":{"passphrase":"secret3"}},
    {"User":"import3", "URL":"facebook.com", "Secrets":{"passphrase":"secret4"}},
  ]

Please make sure you delete the json file after it is imported, since it contains
all your secrets in plaintext.


*/
package main
