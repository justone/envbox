package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type RunCommand struct {
	Vars  []string `short:"e" long:"env" description:"Environment variables to expose" required:"yes"`
	Shell bool     `short:"s" long:"shell" description:"Wrap execution in a call to shell -c ''"`
}

var runCommand RunCommand

func (c *RunCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.RunCommandWithEnv(c.Vars, c.Shell, args)
}

func init() {
	cmd, err := parser.AddCommand("run", "Run a command.", "", &runCommand)

	cmd.Aliases = append(cmd.Aliases, "r")

	if err != nil {
		fmt.Println(err)
	}
}
