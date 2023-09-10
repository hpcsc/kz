package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/hpcsc/kz/internal/tui"
	"github.com/urfave/cli/v2"
	"strings"
)

func newNamespaceSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ns",
		Usage:   "commands to work with Kubernetes namespaces",
		Aliases: []string{"namespace"},
		Action:  switchNamespace,
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

	if len(c.Namespaces) == 0 {
		fmt.Println("no namespaces available")
		return nil
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

func switchNamespace(ctx *cli.Context) error {
	query := ctx.Args().First()
	if len(query) == 0 {
		return fmt.Errorf("namespace name query is required")
	}

	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	destinationNamespaces := cfg.NamespacesMatching(query)
	var namespaceToSwitch string
	if len(destinationNamespaces) == 0 {
		namespaceToSwitch = query
	} else if len(destinationNamespaces) == 1 {
		namespaceToSwitch = destinationNamespaces[0]
	} else {
		namespaceToSwitch, err = tui.ShowDropdown("Please select a namespace", destinationNamespaces)
		if err != nil {
			return err
		}
	}

	if err := kube.SwitchNamespaceTo(namespaceToSwitch); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("switched to namespace %s", namespaceToSwitch))
	return nil
}
