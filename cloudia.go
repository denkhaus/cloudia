package main

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/command"
	"github.com/denkhaus/tcgl/applog"
	"github.com/denkhaus/yamlconfig"
	"os"
)

var releaseVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "cldia"
	app.Version = releaseVersion
	app.Usage = "A command line client for Cloudia - easy clusterwise docker orchestration."
	app.Flags = []cli.Flag{
		cli.StringFlag{"group, g", "", "group or container to restrict the command to"},
		cli.StringFlag{"manifest, m", "", "path to a manifest (.json, .yml, .yaml) file to read from"},
		cli.BoolFlag{"debug, d", "print debug output"},
		//	cli.StringSliceFlag{"peers, C", &cli.StringSlice{}, "a comma-delimited list of machine addresses in the cluster (default: {\"127.0.0.1:4001\"})"},
	}

	cnf := yamlconfig.NewConfig(".cldiarc")
	if err := cnf.Load(func(config *yamlconfig.Config) {
		config.SetDefault("storage:address", "127.0.0.1:6379")
		config.SetDefault("storage:password", "")
		config.SetDefault("storage:prefix", "cldia")
	}, "", false); err != nil {
		applog.Errorf("config error:: load config %s", err.Error())
		return
	}

	cmdr, err := command.NewCommander(app, cnf)
	if err != nil {
		applog.Errorf("startup error:: %s", err.Error())
		return
	}
	cmdr.Run(os.Args)
}
