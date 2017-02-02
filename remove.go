package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type RemoveCommand struct {
	Name string `short:"n" long:"name" description:"Name of environment variable." required:"yes"`
}

var removeCommand RemoveCommand

func (c *RemoveCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.RemoveVariable(c.Name)
}

func init() {
	cmd, err := parser.AddCommand("remove", "Remove an environment variable.", "", &removeCommand)

	cmd.Aliases = append(cmd.Aliases, "rm")

	if err != nil {
		fmt.Println(err)
	}
}
