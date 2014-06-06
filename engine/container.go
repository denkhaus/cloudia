package engine

import (
	"container/list"
	"github.com/denkhaus/tcgl/applog"
	"github.com/fsouza/go-dockerclient"
	"os"
	"strconv"
)

type Container struct {
	id       string
	elm      *list.Element
	node     *Node
	response *docker.Container

	RawName      string   `json:"name" yaml:"name"`
	RawImage     string   `json:"image" yaml:"image"`
	Requirements []string `json:"required" yaml:"required"`
	Run          RunParameters
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Name() string {
	return os.ExpandEnv(container.RawName)
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (container *Container) Image() string {
	return os.ExpandEnv(container.RawImage)
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
//func (container *Container) Id() (id string, err error) {
//	if len(container.id) > 0 {
//		id = container.id
//	} else {
//		// Inspect container, extracting the ID.
//		// This will return gibberish if no container is found.
//		args := []string{"inspect", "--format={{.ID}}", container.Name()}
//		output, outErr := commandOutput("docker", args)
//		if err == nil {
//			id = output
//			container.id = output
//		} else {
//			err = outErr
//		}
//	}
//	return
//}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) exists() bool {
	//	// `ps -a` returns all existant containers
	//	id, err := container.Id()
	//	if err != nil || len(id) == 0 {
	//		return false
	//	}
	//	dockerCmd := []string{"docker", "ps", "--quiet", "--all", "--no-trunc"}
	//	grepCmd := []string{"grep", "-wF", id}
	//	output, err := pipedCommandOutput(dockerCmd, grepCmd)
	//	if err != nil {
	//		return false
	//	}
	//	result := string(output)
	//	if len(result) > 0 {
	//		return true
	//	} else {
	//		return false
	//	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) running() bool {
	//	// `ps` returns all running containers
	//	id, err := container.Id()
	//	if err != nil || len(id) == 0 {
	//		return false
	//	}
	//	dockerCmd := []string{"docker", "ps", "--quiet", "--no-trunc"}
	//	grepCmd := []string{"grep", "-wF", id}
	//	output, err := pipedCommandOutput(dockerCmd, grepCmd)
	//	if err != nil {
	//		return false
	//	}
	//	result := string(output)
	//	if len(result) > 0 {
	//		return true
	//	} else {
	//		return false
	//	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
//func (container *Container) imageExists() bool {
//	dockerCmd := []string{"docker", "images", "--no-trunc"}
//	grepCmd := []string{"grep", "-wF", container.Image()}
//	output, err := pipedCommandOutput(dockerCmd, grepCmd)
//	if err != nil {
//		return false
//	}
//	result := string(output)
//	if len(result) > 0 {
//		return true
//	} else {
//		return false
//	}
//}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) status() error {
	//w * tabwriter.Writer
	//	args := []string{"inspect", "--format={{.State.Running}}\t{{.ID}}\t{{if .NetworkSettings.IPAddress}}{{.NetworkSettings.IPAddress}}{{else}}-{{end}}\t{{range $k,$v := $.NetworkSettings.Ports}}{{$k}},{{end}}", container.Name()}
	//	output, err := commandOutput("docker", args)
	//	if err != nil {
	//		fmt.Fprintf(w, "%s\tError:%v\t%v\n", container.Name(), err, output)
	//		return
	//	}
	//	fmt.Fprintf(w, "%s\t%s\n", container.Name(), output)
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
//// Pull image for container
//func (container *Container) pullImage() {
//	fmt.Printf("Pulling image %s ... ", container.Image())
//	args := []string{"pull", container.Image()}
//	executeCommand("docker", args)
//}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
//// Build image for container
//func (container *Container) buildImage() {
//	fmt.Printf("Building image %s ... ", container.Image())
//	args := []string{"build", "--rm", "--tag=" + container.Image(), container.Dockerfile()}
//	executeCommand("docker", args)
//}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) provision(force bool) error {
	if force || !cnt.exists() {
		cnt.create()
	} else {
		applog.Infof("Container %s does already exist. Use --force to recreate.\n", cnt.Image())
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Run or start container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) runOrStart() error {
	if cnt.exists() {
		return cnt.start()
	} else {
		return cnt.run()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Create container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) create() error {
	config := docker.Config{}

	applog.Infof("Creating container %s ... ", cnt.Name())

	// CPU shares
	if cnt.Run.CpuShares > 0 {
		config.CpuShares = int64(cnt.Run.CpuShares)
	}
	// Dns
	for _, dns := range cnt.Run.Dns() {
		config.Dns = append(config.Dns, dns)
	}
	// Env
	for _, env := range cnt.Run.Env() {
		config.Env = append(config.Env, env)
	}
	// Host
	if len(cnt.Run.Hostname()) > 0 {
		config.Hostname = cnt.Run.Hostname()
	}
	// Memory
	if len(cnt.Run.Memory()) > 0 {
		if mem, err := strconv.ParseInt(cnt.Run.Memory(), 10, 64); err == nil {
			config.Memory = mem
		} else {
			applog.Errorf("run parameters error:: Unable to convert memory param to int64 %v", err)
		}

	}
	// User
	if len(cnt.Run.User()) > 0 {
		config.User = cnt.Run.User()
	}

	// TODO Volumes
	//for _, volume := range cnt.Run.Volumes() {
	//	config.Volumes = append(config.Volumes, volume)
	//}

	// VolumesFrom
	if len(cnt.Run.VolumesFrom()) > 0 {
		config.VolumesFrom = cnt.Run.VolumesFrom()
	}
	// WorkingDir
	if len(cnt.Run.WorkingDir()) > 0 {
		config.WorkingDir = cnt.Run.WorkingDir()
	}
	// Image
	if len(cnt.Image()) > 0 {
		config.Image = cnt.Image()
	}
	// Command
	if cmds, err := cnt.Run.Cmd(); err != nil {
		for _, cmd := range cmds {
			config.Cmd = append(config.Cmd, cmd)
		}
	} else {
		applog.Errorf("run parameters error:: Errror while parsing cmd:: %v", err)
	}

	clst := cnt.node.engine.cluster
	opts := docker.CreateContainerOptions{Config: &config, Name: cnt.Name()}
	id, newCont, err := clst.CreateContainer(opts, "TODO nodesnodes ...string")
	if err != nil {
		return err
	}

	cnt.id = id
	cnt.response = newCont
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Run container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) run() error {
	if cnt.exists() {
		applog.Infof("Container %s does already exist. Use --force to recreate.\n", cnt.Name())
		if !cnt.running() {
			cnt.start()
		}
	} else {
		cnt.create()

		//		if len(container.Run.Cidfile()) > 0 {
		//			args = append(args, "--cidfile", container.Run.Cidfile())
		//		}

		//		// Detach
		//		if container.Run.Detach {
		//			args = append(args, "--detach")
		//		}

		//		// Entrypoint
		//		if len(container.Run.Entrypoint()) > 0 {
		//			args = append(args, "--entrypoint", container.Run.Entrypoint())
		//		}

		//		// Env file
		//		if len(container.Run.EnvFile()) > 0 {
		//			args = append(args, "--env-file", container.Run.EnvFile())
		//		}
		//		// Expose
		//		for _, expose := range container.Run.Expose() {
		//			args = append(args, "--expose", expose)
		//		}

		//		// Interactive
		//		if container.Run.Interactive {
		//			args = append(args, "--interactive")
		//		}
		//		// Link
		//		for _, link := range container.Run.Link() {
		//			args = append(args, "--link", link)
		//		}
		//		// LxcConf
		//		for _, lxcConf := range container.Run.LxcConf() {
		//			args = append(args, "--lxc-conf", lxcConf)
		//		}

		//		// Net
		//		if container.Run.Net() != "bridge" {
		//			args = append(args, "--net", container.Run.Net())
		//		}

		//		// Privileged
		//		if container.Run.Privileged {
		//			args = append(args, "--privileged")
		//		}
		//		// Publish
		//		for _, port := range container.Run.Publish() {
		//			args = append(args, "--publish", port)
		//		}
		//		// PublishAll
		//		if container.Run.PublishAll {
		//			args = append(args, "--publish-all")
		//		}
		//		// Rm
		//		if container.Run.Rm {
		//			args = append(args, "--rm")
		//		}
		//		// Tty
		//		if container.Run.Tty {
		//			args = append(args, "--tty")
		//		}

	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Start container
/////////////////////////////////////////////////////////////////////////////////////////////////

func (cnt Container) start() error {
	//	if container.exists() {
	//		if !container.running() {
	//			fmt.Printf("Starting container %s ... ", container.Name())
	//			args := []string{"start", container.Name()}
	//			executeCommand("docker", args)
	//		}
	//	} else {
	//		print.Error("Container %s does not exist.\n", container.Name())
	//	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Kill container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) kill() error {
	if cnt.running() {
		applog.Infof("Attempt to kill container %s ... ", cnt.Name())
		opts := docker.KillContainerOptions{ID: cnt.id}
		clst := cnt.node.engine.cluster
		return clst.KillContainer(opts)
	} else {
		applog.Infof("Attempt to kill container %s not successfull. Container is not running", cnt.Name())
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Stop container
/////////////////////////////////////////////////////////////////////////////////////////////////

func (container Container) stop() error {
	//	if container.running() {
	//		fmt.Printf("Stopping container %s ... ", container.Name())
	//		args := []string{"stop", container.Name()}
	//		executeCommand("docker", args)
	//	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//Remove container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (container Container) remove() error {
	//	if container.exists() {
	//		if container.running() {
	//			print.Error("Container %s is running and cannot be removed.\n", container.Name())
	//		} else {
	//			fmt.Printf("Removing container %s ... ", container.Name())
	//			args := []string{"rm", container.Name()}
	//			executeCommand("docker", args)
	//		}
	//	}
	return nil
}
