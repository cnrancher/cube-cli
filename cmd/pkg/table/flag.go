package table

import "github.com/urfave/cli"

var outputFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "quiet,q",
		Usage: "Only display IDs",
	},
	cli.StringFlag{
		Name:  "format",
		Usage: "'json' or 'yaml'",
	},
	cli.BoolFlag{
		Name:  "ids",
		Usage: "Include ID column in output",
	},
}

func WriterFlags() []cli.Flag {
	return outputFlags
}
