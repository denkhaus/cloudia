package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//kill will call docker kill on all containers, or the specified one(s).
func (c *Commander) NewKillCommand() {
	c.Register(cli.Command{
		Name:  "kill",
		Usage: "Kill the containers",
		Action: func(ctx *cli.Context) {
			c.Execute(func(node engine.Node) error {
				return node.Kill()
			}, ctx)
		},
	})
}
