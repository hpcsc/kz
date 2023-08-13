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
			{
				Name:   "list",
				Usage:  "list tracked Kubernetes namespaces",
				Action: listNamespaces,
			},
			{
				Name:   "delete",
				Usage:  "delete tracked Kubernetes namespaces",
				Action: deleteNamespaces,
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

func listNamespaces(ctx *cli.Context) error {
	c, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	for _, n := range c.Namespaces {
		fmt.Println(n)
	}

	return nil
}

func deleteNamespaces(ctx *cli.Context) error {
	toBeDeleted := ctx.Args().Slice()
	if len(toBeDeleted) == 0 {
		return fmt.Errorf("no namespaces provided")
	}

	c, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	c.DeleteNamespaces(toBeDeleted...)

	if err := config.SaveToDefaultLocation(c); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("namespace(s) %s deleted", strings.Join(toBeDeleted, ", ")))

	return nil
}
