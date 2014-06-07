package engine

import (
	"errors"
	//"fmt"
	"github.com/denkhaus/tcgl/applog"
	"github.com/tsuru/docker-cluster/cluster"
	"github.com/tsuru/docker-cluster/storage"
)

var (
	errRedisAddressNotSpecified = errors.New("Storage error:: Please specify redis storage address as <scheme>://<host>:<port>.")
	errEmptyNodes               = errors.New("Manifest error:: Please specify at least one node.")
	errEmptyContainerNames      = errors.New("Group error:: No containers available by provided group.")
)

type Engine struct {
	nodes   []Node
	cluster *cluster.Cluster
}

type EngineFunc func(cont Node) error

//type NodeAggregateFunc func(node Node, val interface{}) interface{}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
//func (n Nodes) Aggregate(val interface{}, fn NodeAggregateFunc) interface{} {

//	for _, node := range n {
//		val = fn(node, val)
//	}
//	return val
//}

/////////////////////////////////////////////////////////////////////////////////////////////////
// String
/////////////////////////////////////////////////////////////////////////////////////////////////
//func (n Nodes) String() string {
//	var ret string
//	ret = n.Aggregate(ret, func(node Node, val interface{}) interface{} {
//		res := val.(string)
//		res += fmt.Sprintf("[%s]\n", node)
//		return res
//	}).(string)

//	return ret
//}

///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadFromFile(path, group string) error {
	applog.Infof("Load manifest from file %s", path)

	ld := ManifestLoader{}
	man, err := ld.LoadFromFile(path)
	if err != nil {
		return err
	}

	return e.processManifest(man, group)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadDefaults(group string) error {
	applog.Infof("Load default manifest.")

	ld := ManifestLoader{}
	man, err := ld.LoadDefault()
	if err != nil {
		return err
	}

	return e.processManifest(man, group)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// processManifest
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) processManifest(man *Manifest, group string) error {
	if len(man.Nodes) == 0 {
		return errEmptyNodes
	}

	names := man.GetTemplateNamesByGroup(group)
	if len(names) == 0 {
		return errEmptyContainerNames
	}

	tmps := man.GetTemplatesByNames(names)
	for _, cn := range man.Nodes {
		node, err := NewNode(cn.Id, cn.Address, e, tmps)
		if err != nil {
			return err
		}
		e.nodes = append(e.nodes, *node)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) ForAllNodes(fn EngineFunc) error {
	for _, node := range e.nodes {
		if err := fn(node); err != nil {
			return err
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Register
/////////////////////////////////////////////////////////////////////////////////////////////////
//func (n Nodes) Initialize(e *Engine, man *Manifest, group string) error {
//	//applog.Infof("Initialize node(s) %s", n)

//	for _, node := range n {
//		if err := node.Initialize(e.cluster, cnts); err != nil {
//			return err
//		}
//	}

//	//applog.Infof("Apply group --> %s process containers --> %s", group, e.containers)
//	return nil
//}

///////////////////////////////////////////////////////////////////////////////////////////////
// Execute
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	err := e.ForAllNodes(func(node Node) error {
		if !node.HasContainers() {
			return errors.New("Node error:: No containers loaded")
		}
		return fn(node)
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// NewEngine
///////////////////////////////////////////////////////////////////////////////////////////////
func NewEngine(storageAddress, storagePassword, storagePrefix string) (*Engine, error) {
	eng := &Engine{}

	applog.Debugf("Create cluster store")
	if len(storageAddress) == 0 {
		return nil, errRedisAddressNotSpecified
	}

	var store cluster.Storage
	if len(storagePassword) > 0 {
		store = storage.AuthenticatedRedis(storageAddress, storagePassword, storagePrefix)
	} else {
		store = storage.Redis(storageAddress, storagePrefix)
	}

	//TODO define scheduler ?
	clst, err := cluster.New(nil, store)
	if err != nil {
		return nil, err
	}

	eng.cluster = clst
	return eng, nil
}
