package engine

import (
	"container/list"
	"github.com/denkhaus/tcgl/applog"
	"github.com/fsouza/go-dockerclient"
)

type Container struct {
	id                   string
	elm                  *list.Element
	node                 *Node
	response             *docker.Container
	hostConfig           *docker.HostConfig
	config               *docker.Config
	name                 string
	image                string
	reqmnts              []string
	stopContainerTimeout uint
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func NewContainerFromTemplate(tmp Template, node *Node) (*Container, error) {
	cnt := &Container{node: node}

	cnt.name = tmp.Name()
	cnt.image = tmp.Image()
	cnt.reqmnts = tmp.Requirements

	cnf, err := tmp.Run.CreateDockerConfig()
	if err != nil {
		return nil, err
	}
	cnt.config = cnf

	hstcnf, err := tmp.Run.CreateDockerHostConfig()
	if err != nil {
		return nil, err
	}
	cnt.hostConfig = hstcnf

	cnt.config.Image = cnt.image
	cnt.stopContainerTimeout = 30 //seconds
	return cnt, nil
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
		return cnt.create()
	} else {
		applog.Infof("Container %s does already exist at node\n%s. Use --force to recreate.\n", cnt.image, cnt.node)
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
	applog.Infof("Creating container %s on node %s ", cnt.name, cnt.node)

	clst := cnt.node.engine.cluster
	opts := docker.CreateContainerOptions{Config: cnt.config, Name: cnt.name}
	id, newCont, err := clst.CreateContainer(opts, cnt.node.id)
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
		applog.Infof("Container %s does already exist on node\n%s. Use --force to recreate.", cnt.name, cnt.node)
		if !cnt.running() {
			err := cnt.start()
			if err != nil {

			}
		}
	} else {
		err := cnt.create()
		if err != nil {

		}

	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Start container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) start() error {
	if container.exists() {
		if !container.running() {
			applog.Infof("Attempt to start container %s on node %s", cnt.name, cnt.node)
			clst := cnt.node.engine.cluster
			return clst.StartContainer(cnt.id, cnt.hostConfig)
		} else {
			applog.Infof("Attempt to start container %s on node %s not successfull. Container is already running", cnt.name, cnt.node)
		}
	} else {
		applog.Infof("Attempt to start container %s on node %s not successfull. Container does not exist", cnt.name, cnt.node)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Kill container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt Container) kill() error {
	if cnt.running() {
		applog.Infof("Attempt to kill container %s on node %s", cnt.name, cnt.node)
		opts := docker.KillContainerOptions{ID: cnt.id}
		clst := cnt.node.engine.cluster
		return clst.KillContainer(opts)
	} else {
		applog.Infof("Attempt to kill container %s on node %s not successfull. Container is not running", cnt.name, cnt.node)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Stop container
/////////////////////////////////////////////////////////////////////////////////////////////////

func (cnt Container) stop() error {
	if cnt.running() {
		applog.Infof("Attempt to stop container %s on node %s", cnt.name, cnt.node)
		clst := cnt.node.engine.cluster
		return clst.StopContainer(cnt.id, cnt.stopContainerTimeout)
	} else {
		applog.Infof("Attempt to stop container %s on node %s not successfull. Container is not running", cnt.name, cnt.node)
	}

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
