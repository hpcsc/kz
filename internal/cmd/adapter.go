package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// an adapter function that adapt a no-arguments function to CLI action handler
func noArgumentsAction(f func() error) func(ctx *cli.Context) error {
	return func(_ *cli.Context) error {
		return f()
	}
}

func oneArgumentsAction(f func(string) error, validationMsg string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		query := ctx.Args().First()
		if len(query) == 0 {
			return fmt.Errorf(validationMsg)
		}

		return f(query)
	}
}

func sliceArgumentsAction(f func([]string) error, validationMsg string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		args := ctx.Args().Slice()
		if len(args) == 0 {
			return fmt.Errorf(validationMsg)
		}

		return f(args)
	}
}
