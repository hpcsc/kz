package cmd

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
)

var Version = "main"

func Run() int {
	app := &cli.App{
		Name:    "kz",
		Usage:   "switch Kubernetes namespace and context using partial name",
		Version: Version,
		Action:  switchFromRoot,
		Commands: []*cli.Command{
			newNamespaceSubcommand(),
			newContextSubcommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
		return 1
	}

	return 0
}

func switchFromRoot(ctx *cli.Context) error {
	return switchContext(ctx)
}
