.DEFAULT_GOAL := build
fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	GOOS=darwin GOARCH=arm64 go build -o ./release/films_MacOS -ldflags "-w -s" main.go
	GOOS=windows GOARCH=386 go build -o ./release/films_win32 -ldflags "-w -s" main.go
	GOOS=windows GOARCH=amd64 go build -o ./release/films_win64 -ldflags "-w -s" main.go
	GOOS=linux GOARCH=386 go build -o ./release/films_linux32 -ldflags "-w -s" main.go
	GOOS=linux GOARCH=amd64 go build -o ./release/films_linux64 -ldflags "-w -s" main.go
	GOOS=linux GOARCH=arm64 go build -o ./release/films_linux_arm64 -ldflags "-w -s" main.go 
.PHONY:build