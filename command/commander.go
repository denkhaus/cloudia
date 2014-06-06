package command

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/engine"
	"github.com/denkhaus/tcgl/applog"
	"github.com/denkhaus/yamlconfig"
)

type Commander struct {
	engine *engine.Engine
	config *yamlconfig.Config
	app    *cli.App
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Execute(fn engine.EngineFunc, ctx *cli.Context) {
	path := ctx.String("manifest")
	group := ctx.String("group")

	err := c.engine.LoadFromFile(path, group)
	if err != nil {
		applog.Errorf("manifest error:: %s", err.Error())
		return
	}

	if err = c.engine.Execute(fn); err != nil {
		applog.Errorf("execution error:: %s", err.Error())
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func NewCommander(app *cli.App, cnf *yamlconfig.Config) (*Commander, error) {
	cmd := &Commander{app: app, config: cnf}

	storPrefix := cnf.GetString("storage:prefix")
	storAddress := cnf.GetString("storage:address")
	storPassword := cnf.GetString("storage:password")

	if engine, err := engine.NewEngine(
		storeAddress, storPassword, storPrefix); err != nil {
		cmd.engine = engine
	} else {
		return nil, err
	}

	cmd.NewLiftCommand()
	cmd.NewProvisionCommand()
	cmd.NewRunCommand()
	cmd.NewStartCommand()
	cmd.NewStopCommand()
	cmd.NewRemoveCommand()
	cmd.NewKillCommand()
	cmd.NewStatusCommand()

	return cmd, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Register(cmd cli.Command) {
	c.app.Commands = append(c.app.Commands, cmd)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *Commander) Run(args []string) {
	c.app.Run(args)
}
