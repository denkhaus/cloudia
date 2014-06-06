package engine

import (
	"container/list"
	"fmt"
	"github.com/denkhaus/tcgl/applog"
	"github.com/tsuru/docker-cluster/cluster"
	"math"
)

type ContainerAggregateFunc func(cont Container, val interface{}) interface{}
type ContainerFunc func(cont Container) error

// Node represents a host running Docker. Each node has an ID and an address
// (in the form <scheme>://<host>:<port>/).
type Node struct {
	Id      string `json:"id" yaml:"id"`
	Address string `json:"address" yaml:"address"`
	engine  *Engine
	tree    *Tree
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// ToMap
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ToMap() map[string]string {
	mp := make(map[string]string)
	mp["Address"] = n.Address
	mp["ID"] = n.Id
	return mp
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) String() string {

	ret := fmt.Sprintf("id:%s, address:%s\n", n.Id, n.Address)
	ret = n.Aggregate(ret, func(cont Container, val interface{}) interface{} {
		ret := val.(string)
		ret += fmt.Sprintf("%s, ", cont.Name())
		return ret
	}).(string)
	return ret
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Initialize
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Initialize(c *cluster.Cluster, cnts []Container) error {
	applog.Infof("Node %s :: Apply container template", n)
	if err := n.Apply(cnts); err != nil {
		return err
	}
	applog.Infof("Node %s :: Register with cluster", n)
	if err := c.Register(n.ToMap()); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Apply(conts []Container) error {

	tree := NewTree()
	for _, cont := range conts {
		//check if container is required
		var requiredIns *list.Element
		nRequiredIdx := math.MaxInt64
		for e := tree.Front(); e != nil; e = e.Next() {
			cnt := e.Value.(Container)
			for _, name := range cnt.Requirements {
				if name == cont.Name() {
					idx := tree.GetIndex(cnt)
					if idx < nRequiredIdx {
						requiredIns = cnt.elm
						nRequiredIdx = idx
					}
				}
			}
		}

		deps := cont.Requirements
		if len(deps) == 0 && requiredIns == nil {
			//if is not required and has no requirements
			tree.TreePushBack(cont)
			continue
		}

		//check if container has requirements
		var hasRequirementsIns *list.Element
		nHasRequirementsIdx := math.MaxInt64
		for _, name := range deps {
			c := tree.GetContainerByName(name)
			if cnt, ok := c.(Container); ok {
				idx := tree.GetIndex(cnt)
				if idx < nHasRequirementsIdx {
					hasRequirementsIns = cnt.elm
					nHasRequirementsIdx = idx
				}
			} else {
				tree.AddUnmetRequirement(name)
			}
		}

		// try to insert before nRequiredIdx and after nHasRequirementsIdx

		if hasRequirementsIns != nil && requiredIns == nil { // only has requirements
			tree.TreeInsertAfter(cont, hasRequirementsIns)
		} else if hasRequirementsIns == nil && requiredIns != nil { //only is required
			tree.TreeInsertBefore(cont, requiredIns)
		} else {
			if nRequiredIdx <= nHasRequirementsIdx {
				//
			}
		}
	}

	//TODO check unmet requirements
	n.tree = tree
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) HasContainers() bool {
	return n.tree.Len() != 0
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) First() *list.Element {
	return n.tree.Front()
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Last() *list.Element {
	return n.tree.Back()
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Aggregate(val interface{}, fn ContainerAggregateFunc) interface{} {
	for e := n.First(); e != nil; e = e.Next() {
		val = fn(e.Value.(Container), val)
	}
	return val
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ForAll(fn ContainerFunc) error {
	for e := n.tree.Front(); e != nil; e = e.Next() {
		if err := fn(e.Value.(Container)); err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ForAllReversed(fn ContainerFunc) error {
	for e := n.Last(); e != nil; e = e.Prev() {
		if err := fn(e.Value.(Container)); err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Lift containers (provision + run).
// When forced, this will rebuild all images
// and recreate all containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) lift(force bool, kill bool) error {
	if err := n.provision(force); err != nil {
		return err
	}
	if err := n.runOrStart(force, kill); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Provision containers.
// When forced, this will rebuild all images.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) provision(force bool) error {
	err := n.ForAll(func(cnt Container) error {
		return cnt.provision(force)
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Run containers.
// When forced, removes existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) run(force bool, kill bool) error {
	if force {
		if err := n.remove(force, kill); err != nil {
			return err
		}
	}
	err := n.ForAll(func(cnt Container) error {
		return cnt.run()
	})

	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Run or start containers.
// When forced, removes existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) runOrStart(force bool, kill bool) error {
	if force {
		if err := n.remove(force, kill); err != nil {
			return err
		}
	}
	err := n.ForAll(func(cnt Container) error {
		return cnt.runOrStart()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Start containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) start() error {
	err := n.ForAll(func(cont Container) error {
		return cont.start()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Kill containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) kill() error {
	err := n.ForAllReversed(func(cont Container) error {
		return cont.kill()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Stop containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) stop() error {
	err := n.ForAllReversed(func(cnt Container) error {
		return cnt.stop()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Remove containers.
// When forced, stops existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) remove(force bool, kill bool) error {
	if force {
		if kill {
			if err := n.kill(); err != nil {
				return err
			}
		} else {
			if err := n.stop(); err != nil {
				return err
			}
		}
	}
	err := n.ForAllReversed(func(cnt Container) error {
		return cnt.remove()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Status of containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Status() error {
	err := n.ForAll(func(cont Container) error {
		return cont.status()
	})
	return err
}
