# passmgr

[![Build Status](https://travis-ci.org/urld/passmgr.svg?branch=master)](https://travis-ci.org/urld/passmgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/urld/passmgr)](https://goreportcard.com/report/github.com/urld/passmgr)
[![GoDoc](https://godoc.org/github.com/urld/passmgr/cmd/passmgr?status.svg)](https://godoc.org/github.com/urld/passmgr/cmd/passmgr)

`passmgr` is a simple portable password manager.

## Usage

Just call ```passmgr``` from your command line.
The application will tell you how to proceed.

Read the [command documentation](https://godoc.org/github.com/urld/passmgr/cmd/passmgr)
for more detailed instructions and examples.


## Install

Download the [latest release](https://github.com/urld/passmgr/releases/latest) for your platform,
or build from source:
```go get github.com/urld/passmgr/cmd/passmgr```


### Dependencies

* [github.com/atotto/clipboard](https://github.com/atotto/clipboard)
* [github.com/bgentry/speakeasy](https://github.com/bgentry/speakeasy)
* [golang.org/x/crypto/scrypt](https://godoc.org/golang.org/x/crypto/scrypt)


### Platforms

* OSX
* Windows
* Linux (requires `xclip` or `xsel` command to be installed, probably works on other *nix platforms too)


## Limitations

* no protection against keyloggers
* no protection against cliboard spies


## TODO

* look into improved clipboard handling [hn comment](https://news.ycombinator.com/item?id=14581411)
