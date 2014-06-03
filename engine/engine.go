package engine

import (
	"errors"
	"github.com/tsuru/docker-cluster/cluster"
)

type Engine struct {
	loader     ManifestLoader
	containers Containers
	cluster    cluster.Cluster
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
	names, err := man.GetContainerNamesByGroup(group)
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
		return errors.New("No containers loaded")
	}

	return fn(e.containers)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// NewEngine
///////////////////////////////////////////////////////////////////////////////////////////////
func NewEngine() *Engine {
	eng := &Engine{loader: ManifestLoader{}}
	eng.cluster = cluster.New()
	return eng
}
