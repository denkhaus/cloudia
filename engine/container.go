package engine

import (
	"container/list"
	"github.com/denkhaus/tcgl/applog"
	"github.com/fsouza/go-dockerclient"
	"github.com/tsuru/docker-cluster/cluster"
	"os"
)

type Container struct {
	id       string
	elm      *list.Element
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
func (container *Container) exists() bool {
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
//func (container *Container) running() bool {
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
//}

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
func (container *Container) status() error {
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
func (container Container) provision(force bool) error {
	//	if force || !container.imageExists() {
	//		if len(container.Dockerfile()) > 0 {
	//			container.buildImage()
	//		} else {
	//			container.pullImage()
	//		}
	//	} else {
	//		print.Notice("Image %s does already exist. Use --force to recreate.\n", container.Image())
	//	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Run or start container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (container Container) runOrStart() error {
	if container.exists() {
		return container.start()
	} else {
		return container.run()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Create container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) create(clst *cluster.Cluster) error {
	config := docker.Config{}

	applog.Infof("Creating container %s ... ", cnt.Name())

	// CPU shares
	if cnt.Run.CpuShares > 0 {
		config.CpuSharesargs = strconv.Itoa(cnt.Run.CpuShares)
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
		config.Hostname = container.Run.Hostname()
	}
	// Memory
	if len(cnt.Run.Memory()) > 0 {
		config.Memory = cnt.Run.Memory()
	}
	// User
	if len(cnt.Run.User()) > 0 {
		config.User = cnt.Run.User()
	}
	// Volumes
	for _, volume := range cnt.Run.Volumes() {
		config.Volumes = append(config.Volumes, volume)
	}
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
	for _, cmd := range cnt.Run.Cmd() {
		config.Cmd = append(config.Cmd, cmd)
	}

	opts := docker.CreateContainerOptions{Config: &config}
	opts.Name = container.Name()

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
func (cnt Container) run(clst *cluster.Cluster) error {
	if cnt.exists() {
		applog.Infof("Container %s does already exist. Use --force to recreate.\n", cnt.Name())
		if !cnt.running() {
			cnt.Start()
		}
	} else {
		cnt.create(clst)

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

func (container Container) start() error {
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
func (cnt Container) kill(clst *cluster.Cluster) error {
	if cnt.running() {
		applog.Infof("Attempt to kill container %s ... ", cnt.Name())
		opts := docker.KillContainerOptions{ID: cnt.Id}
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
func (container Container) rm() error {
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
