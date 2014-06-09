package engine

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/denkhaus/tcgl/applog"
	"github.com/fsouza/go-dockerclient"
	"math"
)

var (
	errCircularDependency = errors.New("Manifest error:: circular dependency detected")
)

type ContainerAggregateFunc func(e *list.Element, val interface{}) interface{}
type ContainerFunc func(e *list.Element) error

// Node represents a host running Docker. Each node has an ID and an address
// (in the form <scheme>://<host>:<port>/).
type Node struct {
	id         string
	manifestId string
	address    string
	engine     *Engine
	tree       *Tree
	client     *docker.Client
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) String() string {
	//ret := fmt.Sprintf("id:%s, address:%s, containers: %d", n.id, n.address, n.tree.Length())
	//ret = n.Aggregate(ret, func(cont Container, val interface{}) interface{} {
	//	res := val.(string)
	//	res += fmt.Sprintf("%s, ", cont.name)
	//	return res
	//}).(string)
	return n.id
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Info1
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Info1() string {
	ret := fmt.Sprintf("id:%s, address:%s", n.id, n.address)
	return ret
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// NewNode
/////////////////////////////////////////////////////////////////////////////////////////////////
func NewNode(id, address, manId string, e *Engine, tmps []Template) (*Node, error) {
	node := &Node{id: id, address: address, manifestId: manId, engine: e, tree: NewTree()}

	if client, err := docker.NewClient(address); err != nil {
		return nil, err
	} else {
		node.client = client
	}

	applog.Infof("Apply container templates -> [%s]", node.Info1())
	if err := node.Apply(tmps); err != nil {
		return nil, err
	}

	return node, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ApplyState(cnts []docker.APIContainers) error {
	for _, apiCnt := range cnts {
		for _, name := range apiCnt.Names {
			n.ForAll(func(e *list.Element) error {
				cnt := e.Value.(Container)
				if name[1:] == cnt.FullQualifiedName() { // trim "/"
					applog.Debugf("Apply Id of container %s on node [%s]", cnt, n)
					cnt.SetId(apiCnt.ID)
					e.Value = cnt
				}
				return nil
			})
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Apply(tmps []Template) error {
	tree := n.tree

	for _, tmp := range tmps {
		cont, err := NewContainerFromTemplate(tmp, &n)
		if err != nil {
			return err
		}

		cont.RemoveSelfReference()
		//check if container is required
		var requiredIns *list.Element
		nRequiredIdx := math.MaxInt64
		for e := tree.First(); e != nil; e = e.Next() {
			cnt := e.Value.(Container)
			for _, name := range cnt.reqmnts {
				if name == cont.name {
					idx := tree.GetIndex(cnt)
					if idx < nRequiredIdx {
						requiredIns = cnt.elm
						nRequiredIdx = idx
					}
				}
			}
		}

		deps := cont.reqmnts
		if len(deps) == 0 && requiredIns == nil {
			applog.Debugf("Apply template %s - no requirements, not required", cont.name)
			//if is not required and has no requirements

			cont.elm = tree.TreePushBack(*cont)
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

		if hasRequirementsIns != nil && requiredIns != nil &&
			nRequiredIdx <= nHasRequirementsIdx {
			return errCircularDependency
		}

		applog.Debugf("Apply template %s: reqIdx: %d, hasReqIdx: %d", cont.name, nRequiredIdx, nHasRequirementsIdx)

		if hasRequirementsIns != nil && requiredIns == nil { // only has requirements
			cont.elm = tree.TreeInsertAfter(*cont, hasRequirementsIns)
		} else if hasRequirementsIns == nil && requiredIns != nil { //only is required
			cont.elm = tree.TreeInsertBefore(*cont, requiredIns)
		} else if nRequiredIdx > nHasRequirementsIdx {

		}
	}

	applog.Debugf("Apply templates:: Building node with %d container(s) finished.", tree.Length())
	//TODO check unmet requirements
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) HasContainers() bool {
	return n.tree.Length() != 0
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Aggregate(val interface{}, fn ContainerAggregateFunc) interface{} {
	for e := n.tree.First(); e != nil; e = e.Next() {
		val = fn(e, val)
	}
	return val
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ForAll(fn ContainerFunc) error {
	for e := n.tree.First(); e != nil; e = e.Next() {
		if err := fn(e); err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) ForAllReversed(fn ContainerFunc) error {
	for e := n.tree.Last(); e != nil; e = e.Prev() {
		if err := fn(e); err != nil {
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
func (n Node) Lift(force bool, kill bool) error {
	if err := n.Provision(force); err != nil {
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
func (n Node) Provision(force bool) error {
	err := n.ForAll(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.provision(force)
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Run containers.
// When forced, removes existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Run(force bool, kill bool) error {
	if force {
		if err := n.Remove(force, kill); err != nil {
			return err
		}
	}
	err := n.ForAll(func(e *list.Element) error {
		cnt := e.Value.(Container)
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
		if err := n.Remove(force, kill); err != nil {
			return err
		}
	}
	err := n.ForAll(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.runOrStart()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Start containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Start() error {
	err := n.ForAll(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.start()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Kill containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Kill() error {
	err := n.ForAllReversed(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.kill()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Stop containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Stop() error {
	err := n.ForAllReversed(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.stop()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Remove containers.
// When forced, stops existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Remove(force bool, kill bool) error {
	if force {
		if kill {
			if err := n.Kill(); err != nil {
				return err
			}
		} else {
			if err := n.Stop(); err != nil {
				return err
			}
		}
	}
	err := n.ForAllReversed(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.remove()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Status of containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Node) Status() error {
	err := n.ForAll(func(e *list.Element) error {
		cnt := e.Value.(Container)
		return cnt.status()
	})
	return err
}
