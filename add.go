package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type AddCommand struct {
	Name     string `short:"n" long:"name" description:"Name of environment variable." required:"yes"`
	File     string `short:"f" long:"file" description:"File with contents of variable"`
	Exposed  string `short:"e" long:"exposed" description:"Name of exposed variable, if different than the name."`
	Multiple bool   `short:"m" long:"multiple" description:"Add multiple variables after the first."`
}

var addCommand AddCommand

func (c *AddCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.AddVariable(c.Name, c.Exposed, c.File, c.Multiple)
}

func init() {
	cmd, err := parser.AddCommand("add", "Add an environment variable.", "", &addCommand)

	cmd.Aliases = append(cmd.Aliases, "a")

	if err != nil {
		fmt.Println(err)
	}
}
