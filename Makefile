
all: linux windows

linux:
	CGO_ENABLED=1 \
	GOOS=linux \
	go build -v -o result/VSModUpdater ./cmd/VSModUpdater

windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	go build -v -o result/VSModUpdater.exe ./cmd/VSModUpdater