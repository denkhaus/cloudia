package main

import (
	"github.com/codegangsta/cli"
	"github.com/denkhaus/cloudia/command"
	"os"
)

var releaseVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "cldia"
	app.Version = releaseVersion
	app.Usage = "A command line client for Cloudia - a docker orchestration framework."
	app.Flags = []cli.Flag{
		cli.BoolFlag{"verbose, v", "displays more verbose output"},
		cli.StringFlag{"group, g", "", "group or container to restrict the command to"},
		cli.StringFlag{"manifest, m", "", "path to a cloudia.(json,yml,yaml) file to read from"},

		//	cli.StringSliceFlag{"peers, C", &cli.StringSlice{}, "a comma-delimited list of machine addresses in the cluster (default: {\"127.0.0.1:4001\"})"},
	}

	cmdr, err := command.NewCommander(app)
	if err != nil {
		//TODO handle error
	}
	cmdr.Run(os.Args)
}
