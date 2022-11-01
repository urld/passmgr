PACKAGE      = passmgr
VERSION      = $(shell git log -n1 --pretty='%h')
BUILD_DIR    = build
RELEASE_DIR  = dist
RELEASE_FILE = $(PACKAGE)_$(VERSION)_$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all clean clean_build clean_dist dist build install test


all: test install dist



dist: build shrink
	mkdir -p $(RELEASE_DIR)
	go-licenses save "github.com/urld/passmgr/cmd/passmgr" --save_path="$(BUILD_DIR)/licenses"
	cp LICENSE $(BUILD_DIR)/licenses/passmgr.LICENSE
	tar -cvzf  $(RELEASE_DIR)/$(RELEASE_FILE).tar.gz $(BUILD_DIR) --transform='s/$(BUILD_DIR)/$(RELEASE_FILE)/g'

shrink: build
	strip $(BUILD_DIR)/passmgr*
	upx $(BUILD_DIR)/passmgr*

build: clean_build
	mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && \
	go build github.com/urld/passmgr/cmd/passmgr

test:
	go test github.com/urld/passmgr/...


install:
	go install github.com/urld/passmgr/cmd/passmgr


clean: clean_build clean_dist


clean_build:
	rm -rf $(BUILD_DIR)


clean_dist:
	rm -rf $(RELEASE_DIR)
