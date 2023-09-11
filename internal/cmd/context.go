package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/hpcsc/kz/internal/tui"
	"github.com/urfave/cli/v2"
)

func newContextSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ctx",
		Usage:   "commands to work with Kubernetes contexts",
		Aliases: []string{"context"},
		Action:  oneArgumentsAction(switchContext, "context name query is required"),
		Subcommands: []*cli.Command{
			{
				Name:   "sync",
				Usage:  "sync Kubernetes contexts from kube config files",
				Action: noArgumentsAction(syncContexts),
			},
			{
				Name:   "list",
				Usage:  "list available Kubernetes contexts",
				Action: noArgumentsAction(listContexts),
			},
		},
	}
}

func syncContexts() error {
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

func listContexts() error {
	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	for _, c := range cfg.Contexts {
		fmt.Println(c)
	}

	return nil
}

func switchContext(query string) error {
	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	destinationContexts := cfg.ContextsMatching(query)
	if len(destinationContexts) == 0 {
		return fmt.Errorf("no contexts matched query '%s'", query)
	}

	var contextToSwitch string
	if len(destinationContexts) == 1 {
		contextToSwitch = destinationContexts[0]
	} else {
		contextToSwitch, err = tui.ShowDropdown("Please select a context", destinationContexts)
		if err != nil {
			return err
		}
	}

	if err := kube.SwitchContextTo(contextToSwitch); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("switched to context %s", contextToSwitch))

	return nil
}
