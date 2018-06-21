package table

import (
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

func NewNodeWriter(values [][]string, ctx *cli.Context) *Writer {
	t := &Writer{
		Writer: tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', tabwriter.TabIndent),
		funcMap: map[string]interface{}{
			"json": FormatJSON,
			"yaml": FormatYAML,
		},
	}
	t.HeaderFormat, t.ValueFormat = SimpleFormat(values)

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
