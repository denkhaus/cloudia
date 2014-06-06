package engine

import (
	"errors"
	"fmt"
	"github.com/denkhaus/tcgl/applog"
)

var (
	errEmptyContainerNames = errors.New("group error:: No containers available by provided group.")
)

type Nodes []Node
type NodeFunc func(node Node) error
type NodeAggregateFunc func(node Node, val interface{}) interface{}

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) String() string {
	var ret string
	ret = n.Aggregate(ret, func(node Node, val interface{}) interface{} {
		ret := val.(string)
		ret += fmt.Sprintf("[%s]\n", node)
		return ret
	}).(string)

	return ret
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) ForAll(fn NodeFunc) error {
	for _, node := range n {
		if err := fn(node); err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) Aggregate(val interface{}, fn NodeAggregateFunc) interface{} {
	for _, node := range n {
		val = fn(node, val)
	}
	return val
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Register
/////////////////////////////////////////////////////////////////////////////////////////////////
func (n Nodes) Initialize(e *Engine, man *Manifest, group string) error {
	applog.Infof("Initialize node(s) %s", n)

	names := man.GetContainerNamesByGroup(group)
	if len(names) == 0 {
		return errEmptyContainerNames
	}

	cnts := man.GetContainersByNames(names)
	for _, node := range n {
		if err := node.Initialize(e.cluster, cnts); err != nil {
			return err
		}
	}

	//applog.Infof("Apply group --> %s process containers --> %s", group, e.containers)
	return nil
}
