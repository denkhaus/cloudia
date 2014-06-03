package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//lift will use specified Dockerfiles to build all the containers, or the specified one(s).
//If no Dockerfile is given, it will pull the image(s) from the given registry.
func (c *Commander) NewLiftCommand() {
	c.Register(cli.Command{
		Name:  "lift",
		Usage: "Build or pull images, then run or start the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", false, "rebuild all images"},
			cli.BoolFlag{"kill, k", false, "kill containers"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.lift(c.Bool("force"), c.Bool("kill"))
			}, ctx)
		},
	})
}
