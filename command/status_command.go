package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//Displays the current status of all the containers, or the specified one(s).
func (c *Commander) NewStatusCommand() {
	c.Register(cli.Command{
		Name:  "status",
		Usage: "Displays status of containers",
		Action: func(c *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.status()
			}, c)
		},
	})
}
