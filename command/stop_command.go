package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//stop will call docker stop on all containers, or the specified one(s).
func (c *Commander) NewStopCommand() {
	c.Register(cli.Command{
		Name:  "stop",
		Usage: "Stop the containers",
		Action: func(ctx *cli.Context) {
			c.Execute(func(node engine.Node) error {
				return node.Stop()
			}, ctx)
		},
	})
}
