package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/cnrancher/cube-cli/cmd"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "cube"
	app.Version = VERSION
	app.Usage = "RancherCUBE CLI, managing cube services"
	app.Author = "Rancher Labs, Inc."
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "enable debug logging level",
			EnvVar: "RANCHER_DEBUG",
		},
	}

	app.Commands = []cli.Command{
		cmd.ServerCommand(),
		cmd.NodeCommand(),
		cmd.RKECommand(),
		cmd.PromptCommand(),
	}

	cmd.Flags = app.Flags

	parsed, err := parseArgs(os.Args)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	app.Run(parsed)
}

var singleAlphaLetterRegxp = regexp.MustCompile("[a-zA-Z]")

func parseArgs(args []string) ([]string, error) {
	result := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 1 {
			for i, c := range arg[1:] {
				if string(c) == "=" {
					if i < 1 {
						return nil, errors.New("invalid input with '-' and '=' flag")
					}
					result[len(result)-1] = result[len(result)-1] + arg[i+1:]
					break
				} else if singleAlphaLetterRegxp.MatchString(string(c)) {
					result = append(result, "-"+string(c))
				} else {
					return nil, errors.Errorf("invalid input %v in flag", string(c))
				}
			}
		} else {
			result = append(result, arg)
		}
	}
	return result, nil
}
