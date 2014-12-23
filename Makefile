.PHONY: build dist

default: build

build:
	godep go build -v -a -o ./bin/docker-nginx-reloader ./src

dist:
	GOOS=linux GOARCH=amd64 godep go build -v -a -o ./bin/linux_amd64/docker-nginx-reloader ./src
