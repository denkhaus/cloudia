package command

import (
	"errors"
	"github.com/codegangsta/cli"
)

//provision will use specified Dockerfiles to build all the containers, or the specified one(s).
//If no Dockerfile is given, it will pull the image(s) from the given registry.
func (c *Commander) NewProvisionCommand() {
	c.Register(cli.Command{
		Name:  "provision",
		Usage: "Build or pull images",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", false, "rebuild all images"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(containers Containers) {
				return containers.provision(c.Bool("force"))
			}, ctx)
		},
	})
}
