package cmd

import (
	"fmt"

	"github.com/cnrancher/cube-cli/util"
	"github.com/urfave/cli"
)

const (
	ServerDescription = `
Management RancherCUBE API-SERVER. 
					
Example:
	# Run the RancherCUBE api-server
	$ cube server run --port "9500" --kube-config /example/kube-config.yml
	# Stop the RancherCUBE api-server
	$ cube server stop
	# Remove the RancherCUBE api-server
	$ cube server rm
	# Get the RancherCUBE api-server status
	$ cube server status
`

	ServerPort     = "port"
	ConfigLocation = "kube-config"
)

func ServerCommand() cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"s"},
		Usage:       "Operations with cube api-server",
		Description: ServerDescription,
		Action:      defaultAction(serverStatus),
		Subcommands: []cli.Command{
			{
				Name:        "run",
				Usage:       "Run the RancherCUBE api-server",
				Description: "Run the RancherCUBE api-server",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  ServerPort,
						Usage: "Specify api-server listen port",
					},
					cli.StringFlag{
						Name:  ConfigLocation,
						Usage: "Specify api-server kubernetes config location",
					},
				},
				Action: defaultAction(serverRun),
			},
			{
				Name:        "stop",
				Usage:       "Stop the RancherCUBE api-server",
				Description: "Stop the RancherCUBE api-server",
				Action:      defaultAction(serverStop),
			},
			{
				Name:        "rm",
				Usage:       "Remove the RancherCUBE api-server",
				Description: "Remove the RancherCUBE api-server",
				Action:      defaultAction(serverRm),
			},
			{
				Name:        "status",
				Usage:       "Status the RancherCUBE api-server",
				Description: "Status the RancherCUBE api-server",
				Action:      defaultAction(serverStatus),
			},
		},
	}
}

func serverRun(ctx *cli.Context) error {
	port := ctx.String(ServerPort)
	configLocation := ctx.String(ConfigLocation)
	if "" == port {
		return fmt.Errorf("cube server run: require %v", ServerPort)
	}
	if "" == configLocation {
		return fmt.Errorf("cube server run: require %v", ConfigLocation)
	}

	// generate JWT RSA256 key file
	err := util.GenerateRSA256()
	if err != nil {
		return err
	}

	return nil
}

func serverStop(ctx *cli.Context) error {
	fmt.Println("stop server")
	return nil
}

func serverRm(ctx *cli.Context) error {
	fmt.Println("remove server")
	return nil
}

func serverStatus(ctx *cli.Context) error {
	fmt.Println("status server")
	return nil
}
