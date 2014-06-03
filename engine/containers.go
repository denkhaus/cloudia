package engine

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Containers struct {
	data   []Container
	engine *Engine
}

func (c Containers) Apply(data []Container) {
	//TODO Check dependencies, build dependency tree
	c.data = data
}

func (c Containers) IsEmpty() bool {
	return true
}

func (c Containers) reversed() []Container {
	var reversed []Container
	for i := len(c.data) - 1; i >= 0; i-- {
		reversed = append(reversed, c.data[i])
	}
	return reversed
}

// Lift containers (provision + run).
// When forced, this will rebuild all images
// and recreate all containers.
func (c Containers) lift(force bool, kill bool) {
	c.provision(force)
	c.runOrStart(force, kill)
}

// Provision containers.
// When forced, this will rebuild all images.
func (c Containers) provision(force bool) {
	for _, cont := range c.reversed() {
		//cont.provision(force)
	}
}

// Run containers.
// When forced, removes existing containers first.
func (c Containers) run(force bool, kill bool) {
	if force {
		c.rm(force, kill)
	}
	for _, cont := range c.reversed() {
		//cont.run()
	}
}

// Run or start containers.
// When forced, removes existing containers first.
func (c Containers) runOrStart(force bool, kill bool) {
	if force {
		c.rm(force, kill)
	}
	for _, cont := range c.reversed() {
		//		cont.runOrStart()
	}
}

// Start containers.
func (c Containers) start() {
	for _, cont := range c.reversed() {
		//		cont.start()
	}
}

// Kill containers.
func (c Containers) kill() {
	for _, cont := range c.data {
		//		cont.kill()
	}
}

// Stop containers.
func (c Containers) stop() {
	for _, cont := range c.data {
		//	cont.stop()
	}
}

// Remove containers.
// When forced, stops existing containers first.
func (c Containers) rm(force bool, kill bool) {
	if force {
		if kill {
			c.kill()
		} else {
			c.stop()
		}
	}
	for _, cont := range c.data {
		//		cont.rm()
	}
}

// Status of containers.
func (c Containers) status() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "Name\tRunning\tID\tIP\tPorts")
	for _, container := range c.data {
		//	container.status(w)
	}
	w.Flush()
}
