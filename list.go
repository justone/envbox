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
	key, err := ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	vars, err := LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load vars")
	}

	for key, value := range vars {
		fmt.Printf("%s=%s\n", key, value)
	}

	return nil
}

func init() {
	cmd, err := parser.AddCommand("list", "List environment variables.", "", &listCommand)

	cmd.Aliases = append(cmd.Aliases, "ls")

	if err != nil {
		fmt.Println(err)
	}
}
