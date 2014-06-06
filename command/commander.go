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
	filePath := ctx.GlobalString("manifest")
	group := ctx.GlobalString("group")

	if len(filePath) > 0 {
		err := c.engine.LoadFromFile(filePath, group)
		if err != nil {
			applog.Errorf(err.Error())
			return
		}
	} else {
		err := c.engine.LoadDefaults(group)
		if err != nil {
			applog.Errorf(err.Error())
			return
		}
	}

	if err := c.engine.Execute(fn); err != nil {
		applog.Errorf("Execution error:: %s", err.Error())
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
		storAddress, storPassword, storPrefix); err != nil {
		return nil, err
	} else {
		cmd.engine = engine
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
