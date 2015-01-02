package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

const (
	version = "0.1.0" // Program version

	exitOK  = 0 // Terminate without error
	exitErr = 1 // Terminate with error
)

var (
	// Flags
	host        string // Docker host
	fragment    string // Search string for reloading the right Docker containers
	versionFlag bool   // If true, only print the program version and exit

	// Global functions for unit testing
	newDockerClient newDockerManager = func(host string) (dockerManager, error) { return docker.NewClient(host) }
	exit            exiter           = func(code int) { os.Exit(code) }
	stdout          io.Writer        = os.Stdout
	stderr          io.Writer        = os.Stderr
)

// Takes care of exiting the program
type exiter func(code int)

// Defines the condition of the program exit
type exitCondition struct {
	code int
	err  error
}

// Function type for newDockerClient
type newDockerManager func(host string) (dockerManager, error)

// Interface for dockerClient
type dockerManager interface {
	ListContainers(docker.ListContainersOptions) ([]docker.APIContainers, error)
	KillContainer(docker.KillContainerOptions) error
}

func init() {
	// Flags
	flag.StringVar(&fragment, "fragment", "nginx", "All running Docker containers whose names contain this fragement will receive the SIGHUP signal.")
	flag.StringVar(&host, "host", "unix:///var/run/docker.sock", "The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.")
	flag.BoolVar(&versionFlag, "version", false, "Print the version of docker-nginx-reloader and exit.")
}

func main() {
	defer handleExit() // Graceful and testable exit

	flag.Parse()

	// Does the user only want the version?
	if versionFlag {
		fmt.Fprintln(stdout, version)
		panic(exitCondition{exitOK, nil})
	}

	// Get a Docker client
	dockerClient, err := newDockerClient(host)
	if err != nil {
		panic(exitCondition{exitErr, err})
	}

	// Get a list of all running Docker containers
	containers, err := dockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		panic(exitCondition{exitErr, err})
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if strings.Index(containerName, fragment) >= 0 {
				dockerClient.KillContainer(
					docker.KillContainerOptions{
						ID:     container.ID,
						Signal: docker.SIGHUP,
					},
				)
				fmt.Fprintf(
					stdout,
					"Sent SIGHUP signal to %s (%s)\n",
					containerName[1:],
					container.ID,
				)
			}
		}
	}
}

func handleExit() {
	// Recover the thrown panic
	if err := recover(); err != nil {
		// Check if the panic was thrown by us
		if p, ok := err.(exitCondition); ok {
			// If there was an error, print it
			if p.code == exitErr {
				fmt.Fprintln(stderr, p.err)
			}
			// Terminate the program with the appropriate exit code
			exit(p.code)
		}
	}
}
