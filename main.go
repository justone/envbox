package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type GlobalOptions struct {
	// nothing yet
}

var globalOptions GlobalOptions
var parser = flags.NewParser(&globalOptions, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
