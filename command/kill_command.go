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
			c.Execute(func(containers engine.Containers) error {
				return containers.Kill()
			}, ctx)
		},
	})
}
