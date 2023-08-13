//go:build unit

package kube

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
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

func TestSwitchContextTo(t *testing.T) {
	t.Run("set current context to given context", func(t *testing.T) {
		destinationConfigPath := path.Join(os.TempDir(), fmt.Sprintf("kz-kube-config-%d", time.Now().UnixMilli()))
		copyFile(t, "testdata/kubeconfig-1", destinationConfigPath)
		defer os.Remove(destinationConfigPath)

		err := SwitchContextTo("context-3", destinationConfigPath)
		require.NoError(t, err)

		content, err := os.ReadFile(destinationConfigPath)
		require.NoError(t, err)
		require.Contains(t, string(content), "current-context: context-3")
	})
}

func copyFile(t *testing.T, sourcePath string, destinationPath string) {
	data, err := os.ReadFile(sourcePath)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(destinationPath, data, 0644))
}
