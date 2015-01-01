package main

import (
	"bytes"
	"testing"

	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrintProgramVersion(t *testing.T) {
	// Redirect stdout
	outBuf := &bytes.Buffer{}
	stdout = outBuf

	Convey("Given the command line flag --version", t, func() {
		versionFlag = true

		Convey("When main() is called", func() {
			main()

			Convey("The program version should be printed to stdout", func() {
				So(outBuf.String(), ShouldEqual, version+"\n")
			})
		})

		// Reset
		versionFlag = false
	})
}

// Prerequisites for TestReloadContainers
var reloadedContainers []string

type mockDockerClient struct {
}

func (m mockDockerClient) KillContainer(options docker.KillContainerOptions) error {
	// Remember which containers have been reloaded
	reloadedContainers = append(reloadedContainers, options.ID)

	// There was no error. There never will be.
	return nil
}

func (m mockDockerClient) ListContainers(options docker.ListContainersOptions) ([]docker.APIContainers, error) {
	// Create mock containers
	container1 := docker.APIContainers{
		ID:    "container1",
		Names: []string{"/nginx-debug"},
	}
	container2 := docker.APIContainers{
		ID:    "container2",
		Names: []string{"/gubed-xnign"},
	}

	// Return container slice and no error
	return []docker.APIContainers{container1, container2}, nil
}

func TestReloadContainers(t *testing.T) {
	Convey("Given two Docker containers named nginx-debug and gubed-xnign", t, func() {
		// Redirect stdout
		outBuf := &bytes.Buffer{}
		stdout = outBuf

		Convey("When reloadContainers() is called and no command line flags were set", func() {
			reloadContainers(mockDockerClient{})

			Convey("nginx-debug should be reloaded, but gubed-xnign should not", func() {
				So(reloadedContainers, ShouldResemble, []string{"container1"})
				So(outBuf.String(), ShouldEqual, "Sent SIGHUP signal to nginx-debug (container1)\n")
			})
		})

		Convey("When reloadContainers() is called and the command line flag --fragment=xnign was set", func() {
			fragment = "xnign"
			reloadContainers(mockDockerClient{})

			Convey("nginx-debug should not be reloaded, but gubed-xnign should", func() {
				So(reloadedContainers, ShouldResemble, []string{"container2"})
				So(outBuf.String(), ShouldEqual, "Sent SIGHUP signal to gubed-xnign (container2)\n")
			})

			// Reset
			fragment = ""
		})

		Reset(func() {
			// This reset is run after each `Convey` at the same scope.
			reloadedContainers = []string{}
		})
	})
}
