PACKAGE      = passmgr
VERSION      = $(shell git log -n1 --pretty='%h')
BUILD_DIR    = build
RELEASE_DIR  = dist
RELEASE_FILE = $(PACKAGE)_$(VERSION)_$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all clean clean_build clean_dist dist build install test


all: test install dist



dist: build
	mkdir -p $(RELEASE_DIR)
	mkdir -p $(BUILD_DIR)/licenses
	cp $(GOPATH)/src/github.com/bgentry/speakeasy/LICENSE $(BUILD_DIR)/licenses/bgentry.speakeasy.LICENSE
	cp $(GOPATH)/src/github.com/bgentry/speakeasy/LICENSE_WINDOWS $(BUILD_DIR)/licenses/bgentry.speakeasy.LICENSE_WINDOWS
	cp $(GOPATH)/src/github.com/atotto/clipboard/LICENSE $(BUILD_DIR)/licenses/atotto.clipboard.LICENSE
	cp $(GOPATH)/src/golang.org/x/crypto/LICENSE $(BUILD_DIR)/licenses/golang.crypto.LICENSE
	cp LICENSE $(BUILD_DIR)/licenses/passmgr.LICENSE
	tar -cvzf  $(RELEASE_DIR)/$(RELEASE_FILE).tar.gz $(BUILD_DIR) --transform='s/$(BUILD_DIR)/$(RELEASE_FILE)/g'

build: clean_build
	mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && \
	go build github.com/urld/passmgr/cmd/passmgr

test:
	go test github.com/urld/passmgr/...


install: test
	go install github.com/urld/passmgr/cmd/passmgr


clean: clean_build clean_dist


clean_build:
	rm -rf $(BUILD_DIR)


clean_dist:
	rm -rf $(RELEASE_DIR)

