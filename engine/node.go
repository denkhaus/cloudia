package engine

import (
	"github.com/denkhaus/tcgl/applog"
)

// Node represents a host running Docker. Each node has an ID and an address
// (in the form <scheme>://<host>:<port>/).
type Node struct {
	Id      string `json:"id" yaml:"id"`
	Address string `json:"address" yaml:"address"`
	cnts    Containers
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// ToMap
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ToMap() map[string]string {
	mp := make(map[string]string)
	mp["Address"] = node.Address
	mp["ID"] = node.Id
	return mp
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) String() string {
	return fmt.Sprintf("id:%s, address:%s", n.Id, n.Address)
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Initialize
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Initialize(c *cluster.Cluster, cnts []Container) error {
	applog.Infof("Node %s :: Apply container template", n)
	if err := n.cnts.Apply(cnts); err != nil {
		return err
	}
	applog.Infof("Node %s :: Register with cluster", n)
	if err = e.cluster.Register(n.ToMap()); err != nil {
		return err
	}

	return nil
}
