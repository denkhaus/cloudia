package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//start will call docker start on all containers, or the specified one(s).
func (c *Commander) NewStartCommand() {
	c.Register(cli.Command{
		Name:  "start",
		Usage: "Start the containers",
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers engine.Containers) error {
				return containers.Start()
			}, ctx)
		},
	})
}
