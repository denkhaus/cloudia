package engine

import (
	"container/list"
	"fmt"
	"os"
	"text/tabwriter"
)

type ContainerFunc func(cont Container) error
type Containers struct {
	engine *Engine
	tree   *Tree
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) Apply(conts []Container) error {

	tree := NewTree()
	for cont := range conts {
		//check if container is required
		var requiredIns *Element
		nRequiredIdx := math.MaxInt64
		for e := c.tree.Front(); e != nil; e = e.Next() {
			cnt := e.Value.(Container)
			for name := range cnt.Dependencies {
				if name == cont.Name() {
					idx := t.GetIndex(cnt)
					if idx < nRequiredIdx {
						requiredIns = cnt.elm
						nRequiredIdx = idx
					}
				}
			}
		}

		deps := cont.Dependencies
		if len(deps) == 0 && requiredIns == nil {
			//if is not required and has no requirements
			tree.TreePushBack(cont)
			continue
		}

		//check if container has requirements
		var hasRequirementsIns *Element
		nHasRequirementsIdx := math.MaxInt64
		for name := range deps {
			cnt := tree.GetContainerByName(name)
			if cnt != nil {
				idx := t.GetIndex(cnt)
				if idx < nHasRequirementsIdx {
					hasRequirementsIns = cnt.elm
					nHasRequirementsIdx = idx
				}
			} else {
				tree.AddUnmetDependency(name)
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

	c.tree = tree
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) IsEmpty() bool {
	return c.tree.Len() == 0
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) ForAll(fn ContainerFunc) error {
	for e := c.tree.Front(); e != nil; e = e.Next() {
		if err := fn(e.Value); err != nil {
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
func (c Containers) lift(force bool, kill bool) {
	if err := c.provision(force); err != nil {
		return err
	}
	if err := c.runOrStart(force, kill); err != nil {
		return err
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Provision containers.
// When forced, this will rebuild all images.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) provision(force bool) error {
	err := c.ForAll(func(cont Container) error {
		return cont.provision(force)
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Run containers.
// When forced, removes existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) run(force bool, kill bool) error {
	if force {
		if err := c.rm(force, kill); err != nil {
			return err
		}
	}
	err := c.ForAll(func(cont Container) error {
		return cont.run()
	})

	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Run or start containers.
// When forced, removes existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) runOrStart(force bool, kill bool) error {
	if force {
		if err := c.rm(force, kill); err != nil {
			return err
		}
	}
	err := c.ForAll(func(cont Container) error {
		return cont.runOrStart()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Start containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) start() error {
	err := c.ForAll(func(cont Container) error {
		return cont.start()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Kill containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) kill() error {
	err := c.ForAll(func(cont Container) error {
		return cont.kill()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Stop containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) stop() error {
	err := c.ForAll(func(cont Container) error {
		return cont.stop()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Remove containers.
// When forced, stops existing containers first.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) rm(force bool, kill bool) error {
	if force {
		if kill {
			if err := c.kill(); err != nil {
				return err
			}
		} else {
			if err := c.stop(); err != nil {
				return err
			}
		}
	}
	err := c.ForAll(func(cont Container) error {
		return cont.rm()
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Status of containers.
///////////////////////////////////////////////////////////////////////////////////////////////
func (c Containers) status() error {
	err := c.ForAll(func(cont Container) error {
		return cont.status()
	})
	return err
}
