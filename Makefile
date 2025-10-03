TAG ?= v0.0.0
LDFLAGS = "-s -w -X 'github.com/rafalb8/VSModUpdater/internal/config.VersionNum=$(TAG)'"

.PHONY: all linux windows darwin darwin-amd64 darwin-arm64 release release-combo release-linux release-windows release-darwin

all: linux windows darwin

release: release-linux release-windows release-darwin release-combo

release-combo: all
	cd result/linux && zip ../VSModUpdater-$(TAG).zip VSModUpdater
	cd result/darwin && zip ../VSModUpdater-$(TAG).zip VSModUpdater_macOS
	cd result/windows && zip ../VSModUpdater-$(TAG).zip VSModUpdater.exe

release-linux: linux
	cd result/linux && tar -czvf ../VSModUpdater-Linux.tar.gz VSModUpdater

release-windows: windows
	cd result/windows && zip ../VSModUpdater-Windows.zip VSModUpdater.exe

release-darwin: darwin
	cd result/darwin && tar -czvf ../VSModUpdater-macOS.tar.gz VSModUpdater_macOS

linux:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -ldflags=$(LDFLAGS) -o result/linux/VSModUpdater ./cmd/VSModUpdater

windows:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=windows \
	go build -ldflags=$(LDFLAGS) -o result/windows/VSModUpdater.exe ./cmd/VSModUpdater

darwin: darwin-amd64 darwin-arm64
	cd result/darwin && go tool lipo -output VSModUpdater_macOS -create VSModUpdater_amd64 VSModUpdater_arm64

darwin-amd64:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=darwin \
	go build -ldflags=$(LDFLAGS) -o result/darwin/VSModUpdater_amd64 ./cmd/VSModUpdater

darwin-arm64:
	CGO_ENABLED=0 \
	GOARCH=arm64 \
	GOOS=darwin \
	go build -ldflags=$(LDFLAGS) -o result/darwin/VSModUpdater_arm64 ./cmd/VSModUpdater

clean:
	rm -rf result/*