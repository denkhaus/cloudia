package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//run will call docker run on all containers, or the specified one(s).
func (c *Commander) NewRunCommand() {
	c.Register(cli.Command{
		Name:  "run",
		Usage: "Run the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", false, "stop and remove running containers first"},
			cli.BoolFlag{"kill, k", false, "when using --force, kill containers instead of stopping them"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.run(c.Bool("force"), c.Bool("kill"))
			}, ctx)
		},
	})
}
