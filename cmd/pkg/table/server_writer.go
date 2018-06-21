package table

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/urfave/cli"
)

var (
	idsHeader = [][]string{
		{"ID", "ID"},
	}
)

func NewServerWriter(values [][]string, ctx *cli.Context) *Writer {
	if ctx.Bool("ids") {
		values = append(idsHeader, values...)
	}

	t := &Writer{
		Writer: tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', tabwriter.TabIndent),
		funcMap: map[string]interface{}{
			"id":   FormatContainerID,
			"cmd":  FormatContainerCommand,
			"port": FormatContainerPort,
			"ago":  FormatContainerCreated,
			"name": FormatContainerName,
			"json": FormatJSON,
			"yaml": FormatYAML,
		},
	}
	t.HeaderFormat, t.ValueFormat = SimpleFormat(values)

	if ctx.Bool("quiet") {
		t.HeaderFormat = ""
		t.ValueFormat = "{{.ID}}\n"
	}

	customFormat := ctx.String("format")
	if customFormat == "json" {
		t.HeaderFormat = ""
		t.ValueFormat = "json"
	} else if customFormat == "yaml" {
		t.HeaderFormat = ""
		t.ValueFormat = "yaml"
	} else if customFormat != "" {
		t.ValueFormat = customFormat + "\n"
		t.HeaderFormat = ""
	}

	return t
}

func FormatContainerID(data interface{}) (string, error) {
	containerID, ok := data.(string)
	if !ok {
		return "", nil
	}

	if len(containerID) > 12 {
		return containerID[:12], nil
	}

	return containerID, nil
}

func FormatContainerCommand(data interface{}) (string, error) {
	command, ok := data.(string)
	if !ok {
		return "", nil
	}

	if len(command) > 25 {
		return command[:25], nil
	}

	return command, nil
}

func FormatContainerName(data []string) (string, error) {
	names := ""
	for index, name := range data {
		if index == 0 {
			names += name[1:]
		} else {
			names += " " + name[1:]
		}
	}

	return names, nil
}

func FormatContainerPort(data []types.Port) (string, error) {
	ports := ""
	for index, port := range data {
		if index == 0 {
			ports += port.IP + ":" + strconv.Itoa(int(port.PublicPort)) + "->" + strconv.Itoa(int(port.PrivatePort)) + "/" + port.Type
		} else {
			ports += " " + port.IP + ":" + strconv.Itoa(int(port.PublicPort)) + "->" + strconv.Itoa(int(port.PrivatePort)) + "/" + port.Type
		}
	}

	return fmt.Sprintf("%s", ports), nil
}

func FormatContainerCreated(data interface{}) (string, error) {
	s, ok := data.(int64)
	if !ok {
		return "", nil
	}

	t, err := time.Parse(time.RFC3339, time.Unix(s, 0).Format(time.RFC3339))
	if err != nil {
		return "", err
	}

	return units.HumanDuration(time.Now().UTC().Sub(t)) + " ago", nil
}
