package cmd

import (
	"github.com/urfave/cli/v2"
)

func newNamespaceSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ns",
		Aliases: []string{"namespace"},
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "track a Kubernetes namespace",
				Action: func(*cli.Context) error {
					return nil
				},
			},
		},
	}
}
