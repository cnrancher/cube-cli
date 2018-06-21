package table

import "github.com/urfave/cli"

var outputServerFlags = []cli.Flag{
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

var outputNodeFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "format",
		Usage: "'json' or 'yaml'",
	},
	cli.BoolFlag{
		Name:  "ids",
		Usage: "Include ID column in output",
	},
}

func WriterServerFlags() []cli.Flag {
	return outputServerFlags
}

func WriterNodeFlags() []cli.Flag {
	return outputNodeFlags
}
