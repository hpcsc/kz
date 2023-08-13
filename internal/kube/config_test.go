//go:build unit

package kube

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestContextsFromConfig(t *testing.T) {
	// only able to deterministically test the case where KUBECONFIG is set
	t.Run("return contexts from multiple config files specified in KUBECONFIG variable", func(t *testing.T) {
		os.Setenv("KUBECONFIG", "testdata/kubeconfig-1:testdata/kubeconfig-2")
		defer os.Unsetenv("KUBECONFIG")

		contexts, err := ContextsFromConfig()

		require.NoError(t, err)
		require.ElementsMatch(t, []string{
			"context-1",
			"context-2",
			"context-3",
		}, contexts)
	})
}
