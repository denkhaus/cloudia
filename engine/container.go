package engine

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Container struct {
	id            string
	RawName       string `json:"name" yaml:"name"`
	RawDockerfile string `json:"dockerfile" yaml:"dockerfile"`
	RawImage      string `json:"image" yaml:"image"`
	Run           RunParameters
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Name() string {
	return os.ExpandEnv(container.RawName)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Dockerfile() string {
	return os.ExpandEnv(container.RawDockerfile)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Image() string {
	return os.ExpandEnv(container.RawImage)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Id() (id string, err error) {
	if len(container.id) > 0 {
		id = container.id
	} else {
		// Inspect container, extracting the ID.
		// This will return gibberish if no container is found.
		args := []string{"inspect", "--format={{.ID}}", container.Name()}
		output, outErr := commandOutput("docker", args)
		if err == nil {
			id = output
			container.id = output
		} else {
			err = outErr
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) exists() bool {
	// `ps -a` returns all existant containers
	id, err := container.Id()
	if err != nil || len(id) == 0 {
		return false
	}
	dockerCmd := []string{"docker", "ps", "--quiet", "--all", "--no-trunc"}
	grepCmd := []string{"grep", "-wF", id}
	output, err := pipedCommandOutput(dockerCmd, grepCmd)
	if err != nil {
		return false
	}
	result := string(output)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) running() bool {
	// `ps` returns all running containers
	id, err := container.Id()
	if err != nil || len(id) == 0 {
		return false
	}
	dockerCmd := []string{"docker", "ps", "--quiet", "--no-trunc"}
	grepCmd := []string{"grep", "-wF", id}
	output, err := pipedCommandOutput(dockerCmd, grepCmd)
	if err != nil {
		return false
	}
	result := string(output)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) imageExists() bool {
	dockerCmd := []string{"docker", "images", "--no-trunc"}
	grepCmd := []string{"grep", "-wF", container.Image()}
	output, err := pipedCommandOutput(dockerCmd, grepCmd)
	if err != nil {
		return false
	}
	result := string(output)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) status(w *tabwriter.Writer) {
	args := []string{"inspect", "--format={{.State.Running}}\t{{.ID}}\t{{if .NetworkSettings.IPAddress}}{{.NetworkSettings.IPAddress}}{{else}}-{{end}}\t{{range $k,$v := $.NetworkSettings.Ports}}{{$k}},{{end}}", container.Name()}
	output, err := commandOutput("docker", args)
	if err != nil {
		fmt.Fprintf(w, "%s\tError:%v\t%v\n", container.Name(), err, output)
		return
	}
	fmt.Fprintf(w, "%s\t%s\n", container.Name(), output)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Pull image for container
func (container *Container) pullImage() {
	fmt.Printf("Pulling image %s ... ", container.Image())
	args := []string{"pull", container.Image()}
	executeCommand("docker", args)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Build image for container
func (container *Container) buildImage() {
	fmt.Printf("Building image %s ... ", container.Image())
	args := []string{"build", "--rm", "--tag=" + container.Image(), container.Dockerfile()}
	executeCommand("docker", args)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (container Container) provision(force bool) {
	if force || !container.imageExists() {
		if len(container.Dockerfile()) > 0 {
			container.buildImage()
		} else {
			container.pullImage()
		}
	} else {
		print.Notice("Image %s does already exist. Use --force to recreate.\n", container.Image())
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Run or start container
func (container Container) runOrStart() {
	if container.exists() {
		container.start()
	} else {
		container.run()
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Run container
func (container Container) run() {
	if container.exists() {
		print.Notice("Container %s does already exist. Use --force to recreate.\n", container.Name())
		if !container.running() {
			container.start()
		}
	} else {
		fmt.Printf("Running container %s ... ", container.Name())
		// Assemble command arguments
		args := []string{"run"}
		// Cidfile
		if len(container.Run.Cidfile()) > 0 {
			args = append(args, "--cidfile", container.Run.Cidfile())
		}
		// CPU shares
		if container.Run.CpuShares > 0 {
			args = append(args, "--cpu-shares", strconv.Itoa(container.Run.CpuShares))
		}
		// Detach
		if container.Run.Detach {
			args = append(args, "--detach")
		}
		// Dns
		for _, dns := range container.Run.Dns() {
			args = append(args, "--dns", dns)
		}
		// Entrypoint
		if len(container.Run.Entrypoint()) > 0 {
			args = append(args, "--entrypoint", container.Run.Entrypoint())
		}
		// Env
		for _, env := range container.Run.Env() {
			args = append(args, "--env", env)
		}
		// Env file
		if len(container.Run.EnvFile()) > 0 {
			args = append(args, "--env-file", container.Run.EnvFile())
		}
		// Expose
		for _, expose := range container.Run.Expose() {
			args = append(args, "--expose", expose)
		}
		// Host
		if len(container.Run.Hostname()) > 0 {
			args = append(args, "--hostname", container.Run.Hostname())
		}
		// Interactive
		if container.Run.Interactive {
			args = append(args, "--interactive")
		}
		// Link
		for _, link := range container.Run.Link() {
			args = append(args, "--link", link)
		}
		// LxcConf
		for _, lxcConf := range container.Run.LxcConf() {
			args = append(args, "--lxc-conf", lxcConf)
		}
		// Memory
		if len(container.Run.Memory()) > 0 {
			args = append(args, "--memory", container.Run.Memory())
		}
		// Net
		if container.Run.Net() != "bridge" {
			args = append(args, "--net", container.Run.Net())
		}
		// Privileged
		if container.Run.Privileged {
			args = append(args, "--privileged")
		}
		// Publish
		for _, port := range container.Run.Publish() {
			args = append(args, "--publish", port)
		}
		// PublishAll
		if container.Run.PublishAll {
			args = append(args, "--publish-all")
		}
		// Rm
		if container.Run.Rm {
			args = append(args, "--rm")
		}
		// Tty
		if container.Run.Tty {
			args = append(args, "--tty")
		}
		// User
		if len(container.Run.User()) > 0 {
			args = append(args, "--user", container.Run.User())
		}
		// Volumes
		for _, volume := range container.Run.Volume() {
			args = append(args, "--volume", volume)
		}
		// VolumesFrom
		for _, volumeFrom := range container.Run.VolumesFrom() {
			args = append(args, "--volumes-from", volumeFrom)
		}
		// Workdir
		if len(container.Run.Workdir()) > 0 {
			args = append(args, "--workdir", container.Run.Workdir())
		}
		// Name
		args = append(args, "--name", container.Name())
		// Image
		args = append(args, container.Image())
		// Command
		args = append(args, container.Run.Cmd()...)
		// Execute command
		executeCommand("docker", args)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Start container
func (container Container) start() {
	if container.exists() {
		if !container.running() {
			fmt.Printf("Starting container %s ... ", container.Name())
			args := []string{"start", container.Name()}
			executeCommand("docker", args)
		}
	} else {
		print.Error("Container %s does not exist.\n", container.Name())
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Kill container
func (container Container) kill() {
	if container.running() {
		fmt.Printf("Killing container %s ... ", container.Name())
		args := []string{"kill", container.Name()}
		executeCommand("docker", args)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Stop container
func (container Container) stop() {
	if container.running() {
		fmt.Printf("Stopping container %s ... ", container.Name())
		args := []string{"stop", container.Name()}
		executeCommand("docker", args)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
// Remove container
func (container Container) rm() {
	if container.exists() {
		if container.running() {
			print.Error("Container %s is running and cannot be removed.\n", container.Name())
		} else {
			fmt.Printf("Removing container %s ... ", container.Name())
			args := []string{"rm", container.Name()}
			executeCommand("docker", args)
		}
	}
}
