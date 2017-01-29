package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type AddCommand struct {
	Name string `short:"n" long:"name" description:"Name of environment variable." required:"yes"`
}

var addCommand AddCommand

func (c *AddCommand) Execute(args []string) error {
	// TODO: check for duplicate name

	key, err := ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	value, err := PromptForValue()
	if err != nil {
		return errors.Wrap(err, "error reading value")
	}

	return AddVariable(key, c.Name, value)
}

func init() {
	cmd, err := parser.AddCommand("add", "Add an environment variable.", "", &addCommand)

	cmd.Aliases = append(cmd.Aliases, "a")

	if err != nil {
		fmt.Println(err)
	}
}
