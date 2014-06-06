package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

//rm will call docker rm on all containers, or the specified one(s).
func (c *Commander) NewRemoveCommand() {
	c.Register(cli.Command{
		Name:  "rm",
		Usage: "Remove the containers",
		Flags: []cli.Flag{
			cli.BoolFlag{"force, f", "stop running containers first"},
			cli.BoolFlag{"kill, k", "when using --force, kill containers instead of stopping them"},
		},
		Action: func(ctx *cli.Context) {
			c.Execute(func(node engine.Node) error {
				return node.Remove(ctx.Bool("force"), ctx.Bool("kill"))
			}, ctx)
		},
	})
}
