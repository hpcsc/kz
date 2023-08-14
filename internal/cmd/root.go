package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/config"
	"github.com/hpcsc/kz/internal/kube"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/tools/clientcmd"
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
	if ctx.Args().Len() == 1 {
		return switchContext(ctx)
	} else if ctx.Args().Len() > 1 {
		return switchContextAndNamespace(ctx)
	}

	return fmt.Errorf("context name query is required")
}

func switchContextAndNamespace(ctx *cli.Context) error {
	contextQuery := ctx.Args().Get(0)
	namespaceQuery := ctx.Args().Get(1)

	cfg, err := config.LoadFromDefaultLocation()
	if err != nil {
		return err
	}

	destinationContexts := cfg.ContextsMatching(contextQuery)
	if len(destinationContexts) == 0 {
		return fmt.Errorf("no contexts matched query '%s'", contextQuery)
	}

	destinationNamespaces := cfg.NamespacesMatching(namespaceQuery)
	var namespaceToSwitch string
	if len(destinationNamespaces) == 0 {
		namespaceToSwitch = namespaceQuery
	} else {
		namespaceToSwitch = destinationNamespaces[0]
	}

	// switched to 1st matching context and namespace for now
	if err := kube.SwitchContextAndNamespace(destinationContexts[0], namespaceToSwitch, clientcmd.RecommendedHomeFile); err != nil {
		return err
	}

	color.Green(fmt.Sprintf("switched to context %s, namespace %s", destinationContexts[0], namespaceToSwitch))

	return nil
}
