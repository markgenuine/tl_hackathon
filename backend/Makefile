LOCAL_BIN:=$(CURDIR)/bin

build-linux:
	GOOS=linux GOARCH=amd64 go build -o service_linux ./main.go

build:
	go build -o service ./main.go

run:
	go run *.go