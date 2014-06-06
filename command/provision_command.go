package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//provision will use specified Dockerfiles to build all the containers, or the specified one(s).
//If no Dockerfile is given, it will pull the image(s) from the given registry.
func (c *Commander) NewProvisionCommand() {
	c.Register(cli.Command{
		Name:  "provision",
		Usage: "Build or pull images",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", "rebuild all images"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(node engine.Node) error {
				return node.Provision(ctx.Bool("force"))
			}, ctx)
		},
	})
}
