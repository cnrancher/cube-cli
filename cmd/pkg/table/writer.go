package table

import (
	"github.com/urfave/cli"
)

type WriterInterface interface {
	NewWriter(values [][]string, ctx *cli.Context) *Writer
}

type FormatFunc interface{}
