package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/urfave/cli/v2"
	"strings"
)

func newNamespaceSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ns",
		Aliases: []string{"namespace"},
		Subcommands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "track a Kubernetes namespace",
				Action: addNamespaces,
			},
		},
	}
}

func addNamespaces(ctx *cli.Context) error {
	toBeAdded := ctx.Args().Slice()
	if len(toBeAdded) == 0 {
		return fmt.Errorf("no namespaces provided")
	}

	c, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	c.AddNamespaces(toBeAdded...)

	if err := config.SaveToDefaultLocation(c); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("namespace(s) %s added", strings.Join(toBeAdded, ", ")))

	return nil
}
