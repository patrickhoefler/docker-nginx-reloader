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
	version = "0.1.0-beta3"
)

var (
	// Flags
	host        string
	fragment    string
	versionFlag bool

	// Used for testing
	stdout io.Writer = os.Stdout
)

func init() {
	// Flags
	flag.StringVar(&fragment, "fragment", "nginx", "All running Docker containers whose names contain this fragement will receive the SIGHUP signal.")
	flag.StringVar(&host, "host", "unix:///var/run/docker.sock", "The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.")
	flag.BoolVar(&versionFlag, "version", false, "Print the version of docker-nginx-reloader and exit.")
}

func main() {
	flag.Parse()

	// Does the user only want the version?
	if versionFlag {
		fmt.Fprintln(stdout, version)
	} else {

		// Get a Docker client
		client, err := docker.NewClient(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Reload the matching containers
		reloadContainers(client)
	}
}

type dockerManager interface {
	ListContainers(docker.ListContainersOptions) ([]docker.APIContainers, error)
	KillContainer(docker.KillContainerOptions) error
}

func reloadContainers(client dockerManager) {
	// Get a list of all running Docker containers
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, container := range containers {

		for _, containerName := range container.Names {

			if strings.Index(containerName, fragment) >= 0 {
				client.KillContainer(
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
