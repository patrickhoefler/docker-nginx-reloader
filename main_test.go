package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrintProgramVersion(t *testing.T) {
	// Redirect exit()
	exit = func(code int) {}
	// Redirect stdout
	outBuf := &bytes.Buffer{}
	stdout = outBuf

	Convey("When main() is called with the command line flag --version", t, func() {
		versionFlag = true
		main()

		Convey("The program version should be printed to stdout", func() {
			So(outBuf.String(), ShouldEqual, version+"\n")
		})

		// Reset
		versionFlag = false
	})
}

type dockerClient struct {
	reloadedContainers []string
	listErr            error
}

func (m *dockerClient) KillContainer(options docker.KillContainerOptions) error {
	// Remember which containers have been reloaded
	m.reloadedContainers = append(m.reloadedContainers, options.ID)

	// Return no error
	return nil
}

func (m *dockerClient) ListContainers(options docker.ListContainersOptions) ([]docker.APIContainers, error) {
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
	return []docker.APIContainers{container1, container2}, m.listErr
}

func TestNewDockerClientError(t *testing.T) {
	// Redirect exit()
	exit = func(code int) {}
	// Redirect stderr
	errBuf := &bytes.Buffer{}
	stderr = errBuf

	Convey("When main() is called and docker.NewClient() throws an error", t, func() {
		// Mock unsuccessful docker.NewClient()
		newDockerClient = func(host string) (dockerManager, error) {
			return nil, errors.New("No Docker client for you!")
		}
		main()

		Convey("It should be printed to Stderr", func() {
			So(errBuf.String(), ShouldEqual, "No Docker client for you!\n")
		})
	})
}

func TestListContainersError(t *testing.T) {
	// Redirect exit()
	exit = func(code int) {}
	// Redirect stderr
	errBuf := &bytes.Buffer{}
	stderr = errBuf

	Convey("When main() is called and dockerClient.ListContainers() throws an error", t, func() {
		// Create new mock Docker client
		mockDockerClient := &dockerClient{listErr: errors.New("No containers for you!")}
		// Mock docker.NewClient()
		newDockerClient = func(host string) (dockerManager, error) {
			return mockDockerClient, nil
		}
		main()

		Convey("It should be printed to Stderr", func() {
			So(errBuf.String(), ShouldEqual, "No containers for you!\n")
		})
	})
}

func TestReloadContainers(t *testing.T) {
	Convey("Given two Docker containers named nginx-debug and gubed-xnign", t, func() {
		// Create new mock Docker client
		mockDockerClient := &dockerClient{}
		// Mock docker.NewClient()
		newDockerClient = func(host string) (dockerManager, error) {
			return mockDockerClient, nil
		}
		// Redirect stdout
		outBuf := &bytes.Buffer{}
		stdout = outBuf

		Convey("When main() is called with no command line flags", func() {
			main()

			Convey("nginx-debug should be reloaded, but gubed-xnign should not", func() {
				So(mockDockerClient.reloadedContainers, ShouldResemble, []string{"container1"})
				So(outBuf.String(), ShouldEqual, "Sent SIGHUP signal to nginx-debug (container1)\n")
			})
		})

		Convey("When main() is called with the command line flag --fragment=xnign", func() {
			fragment = "xnign"
			main()

			Convey("nginx-debug should not be reloaded, but gubed-xnign should", func() {
				So(mockDockerClient.reloadedContainers, ShouldResemble, []string{"container2"})
				So(outBuf.String(), ShouldEqual, "Sent SIGHUP signal to gubed-xnign (container2)\n")
			})

			// Reset
			fragment = ""
		})
	})
}
