package table

import (
	"encoding/json"
	"io"
	"text/template"

	"gopkg.in/yaml.v2"
)

type FormatFunc interface{}

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
