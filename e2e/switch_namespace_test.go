//go:build e2e

package e2e

import (
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"testing"
)

func TestSwitchNamespace(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/switch_namespace",
		Setup: func(env *testscript.Env) error {
			env.Setenv("HOME", os.TempDir())
			return nil
		},
	})
}
