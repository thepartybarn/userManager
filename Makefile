# Go parameters

BuildDate := $(shell date -u +%Y%m%d.%H%M%S)
GitVersion := $(shell git describe --always --long --dirty=-test)
MakeFilePath := $(abspath $(lastword $(MAKEFILE_LIST)))
FolderName := $(notdir $(patsubst %/,%,$(dir $(MakeFilePath))))


all: format test checkin build push

format:
	go get ./...
	go fmt *.go

test: format
	go test -v ./...

checkin: format test
	-rm $(FolderName)
	-git pull
	-git commit -a
	-git push

build: format
	go build -ldflags "-X main._buildDate=$(BuildDate) -X main._gitVersion=$(GitVersion)" -o $(FolderName) *.go

push: format test checkin build
	sudo docker build -t thepartybarn/production:$(FolderName) .
	sudo docker login
	sudo docker push thepartybarn/production:$(FolderName)

run: format build
	sudo ./$(FolderName)
