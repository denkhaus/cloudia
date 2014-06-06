package engine

import (
	"errors"
	"github.com/denkhaus/tcgl/applog"
	"github.com/tsuru/docker-cluster/cluster"
	"github.com/tsuru/docker-cluster/storage"
)

var (
	errRedisAddressNotSpecified = errors.New("config error:: Please specify redis storage address as <scheme>://<host>:<port>.")
	errEmptyNodes               = errors.New("manifest error:: Please specify at least one node.")
)

type Engine struct {
	loader  ManifestLoader
	nodes   Nodes
	cluster *cluster.Cluster
}

type EngineFunc func(cont Node) error

///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadFromFile(path, group string) error {
	applog.Infof("Load manifest from file %s", path)

	man, err := e.loader.LoadFromFile(path)
	if err != nil {
		return err
	}

	if len(man.Nodes) == 0 {
		return errEmptyNodes
	}

	e.nodes = man.Nodes
	if err := e.nodes.Initialize(e, man, group); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Execute
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	err := e.nodes.ForAll(func(node Node) error {
		if !node.HasContainers() {
			return errors.New("node error:: No containers loaded")
		}
		return fn(node)
	})
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
// NewEngine
///////////////////////////////////////////////////////////////////////////////////////////////
func NewEngine(storageAddress, storagePassword, storagePrefix string) (*Engine, error) {
	eng := &Engine{loader: ManifestLoader{}}

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
