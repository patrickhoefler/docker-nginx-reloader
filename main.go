package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	const (
		version = "0.1.0-beta2"
	)

var (
	// Flags
	host        string
	fragment    string
	versionFlag bool


func init() {
	// Flags
	flag.StringVar(&fragment, "fragment", "nginx", "All running Docker containers whose names contain this fragement will receive the SIGHUP signal.")
	flag.StringVar(&host, "host", "unix:///var/run/docker.sock", "The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.")
	flag.BoolVar(&versionFlag, "version", false, "Print the version of docker-nginx-reloader and exit.")
}

	if *printVersion {
		fmt.Fprintln(os.Stdout, version)
		os.Exit(0)
	}
	flag.Parse()

	// Get a Docker client
	client, err := docker.NewClient(*host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Get a list of all running Docker containers
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, container := range containers {

		for _, containerName := range container.Names {

			if strings.Index(containerName, *fragment) >= 0 {
				client.KillContainer(
					docker.KillContainerOptions{
						ID:     container.ID,
						Signal: docker.SIGHUP,
					},
				)
				fmt.Printf(
					"Sent SIGHUP signal to %s (%s)\n",
					containerName[1:],
					container.ID,
				)
			}
		}
	}
}
