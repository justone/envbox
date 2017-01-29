package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type GenerateKeyCommand struct {
	Set bool `short:"s" long:"set" description:"Set the new key as the one to use on this system."`
}

type SetKeyCommand struct{}

type KeyCommand struct {
	Generate GenerateKeyCommand `command:"generate" alias:"gen" description:"Generate new key."`
	Set      SetKeyCommand      `command:"set" description:"Set key."`
}

func (r *GenerateKeyCommand) Execute(args []string) error {
	key, err := GenerateNewKey()
	if err != nil {
		return errors.Wrap(err, "unable to generate key")
	}
	fmt.Println(key)

	if r.Set {
		// TODO: warn when overriding existing key
		return StoreKey(key)
	}
	return nil
}

func (r *SetKeyCommand) Execute(args []string) error {
	key, err := PromptForKey()
	if err != nil {
		return errors.Wrap(err, "unable to prompt for key")
	}

	return StoreKey(key)
}

func init() {
	var keyCommand KeyCommand

	cmd, err := parser.AddCommand("key", "Manage key.", "", &keyCommand)

	cmd.Aliases = append(cmd.Aliases, "k")

	if err != nil {
		fmt.Println(err)
	}
}
