package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/docker/docker/pkg/mflag"
	"github.com/fsouza/go-dockerclient"
)

// Command line flags
var (
	flHost     = mflag.String([]string{"H", "-host"}, "unix:///var/run/docker.sock", "The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.")
	flFragment = mflag.String([]string{"-fragment"}, "nginx", "All running Docker containers whose names contains this fragement will receive the SIGHUP signal.")
	flVersion  = mflag.Bool([]string{"v", "-version"}, false, "Print the version of docker-nginx-reloader and exit.")
)

func init() {
	mflag.Usage = func() {
		message := "usage: docker-nginx-reloader [options]\n\nSends a SIGHUP signal to all running Docker containers whose name contains the given fragment.\n\nOptions:\n"
		fmt.Fprint(os.Stderr, message)
		mflag.PrintDefaults()
	}

	mflag.Parse()

	if *flVersion {
		fmt.Fprintln(os.Stdout, Version)
		os.Exit(0)
	}
}

func main() {
	// Get a Docker client
	client, err := docker.NewClient(*flHost)
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

			matched, err := regexp.MatchString(*flFragment, containerName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if matched {
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
