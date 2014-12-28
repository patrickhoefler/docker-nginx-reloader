.PHONY: build dist

default: build

build:
	go build -v -o ./bin/docker-nginx-reloader

dist:
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/linux-amd64/docker-nginx-reloader

test: dist
	docker build -t nginx-debug testing/nginx-debug
	docker run --name nginx-debug -d nginx-debug

	docker build -t docker-nginx-reloader .
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock docker-nginx-reloader
	docker logs nginx-debug

	docker rm -f nginx-debug

update:
	go get -u -f -v
	godep save
