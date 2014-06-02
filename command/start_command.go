package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//start will call docker start on all containers, or the specified one(s).
func (c *Commander) NewStartCommand() {
	c.Register(cli.Command{
		Name:  "start",
		Usage: "Start the containers",
		Action: func(c *cli.Context) {
			containersCommand(func(containers Containers) {
				containers.start()
			}, c)
		},
	})
}
