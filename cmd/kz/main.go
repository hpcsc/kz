package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
)

var Version = "main"

func main() {
	app := &cli.App{
		Name:     "kz",
		Usage:    "switch Kubernetes namespace and context using partial name",
		Version:  Version,
		Commands: []*cli.Command{},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
