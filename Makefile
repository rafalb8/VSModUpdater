TAG ?= v0.0.0
LDFLAGS = "-s -w -X 'github.com/rafalb8/VSModUpdater/internal/config.VersionNum=$(TAG)'"

all: linux windows

release: all
	cd result && zip VSModUpdater-$(TAG).zip VSModUpdater VSModUpdater.exe
	cd result && zip VSModUpdater-Windows.zip VSModUpdater.exe
	cd result && tar -czvf VSModUpdater-Linux.tar.gz VSModUpdater

linux:
	CGO_ENABLED=0 \
	GOOS=linux \
	go build -ldflags=$(LDFLAGS) -o result/VSModUpdater ./cmd/VSModUpdater

windows:
	CGO_ENABLED=0 \
	GOOS=windows \
	go build -ldflags=$(LDFLAGS) -o result/VSModUpdater.exe ./cmd/VSModUpdater

