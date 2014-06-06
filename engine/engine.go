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

type EngineFunc func(cont Containers) error

///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadFromFile(path, group string) error {
	applog.Infof("Load manifest frome file %s", path)

	man, err := e.loader.LoadFromFile(path)
	if err != nil {
		return err
	}

	if len(man.Nodes) == 0 {
		return errEmptyNodes
	}

	e.Nodes = man.Nodes
	if err := e.Nodes.Initialize(e, man, group); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Execute
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	if e.containers.IsEmpty() {
		return errors.New("no containers loaded")
	}

	return fn(e.containers)
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
