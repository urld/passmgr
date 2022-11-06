#!/bin/sh

set -e 
MAKE="make -e"
export VERSION=$(git describe --exact-match --tags 2>/dev/null || git log -n1 --pretty='%h')


$MAKE test

GOOS=linux GOARCH=amd64 $MAKE shrink dist
GOOS=windows GOARCH=amd64 $MAKE dist
