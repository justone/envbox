package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

type GlobalOptions struct {
	Quiet   func() `short:"q" long:"quiet" description:"Show as little information as possible."`
	Verbose func() `short:"v" long:"verbose" description:"Show verbose debug information."`
}

var globalOptions GlobalOptions
var parser = flags.NewParser(&globalOptions, flags.Default)

func main() {
	// configure logging
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// options to change log level
	globalOptions.Quiet = func() {
		logrus.SetLevel(logrus.WarnLevel)
	}
	globalOptions.Verbose = func() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
