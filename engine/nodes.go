package engine

import (
	"github.com/denkhaus/tcgl/applog"
)

var (
	errEmptyContainerNames = errors.New("group error:: No containers available by provided group.")
)

type Nodes []Node

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) String() string {
	var res string
	for node := range n {
		res = append(res, fmt.Sprintf("%s\n", node))
	}
	return res
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Register
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) Initialize(e *Engine, man *Manifest, group string) error {
	applog.Infof("Initialize nodes %s", n)

	names := man.GetContainerNamesByGroup(group)
	if len(names) == 0 {
		return errEmptyContainerNames
	}

	cnts := man.GetContainersByNames(names)
	for node := range n {
		if err := node.Initialize(e.cluster, cnts); err != nil {
			return err
		}
	}

	//applog.Infof("Apply group --> %s process containers --> %s", group, e.containers)
}
