package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/hpcsc/kz/internal/tui"
	"github.com/urfave/cli/v2"
	"os"
)

var Version = "main"

func Run() int {
	app := &cli.App{
		Name:                 "kz",
		Usage:                "switch Kubernetes namespace and context using partial name",
		Version:              Version,
		EnableBashCompletion: true,
		Action:               switchFromRoot,
		Commands: []*cli.Command{
			newNamespaceSubcommand(),
			newContextSubcommand(),
			newUpdateSubcommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
		return 1
	}

	return 0
}

func switchFromRoot(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		query := ctx.Args().First()
		return switchContext(query)
	}

	if ctx.Args().Len() > 1 {
		contextQuery := ctx.Args().Get(0)
		namespaceQuery := ctx.Args().Get(1)

		if contextQuery == "-" {
			return switchNamespace(namespaceQuery)
		}

		return switchContextAndNamespace(contextQuery, namespaceQuery)
	}

	return fmt.Errorf("context name query is required")
}

func switchContextAndNamespace(contextQuery string, namespaceQuery string) error {
	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	destinationContexts := cfg.ContextsMatching(contextQuery)
	if len(destinationContexts) == 0 {
		return fmt.Errorf("no contexts matched query '%s'", contextQuery)
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

	destinationNamespaces := cfg.NamespacesMatching(namespaceQuery)
	var namespaceToSwitch string
	if len(destinationNamespaces) == 0 {
		namespaceToSwitch = namespaceQuery
	} else if len(destinationNamespaces) == 1 {
		namespaceToSwitch = destinationNamespaces[0]
	} else {
		namespaceToSwitch, err = tui.ShowDropdown("Please select a namespace", destinationNamespaces)
		if err != nil {
			return err
		}
	}

	if err := kube.SwitchContextAndNamespace(contextToSwitch, namespaceToSwitch); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("switched to context %s, namespace %s", contextToSwitch, namespaceToSwitch))

	return nil
}
