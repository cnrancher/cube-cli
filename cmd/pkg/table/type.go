package table

import (
	"text/tabwriter"
)

type Writer struct {
	quite         bool
	HeaderFormat  string
	ValueFormat   string
	err           error
	headerPrinted bool
	Writer        *tabwriter.Writer
	funcMap       map[string]interface{}
}
