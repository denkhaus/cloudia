package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//lift will use specified Dockerfiles to build all the containers, or the specified one(s).
//If no Dockerfile is given, it will pull the image(s) from the given registry.
func (c *Commander) NewLiftCommand() {
	c.Register(cli.Command{
		Name:  "lift",
		Usage: "Build or pull images, then run or start the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", "rebuild all images"},
			cli.BoolFlag{"kill, k", "kill containers"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers engine.Containers) error {
				return containers.Lift(ctx.Bool("force"), ctx.Bool("kill"))
			}, ctx)
		},
	})
}
