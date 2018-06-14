package cmd

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/urfave/cli"
)

func PromptCommand() cli.Command {
	return cli.Command{
		Name:      "prompt",
		Usage:     "Enter rancher cli auto-prompt mode",
		ArgsUsage: "None",
		Action:    promptAction,
		Flags:     []cli.Flag{},
	}
}

func promptAction(ctx *cli.Context) error {
	fmt.Print("cube cli auto-completion mode")
	defer fmt.Println("Goodbye!")
	p := prompt.New(
		Executor,
		Completer,
		prompt.OptionTitle("cube-prompt: interactive cube cli"),
		prompt.OptionPrefix("cube "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionMaxSuggestion(20),
	)
	p.Run()
	return nil
}
