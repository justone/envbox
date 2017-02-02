package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type ListCommand struct {
	// nothing
}

var listCommand ListCommand

func (c *ListCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.ListVariables()
}

func init() {
	cmd, err := parser.AddCommand("list", "List environment variables.", "", &listCommand)

	cmd.Aliases = append(cmd.Aliases, "ls")

	if err != nil {
		fmt.Println(err)
	}
}
