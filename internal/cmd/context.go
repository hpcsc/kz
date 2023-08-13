package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/urfave/cli/v2"
)

func newContextSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ctx",
		Aliases: []string{"context"},
		Subcommands: []*cli.Command{
			{
				Name:   "sync",
				Usage:  "sync Kubernetes contexts from kube config files",
				Action: syncContexts,
			},
		},
	}
}

func syncContexts(ctx *cli.Context) error {
	c, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	contexts, err := kube.ContextsFromConfig()
	if err != nil {
		return err
	}

	c.Contexts = contexts

	if err := config.SaveToDefaultLocation(c); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("%d contexts synced", len(contexts)))

	return nil
}
