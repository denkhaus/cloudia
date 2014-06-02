package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//kill will call docker kill on all containers, or the specified one(s).
func (c *Commander) NewKillCommand() {
	c.Register(cli.Command{
		Name:  "kill",
		Usage: "Kill the containers",
		Action: func(c *cli.Context) {
			containersCommand(func(containers Containers) {
				containers.kill()
			}, c)
		},
	})
}
