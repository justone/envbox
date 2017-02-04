package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type GenerateKeyCommand struct {
	Set bool `short:"s" long:"set" description:"Set the new key as the one to use on this system."`
}

type SetKeyCommand struct{}

type ShowKeyCommand struct{}

type KeyCommand struct {
	Generate GenerateKeyCommand `command:"generate" alias:"gen" description:"Generate new key."`
	Set      SetKeyCommand      `command:"set" description:"Set key."`
	Show     ShowKeyCommand     `command:"show" description:"Show key."`
}

func (r *GenerateKeyCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.GenerateNewKey(r.Set)
}

func (r *SetKeyCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.PromptAndStoreKey()
}

func (r *ShowKeyCommand) Execute(args []string) error {
	box, err := NewEnvBox()
	if err != nil {
		return errors.Wrap(err, "unable to create env box")
	}

	return box.ShowKey()
}

func init() {
	var keyCommand KeyCommand

	cmd, err := parser.AddCommand("key", "Manage key.", "", &keyCommand)

	cmd.Aliases = append(cmd.Aliases, "k")

	if err != nil {
		fmt.Println(err)
	}
}
