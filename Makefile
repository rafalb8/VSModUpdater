
all: linux windows

release: all
	cd result && zip VSModUpdater-Windows.zip VSModUpdater.exe
	cd result && tar -czvf VSModUpdater-Linux.tar.gz VSModUpdater

linux:
	CGO_ENABLED=1 \
	GOOS=linux \
	go build -ldflags='-s -w' -o result/VSModUpdater ./cmd/VSModUpdater

windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	go build -ldflags='-s -w' -o result/VSModUpdater.exe ./cmd/VSModUpdater
