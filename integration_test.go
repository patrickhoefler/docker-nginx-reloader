// +build integration

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Manually parse the command line flags to get
	// access to testing.Verbose() outside of m.Run()
	flag.Parse()

	// Test setup
	setup()

	// Run the tests and exit
	os.Exit(m.Run())
}

func setup() {
	if testing.Verbose() {
		log.Println("Building docker-nginx-reloader")
	}
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	runCommand("go build -v")

	if testing.Verbose() {
		log.Println("Building patrickhoefler/nginx-debug image")
	}
	runCommand("docker build --tag=patrickhoefler/nginx-debug testing/nginx-debug")

	if testing.Verbose() {
		log.Println("Building patrickhoefler/docker-nginx-reloader image")
	}
	runCommand("docker build --tag=patrickhoefler/docker-nginx-reloader .")

	if testing.Verbose() {
		log.Println("Removing untagged Docker images")
	}
	danglingImages := runCommand("docker images -q --filter=dangling=true")
	for _, danglingImage := range strings.Split(danglingImages, "\n") {
		runCommandThatMayFail("docker rmi " + danglingImage)
	}
}

func runCommand(command string) string {
	splitCommand := strings.Split(command, " ")
	output, err := exec.Command(splitCommand[0], splitCommand[1:]...).CombinedOutput()
	if err != nil {
		log.Println(string(output))
		log.Fatal(err)
	}
	return string(output)
}

func runCommandThatMayFail(command string) string {
	splitCommand := strings.Split(command, " ")
	output, err := exec.Command(splitCommand[0], splitCommand[1:]...).CombinedOutput()
	if testing.Verbose() && err != nil {
		log.Println(string(output))
		log.Println(err)
	}
	return string(output)
}

func runTestCommand(t *testing.T, command string) string {
	splitCommand := strings.Split(command, " ")
	output, err := exec.Command(splitCommand[0], splitCommand[1:]...).CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Error(err)
	}
	return string(output)
}

func runIntegrationTest(t *testing.T, dockerCommand string) {
	// Make sure that a potentially leftover nginx-debug container is removed
	runCommandThatMayFail("docker rm --force nginx-debug")

	// Start nginx-debug container
	runTestCommand(t, "docker run --name=nginx-debug --detach patrickhoefler/nginx-debug")

	// Always remove nginx-debug container
	defer runTestCommand(t, "docker rm --force nginx-debug")

	// Run docker-nginx-reloader
	runTestCommand(t, dockerCommand)

	output := runTestCommand(t, "docker logs nginx-debug")
	if strings.Index(output, "signal 1 (SIGHUP) received") < 0 {
		t.Error("nginx didn't receive the SIGHUP signal")
	}
}

func TestWithoutFlags(t *testing.T) {
	// Run docker-nginx-reloader without flags
	runIntegrationTest(t, "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader")
}

func TestWithCustomFragment(t *testing.T) {
	// Run docker-nginx-reloader with fragment flag
	runIntegrationTest(t, "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock patrickhoefler/docker-nginx-reloader /docker-nginx-reloader --fragment=debug")
}
