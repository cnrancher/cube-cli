package cmd

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/urfave/cli"
)

var (
	Commands = map[string]cli.Command{
		"server": ServerCommand(),
		"node":   NodeCommand(),
	}
	Flags = []cli.Flag{}
)

func Completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}

	args := strings.Split(d.TextBeforeCursor(), " ")
	w := d.GetWordBeforeCursor()

	// If PIPE is in text before the cursor, returns empty suggestions.
	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}

	// If word before the cursor starts with "-", returns CLI flag options.
	if strings.HasPrefix(w, "-") {
		return optionCompleter(args, strings.HasPrefix(w, "--"))
	}

	return argumentsCompleter(excludeOptions(args))
}

func argumentsCompleter(args []string) []prompt.Suggest {
	suggests := []prompt.Suggest{}
	for name, command := range Commands {
		if command.Name != "prompt" {
			suggests = append(suggests, prompt.Suggest{
				Text:        name,
				Description: command.Usage,
			})
		}
	}

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(suggests, args[0], true)
	}

	switch args[0] {
	case "server":
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "run", Description: "Run RancherCUBE api-server"},
				{Text: "status", Description: "Status the RancherCUBE api-server"},
				{Text: "stop", Description: "Stop the RancherCUBE api-server"},
				{Text: "rm", Description: "Remove the RancherCUBE api-server"},
			}
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}
	case "node":
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "ls", Description: "List the Rancher Kubernetes Engine Nodes"},
				{Text: "add", Description: "Add the Rancher Kubernetes Engine Node"},
				{Text: "rm", Description: "Remove the Rancher Kubernetes Engine Node"},
			}
			return prompt.FilterHasPrefix(subcommands, args[1], true)
		}
	default:
		if len(args) == 2 {
			return prompt.FilterHasPrefix(getSubcommandSuggest(args[0]), args[1], true)
		}
	}
	return []prompt.Suggest{}
}

func getSubcommandSuggest(name string) []prompt.Suggest {
	subcommands := []prompt.Suggest{}
	for _, com := range Commands[name].Subcommands {
		subcommands = append(subcommands, prompt.Suggest{
			Text:        com.Name,
			Description: com.Usage,
		})
	}
	return subcommands
}
