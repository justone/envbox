package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type ShowCommand struct {
	Name   string `short:"n" long:"name" description:"Name of environment variable." required:"yes"`
	Export bool   `short:"e" long:"export" description:"Instead of human readable, format for shell eval"`
}

var showCommand ShowCommand

func (c *ShowCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	if c.Export {
		return box.ExportVariable(c.Name)
	}
	return box.ShowVariable(c.Name)
}

func init() {
	_, err := parser.AddCommand("show", "Show an environment variable.", "", &showCommand)

	if err != nil {
		fmt.Println(err)
	}
}
