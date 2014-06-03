package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//rm will call docker rm on all containers, or the specified one(s).
func (c *Commander) NewRemoveCommand() {
	c.Register(cli.Command{
		Name:  "rm",
		Usage: "Remove the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", false, "stop running containers first"},
			cli.BoolFlag{"kill, k", false, "when using --force, kill containers instead of stopping them"},
		},
		Action: func(c *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.remove(c.Bool("force"), c.Bool("kill"))
			}, c)
		},
	})
}
