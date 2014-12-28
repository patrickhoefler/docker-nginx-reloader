package main

import (
	"fmt"
	"os"

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
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, container := range containers {
		fmt.Println("ID: ", container.ID)
		fmt.Println("Image: ", container.Image)
		fmt.Println("Command: ", container.Command)
		fmt.Println("Created: ", container.Created)
		fmt.Println("Status: ", container.Status)
		fmt.Println("Ports: ", container.Ports)
	}

	// Send the SIGHUP signal to all matching containers
	fmt.Println(*flFragment)
}
