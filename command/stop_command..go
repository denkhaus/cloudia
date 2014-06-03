package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//stop will call docker stop on all containers, or the specified one(s).
func (c *Commander) NewStartCommand() {
	c.Register(cli.Command{
		Name:  "stop",
		Usage: "Stop the containers",
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.stop()
			}, ctx)
		},
	})
}
