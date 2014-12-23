.PHONY: build dist

default: build

build:
	godep go build -v -o ./bin/docker-nginx-reloader ./src

dist:
	GOOS=linux GOARCH=amd64 godep go build -v -o ./bin/linux_amd64/docker-nginx-reloader ./src

update:
	go get -u -f ./src
	godep save ./src
