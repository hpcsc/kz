package tui

import "github.com/pterm/pterm"

func ShowDropdown(label string, options []string) (string, error) {
	return pterm.DefaultInteractiveSelect.
		WithDefaultText(label).
		WithOptions(options).
		Show()
}
