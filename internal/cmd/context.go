package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/hpcsc/kz/internal/tui"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/tools/clientcmd"
)

func newContextSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "ctx",
		Aliases: []string{"context"},
		Action:  switchContext,
		Subcommands: []*cli.Command{
			{
				Name:   "sync",
				Usage:  "sync Kubernetes contexts from kube config files",
				Action: syncContexts,
			},
			{
				Name:   "list",
				Usage:  "list available Kubernetes contexts",
				Action: listContexts,
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

func listContexts(ctx *cli.Context) error {
	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	for _, c := range cfg.Contexts {
		fmt.Println(c)
	}

	return nil
}

func switchContext(ctx *cli.Context) error {
	query := ctx.Args().First()
	if len(query) == 0 {
		return fmt.Errorf("context name query is required")
	}

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

	if err := kube.SwitchContextTo(contextToSwitch, clientcmd.RecommendedHomeFile); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("switched to context %s", contextToSwitch))

	return nil
}
