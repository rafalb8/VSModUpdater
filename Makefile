TAG ?= v0.0.0
LDFLAGS = "-s -w -X github.com/rafalb8/VSModUpdater/internal/config.version=$(TAG)"

.PHONY: all linux windows darwin darwin-amd64 darwin-arm64 build build-combo build-linux build-windows build-darwin

all: linux windows darwin

build: build-linux build-windows build-darwin build-combo

build-combo: all
	cd bin/linux && zip ../VSModUpdater-$(TAG).zip VSModUpdater
	cd bin/darwin && zip ../VSModUpdater-$(TAG).zip VSModUpdater_macOS
	cd bin/windows && zip ../VSModUpdater-$(TAG).zip VSModUpdater.exe

build-linux: linux
	cd bin/linux && tar -czvf ../VSModUpdater-Linux.tar.gz VSModUpdater

build-windows: windows
	cd bin/windows && zip ../VSModUpdater-Windows.zip VSModUpdater.exe

build-darwin: darwin
	cd bin/darwin && tar -czvf ../VSModUpdater-macOS.tar.gz VSModUpdater_macOS

linux:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -trimpath -ldflags=$(LDFLAGS) -o bin/linux/VSModUpdater .

windows:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=windows \
	go build -trimpath -ldflags=$(LDFLAGS) -o bin/windows/VSModUpdater.exe .

darwin: darwin-amd64 darwin-arm64
	cd bin/darwin && go tool lipo -output VSModUpdater_macOS -create VSModUpdater_amd64 VSModUpdater_arm64

darwin-amd64:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=darwin \
	go build -trimpath -ldflags=$(LDFLAGS) -o bin/darwin/VSModUpdater_amd64 .

darwin-arm64:
	CGO_ENABLED=0 \
	GOARCH=arm64 \
	GOOS=darwin \
	go build -trimpath -ldflags=$(LDFLAGS) -o bin/darwin/VSModUpdater_arm64 .

clean:
	rm -rf bin/*
