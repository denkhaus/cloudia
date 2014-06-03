package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
)

type Commander struct {
	engine engine.Engine
	app    *cli.App
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Execute(fn EngineFunc, cont *cli.Context) {
	path := c.String("manifest")
	group := c.String("group")

	err := c.engine.LoadFromFile(path, group)
	if err != nil {
		//TODO Handle Error
	}

	err := c.engine.Execute(fn)
	if err != nil {
		//TODO Handle Error
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func NewCommander(app *cli.App) *Commander {
	cmd = &Commander{app: app, engine: engine.NewEngine()}

	cmd.NewLiftCommand()
	cmd.NewProvisionCommand()
	cmd.NewRunCommand()
	cmd.NewStartCommand()
	cmd.NewStopCommand()
	cmd.NewRemoveCommand()
	cmd.NewKillCommand()
	cmd.NewStatusCommand()

	return cmd
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Register(cmd *cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Run(args []string) {
	c.app.Run(args)
}
