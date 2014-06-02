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
		Action: func(c *cli.Context) {
			containersCommand(func(containers Containers) {
				containers.stop()
			}, c)
		},
	})
}
