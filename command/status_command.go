package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//Displays the current status of all the containers, or the specified one(s).
func (c *Commander) NewStatusCommand() {
	c.Register(cli.Command{
		Name:  "status",
		Usage: "Displays status of containers",
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers engine.Containers) error {
				return containers.Status()
			}, ctx)
		},
	})
}
