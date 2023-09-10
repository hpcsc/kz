package cmd

import (
	"github.com/fatih/color"
	"github.com/hpcsc/kz/internal/gateway"
	"github.com/hpcsc/kz/internal/updater"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

func newUpdateSubcommand() *cli.Command {
	return &cli.Command{
		Name:        "update",
		Usage:       "Update to latest release version",
		Description: "Update to latest release version",
		Action:      update,
	}
}

func update(ctx *cli.Context) error {
	currentExecutable, err := os.Executable()
	if err != nil {
		return err
	}

	gw := gateway.NewGithubGateway()
	u := updater.New(runtime.GOARCH, currentExecutable, gw)
	msg, err := u.UpdateFrom(Version)
	if err != nil {
		return err
	}

	if msg != "" {
		color.Green(msg)
	}

	return nil
}
