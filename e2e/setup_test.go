//go:build e2e

package e2e

import (
	"github.com/hpcsc/kz/internal/cmd"
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"kz": cmd.Run,
	}))
}
