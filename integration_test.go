// +build integration

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func runCommand(command string) (string, error) {
	splitCommand := strings.Split(command, " ")
	output, err := exec.Command(splitCommand[0], splitCommand[1:]...).CombinedOutput()
	return string(output), err
}

func TestIntegrationWithDockerClient(t *testing.T) {
	Convey("Given the main binary is built and all Docker images are built and started", t, func() {
		//Build main binary
		os.Setenv("GOOS", "linux")
		os.Setenv("GOARCH", "amd64")
		runCommand("go build -v")

		// Build docker image for main binary
		runCommand("docker build --tag=patrickhoefler/docker-nginx-reloader .")

		// Build nginx docker image with debugging enable
		runCommand("docker build --tag=patrickhoefler/nginx-debug testing/nginx-debug")

		// Remove untagged Docker images
		danglingImages, _ := runCommand("docker images -q --filter=dangling=true")
		for _, danglingImage := range strings.Split(strings.TrimSpace(danglingImages), "\n") {
			runCommand("docker rmi " + danglingImage)
		}

		// Start containers
		runCommand("docker run --name=nginx-debug --detach patrickhoefler/nginx-debug")
		runCommand("docker run --name=gubed-xnign --detach patrickhoefler/nginx-debug")

		Convey("When docker-nginx-reloader is run with no command line flags", func() {
			runCommand("docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader")

			Convey("nginx-debug should be reloaded, but gubed-xnign should not", func() {
				output, _ := runCommand("docker logs nginx-debug")
				So(output, ShouldContainSubstring, "signal 1 (SIGHUP) received")

				output, _ = runCommand("docker logs gubed-xnign")
				So(output, ShouldNotContainSubstring, "signal 1 (SIGHUP) received")
			})
		})

		Convey("When docker-nginx-reloader is run with the command line flag --fragment=xnign", func() {
			runCommand("docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader --fragment=xnign")

			Convey("nginx-debug should not be reloaded, but gubed-xnign should", func() {
				output, _ := runCommand("docker logs nginx-debug")
				So(output, ShouldNotContainSubstring, "signal 1 (SIGHUP) received")

				output, _ = runCommand("docker logs gubed-xnign")
				So(output, ShouldContainSubstring, "signal 1 (SIGHUP) received")
			})
		})

		Reset(func() {
			// This reset is run after each `Convey` at the same scope.
			// Stop and remove containers
			runCommand("docker rm --force nginx-debug")
			runCommand("docker rm --force gubed-xnign")
		})
	})
}
