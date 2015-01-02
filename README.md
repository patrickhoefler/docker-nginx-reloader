docker-nginx-reloader
=====================

[![Build Status](https://travis-ci.org/patrickhoefler/docker-nginx-reloader.svg?branch=master)](https://travis-ci.org/patrickhoefler/docker-nginx-reloader) [![Coverage Status](https://img.shields.io/coveralls/patrickhoefler/docker-nginx-reloader.svg)](https://coveralls.io/r/patrickhoefler/docker-nginx-reloader?branch=master)

Send a SIGHUP signal to Docker containers without knowing their exact name. Especially useful in combination with [Kubernetes](https://github.com/googlecloudplatform/kubernetes).

## Example

Let's say that in your local Docker instance you have two running nginx containers called `dev-nginx` and `test-nginx`. If you run `docker-nginx-reloader` without any command line flags, both nginx containers will receive a SIGHUP signal and the contained nginx instances will reload. If you run `docker-nginx-reloader --fragment=dev`, only the `dev-nginx` instance will be reloaded.

## Usage

```
$ ./docker-nginx-reloader --help
Usage of ./docker-nginx-reloader:
-fragment="nginx": All running Docker containers whose names contain this fragement will receive the SIGHUP signal.
-host="unix:///var/run/docker.sock": The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.
-version=false: Print the version of docker-nginx-reloader and exit.
```

## Usage in a Docker container

You can run the latest dockerized version of docker-nginx-reloader via the following command:

`docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader`

The dockerized version takes the same command line flags, so to stay with the example above you could run:

`docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader --fragment=dev`

You might need to change the `-v /var/run/docker.sock:/var/run/docker.sock` part according to your local setup.

## Building

- `go get -d github.com/patrickhoefler/docker-nginx-reloader`
- `go get github.com/tools/godep`
- `cd $GOPATH/src/github.com/patrickhoefler/docker-nginx-reloader`
- `godep go build` creates the `docker-nginx-reloader` binary in the current directory.
- `godep go test` runs the unit tests.
- `godep go test --tags=integration` runs the unit tests as well as the Docker integration tests and builds a Docker image called `patrickhoefler/docker-nginx-reloader` in your locally configured Docker environment.

## License

[MIT](https://github.com/patrickhoefler/docker-nginx-reloader/blob/master/LICENSE)
