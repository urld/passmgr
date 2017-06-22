# passmgr

[![Build Status](https://travis-ci.org/urld/passmgr.svg?branch=master)](https://travis-ci.org/urld/passmgr)
[![Go Report Card](https://goreportcard.com/badge/github.com/urld/passmgr)](https://goreportcard.com/report/github.com/urld/passmgr)
[![GoDoc](https://godoc.org/github.com/urld/passmgr/cmd/passmgr?status.svg)](https://godoc.org/github.com/urld/passmgr/cmd/passmgr)

`passmgr` is a simple portable password manager.

## Usage

see [GoDoc](https://godoc.org/github.com/urld/passmgr/cmd/passmgr)


## Install

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
