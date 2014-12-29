package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/pkg/mflag"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	const (
		version = "0.1.0-beta2"
	)

	var (
		host         = mflag.String([]string{"H", "-host"}, "unix:///var/run/docker.sock", "The Docker socket to connect to, specified using tcp://host:port or unix:///path/to/socket.")
		fragment     = mflag.String([]string{"-fragment"}, "nginx", "All running Docker containers whose names contains this fragement will receive the SIGHUP signal.")
		printVersion = mflag.Bool([]string{"v", "-version"}, false, "Print the version of docker-nginx-reloader and exit.")
	)

	mflag.Usage = func() {
		message := "usage: docker-nginx-reloader [options]\n\nSends a SIGHUP signal to all running Docker containers whose name contains the given fragment.\n\nOptions:\n"
		fmt.Fprint(os.Stderr, message)
		mflag.PrintDefaults()
	}

	mflag.Parse()

	if *printVersion {
		fmt.Fprintln(os.Stdout, version)
		os.Exit(0)
	}

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
