package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//run will call docker run on all containers, or the specified one(s).
func (c *Commander) NewRunCommand() {
	c.Register(cli.Command{
		Name:  "run",
		Usage: "Run the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", "stop and remove running containers first"},
			cli.BoolFlag{"kill, k", "when using --force, kill containers instead of stopping them"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(node engine.Node) error {
				return node.Run(ctx.Bool("force"), ctx.Bool("kill"))
			}, ctx)
		},
	})
}
