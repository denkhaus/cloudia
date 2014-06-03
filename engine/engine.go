package engine

import (
	"github.com/tsuru/docker-cluster/cluster"
)

type Engine struct {
	loader     ManifestLoader
	containers Containers
	cluster    cluster.Cluster
}

type EngineFunc func(containers Containers) error
///////////////////////////////////////////////////////////////////////////////////////////////
// LoadFromFile
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) LoadFromFile(path, group string) error {
	man, err := e.loader.Load(path)
	if err != nil {
		e.containers = nil
		return err
	}
	//TODO Do we have nodes
	names, err := man.GetContainerNamesByGroup(group)
	if err != nil {
		e.containers = nil
		return err
	}

	conts = man.GetContainersByNames(names)	
	err =: e.cluster.Register(man.Nodes)	
	if err != nil {
		e.containers = nil
		return err
	}
	
	e.containers = conts
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Execute
///////////////////////////////////////////////////////////////////////////////////////////////
func (e *Engine) Execute(fn EngineFunc) error {
	if e.containers == nil {
		return StatusError{errors.New("No containers loaded"), 22}
	}

	return fn(e.containers)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// NewEngine
///////////////////////////////////////////////////////////////////////////////////////////////
func NewEngine() *Engine {
	eng := &Engine{loader: manifest.ManifestLoader{}}
	eng.cluster = cluster.New()
	return eng
}
