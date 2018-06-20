package table

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	idsHeader = [][]string{
		{"ID", "ID"},
	}
)

func (w *Writer) NewWriter(values [][]string, ctx *cli.Context) *Writer {
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

func (w *Writer) AddFormatFunc(name string, f FormatFunc) {
	w.funcMap[name] = f
}

func (w *Writer) Err() error {
	return w.err
}

func (w *Writer) writeHeader() {
	if w.HeaderFormat != "" && !w.headerPrinted {
		w.headerPrinted = true
		w.err = w.printTemplate(w.Writer, w.HeaderFormat, struct{}{})
		if w.err != nil {
			return
		}
	}
}

func (w *Writer) Write(obj interface{}) {
	if w.err != nil {
		return
	}

	w.writeHeader()
	if w.err != nil {
		return
	}

	if w.ValueFormat == "json" {
		content, err := json.Marshal(obj)
		w.err = err
		if w.err != nil {
			return
		}
		_, w.err = w.Writer.Write(append(content, byte('\n')))
	} else if w.ValueFormat == "yaml" {
		content, err := yaml.Marshal(obj)
		w.err = err
		if w.err != nil {
			return
		}
		w.Writer.Write([]byte("---\n"))
		_, w.err = w.Writer.Write(append(content, byte('\n')))
	} else {
		w.err = w.printTemplate(w.Writer, w.ValueFormat, obj)
	}
}

func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}
	w.writeHeader()
	if w.err != nil {
		return w.err
	}
	return w.Writer.Flush()
}

func (w *Writer) printTemplate(out io.Writer, templateContent string, obj interface{}) error {
	tmpl, err := template.New("").Funcs(w.funcMap).Parse(templateContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(out, obj)
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
