package cmd

import "github.com/urfave/cli/v2"

func newContextSubcommand() *cli.Command {
	return &cli.Command{
		Name:        "ctx",
		Aliases:     []string{"context"},
		Subcommands: []*cli.Command{},
	}
}
