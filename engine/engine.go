package engine

import (
	"errors"
	"github.com/tsuru/docker-cluster/cluster"
)

type Engine struct {
	loader     ManifestLoader
	containers Containers
	cluster    *cluster.Cluster
}

type EngineFunc func(cont Containers) error

///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadFromFile(path, group string) error {
	man, err := e.loader.LoadFromFile(path)
	if err != nil {
		return err
	}

	//TODO Do we have nodes
	names := man.GetContainerNamesByGroup(group)
	//, fmt.Errorf("no group nor container matching `%s`", group)
	if err != nil {
		return err
	}

	conts := man.GetContainersByNames(names)
	err = e.cluster.Register(man.Nodes)
	if err != nil {
		return err
	}

	e.containers.Apply(conts)
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
func NewEngine() (*Engine, error) {
	eng := &Engine{loader: ManifestLoader{}}
	//TODO define storage, scheduler
	cluster, err := cluster.New(nil, nil)
	if err != nil {
		return nil, err
	}
	eng.cluster = cluster
	return eng, nil
}
