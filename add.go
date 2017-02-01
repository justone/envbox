package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type AddCommand struct {
	Name    string `short:"n" long:"name" description:"Name of environment variable." required:"yes"`
	File    string `short:"f" long:"file" description:"File with contents of variable"`
	Exposed string `short:"e" long:"exposed" description:"Name of exposed variable, if different than the name."`
}

var addCommand AddCommand

func (c *AddCommand) Execute(args []string) error {
	var err error

	key, err := ReadKey()
	if err != nil {
		return errors.Wrap(err, "unable to read key")
	}

	// check for duplicate name
	vars, err := LoadEnvVars(key)
	if err != nil {
		return errors.Wrap(err, "unable to load vars")
	}

	if _, ok := vars[c.Name]; ok {
		return fmt.Errorf("var %s already exists", c.Name)
	}

	exposed := c.Exposed
	if len(c.Exposed) == 0 {
		exposed = c.Name
	}

	var value string
	if len(c.File) > 0 {
		data, err := ioutil.ReadFile(c.File)
		if err != nil {
			return errors.Wrap(err, "error reading file")
		}
		value = strings.TrimSpace(string(data))
	} else {
		value, err = PromptForValue()
		if err != nil {
			return errors.Wrap(err, "error reading value")
		}
	}

	return AddVariable(key, c.Name, exposed, value)
}

func init() {
	cmd, err := parser.AddCommand("add", "Add an environment variable.", "", &addCommand)

	cmd.Aliases = append(cmd.Aliases, "a")

	if err != nil {
		fmt.Println(err)
	}
}
