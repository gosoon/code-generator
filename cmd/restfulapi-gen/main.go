package main

import (
	"flag"

	generatorargs "github.com/gosoon/code-generator/cmd/args"
	"github.com/gosoon/code-generator/cmd/generators"

	"github.com/gosoon/glog"
	"github.com/spf13/pflag"
)

func main() {
	genericArgs, customArgs := generatorargs.NewDefaults()

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := generatorargs.Validate(genericArgs); err != nil {
		glog.Fatalf("Error: %v", err)
	}

	// Run it.
	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		glog.Fatalf("Error: %v", err)
	}
	glog.Info("Completed successfully.")
}
