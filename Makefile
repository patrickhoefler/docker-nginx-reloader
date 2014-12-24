.PHONY: build dist

default: build

build:
	go build -v -o ./bin/docker-nginx-reloader

dist:
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/linux_amd64/docker-nginx-reloader

update:
	go get -u -f -v
	godep save
