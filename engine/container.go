package engine

import (
	"container/list"
	"fmt"
	"github.com/denkhaus/tcgl/applog"
	"github.com/fsouza/go-dockerclient"
)

type Container struct {
	id                           string
	manifestId                   string
	elm                          *list.Element
	node                         *Node
	response                     *docker.Container
	hostConfig                   *docker.HostConfig
	config                       *docker.Config
	name                         string
	image                        string
	reqmnts                      []string
	stopContainerTimeout         uint
	removeContainerForce         bool
	removeContainerRemoveVolumes bool
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func NewContainerFromTemplate(tmp Template, node *Node) (*Container, error) {
	cnt := &Container{node: node, manifestId: node.manifestId}

	cnt.removeContainerForce = false
	cnt.removeContainerRemoveVolumes = false
	cnt.stopContainerTimeout = 10 //seconds

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
	return cnt, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//// RemoveSelfReference
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) RemoveSelfReference() {
	var rq []string
	for _, ref := range cnt.reqmnts {
		if cnt.name != ref {
			rq = append(rq, ref)
		}
	}
	cnt.reqmnts = rq
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//// SetId
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) SetId(id string) {
	cnt.id = id
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//// FullQualifiedName
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) FullQualifiedName() string {
	return fmt.Sprintf("%s-%s", cnt.name, cnt.manifestId)
}

////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) String() string {
	return cnt.name
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) retrieveState() bool {

	if len(cnt.id) == 0 {
		applog.Errorf("Container error:: cannot retrieve infos -> id not set")
		return false
	}

	client := cnt.node.client
	cont, err := client.InspectContainer(cnt.id)
	if err != nil {
		applog.Errorf("Container error:: while checking status -> %s", err.Error())
		return false
	}
	cnt.response = cont
	return true
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) exists() bool {
	return cnt.retrieveState()
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) running() bool {
	if ok := cnt.retrieveState(); ok {
		state := cnt.response.State
		return state.Running
	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) status() error {
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
////
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) provision(force bool) error {
	if force || !cnt.exists() {
		if err := cnt.create(); err != nil {
			applog.Infof("Creating container %s on node %s not successfull.", cnt.name, cnt.node)
			return err
		} else {
			applog.Infof("Container %s successfull created on node %s.", cnt.name, cnt.node)
		}
	} else {
		applog.Infof("Container %s already exists on node %s.", cnt.name, cnt.node)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Run or start container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) runOrStart() error {
	if cnt.exists() {
		return cnt.start()
	} else {
		return cnt.run()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Create container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) create() error {

	client := cnt.node.client
	opts := docker.CreateContainerOptions{Config: cnt.config, Name: cnt.FullQualifiedName()}

	applog.Infof("Pull image %s for container %s on node %s ", opts.Config.Image, cnt, cnt.node)
	err := client.PullImage(docker.PullImageOptions{Repository: opts.Config.Image}, docker.AuthConfiguration{})
	if err != nil {
		return err
	}

	applog.Infof("Creating container %s on node %s. Please wait...", cnt, cnt.node)
	newCont, err := client.CreateContainer(opts)
	if err != nil {
		return err
	}

	cnt.id = newCont.ID
	cnt.response = newCont
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Run container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) run() error {
	if cnt.exists() {
		if !cnt.running() {
			return cnt.start()
		}
	} else {
		if err := cnt.create(); err != nil {
			return err
		}
		return cnt.start()
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Start container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) start() error {
	if cnt.exists() {
		if !cnt.running() {
			client := cnt.node.client
			if err := client.StartContainer(cnt.id, cnt.hostConfig); err != nil {
				applog.Infof("Starting container %s on node %s not successfull.", cnt.name, cnt.node)
				return err
			}
			applog.Infof("Container %s successfull started on node %s.", cnt.name, cnt.node)

		} else {
			applog.Infof("Container %s is already running on node %s.", cnt.name, cnt.node)
		}
	} else {
		applog.Infof("Container %s does not exist on node %s.", cnt.name, cnt.node)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Kill container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) kill() error {
	if cnt.exists() {
		if cnt.running() {
			opts := docker.KillContainerOptions{ID: cnt.id}
			client := cnt.node.client
			if err := client.KillContainer(opts); err != nil {
				applog.Infof("Killing container %s on node %s not successfull.", cnt.name, cnt.node)
				return err
			} else {
				applog.Infof("Container %s successfull killed on node %s.", cnt.name, cnt.node)
			}
		} else {
			applog.Infof("Container %s is not running on node %s.", cnt.name, cnt.node)
		}
	} else {
		applog.Infof("Container %s does not exist on node %s.", cnt.name, cnt.node)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Stop container
/////////////////////////////////////////////////////////////////////////////////////////////////

func (cnt *Container) stop() error {
	if cnt.exists() {
		if cnt.running() {
			client := cnt.node.client
			if err := client.StopContainer(cnt.id, cnt.stopContainerTimeout); err != nil {
				applog.Infof("Stopping container %s on node %s not successfull.", cnt.name, cnt.node)
				return err
			} else {
				applog.Infof("Container %s successfull stopped on node %s.", cnt.name, cnt.node)
			}

		} else {
			applog.Infof("Container %s is not running on node %s.", cnt.name, cnt.node)
		}
	} else {
		applog.Infof("Container %s does not exist on node %s.", cnt.name, cnt.node)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//Remove container
/////////////////////////////////////////////////////////////////////////////////////////////////
func (cnt *Container) remove() error {
	if cnt.exists() {
		if cnt.running() {
			applog.Infof("Removing container %s on node %s not successfull. Container is running", cnt.name, cnt.node)
		} else {
			opts := docker.RemoveContainerOptions{ID: cnt.id}
			opts.Force = cnt.removeContainerForce
			opts.RemoveVolumes = cnt.removeContainerRemoveVolumes

			client := cnt.node.client
			if err := client.RemoveContainer(opts); err != nil {
				applog.Infof("Removing container %s on node %s not successfull.", cnt.name, cnt.node)
				return err
			} else {
				applog.Infof("Container %s successfull removed on node %s.", cnt.name, cnt.node)
			}
		}
	} else {
		applog.Infof("Container %s does not exist on node %s.", cnt.name, cnt.node)
	}
	return nil
}
