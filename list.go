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

	for name, envVar := range vars {
		// TODO: figure out a better way to list these
		fmt.Print(name)
		if envVar.Exposed != envVar.Name {
			fmt.Printf("(%s)", envVar.Exposed)
		}
		fmt.Printf("=%s", envVar.Value)
		fmt.Println()
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
