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
	t.Run("return error when context to switch is empty", func(t *testing.T) {
		err := SwitchContextTo("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "context to switch to is required")
	})

	t.Run("return error when context to switch to not exists in config files", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchContextTo("context-3")

		require.Error(t, err)
		require.Contains(t, err.Error(), "context with name context-3 does not exist in kube config file(s)")
	})

	t.Run("set current context to given context", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchContextTo("context-2")
		require.NoError(t, err)

		content, err := os.ReadFile(destinationConfigPath)
		require.NoError(t, err)
		require.Contains(t, string(content), "current-context: context-2")
	})

	t.Run("set current context to given context when multiple config files available", func(t *testing.T) {
		config1Path := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(config1Path)
		config2Path := copyFileToTmp(t, "testdata/kubeconfig-2")
		defer os.Remove(config2Path)

		os.Setenv("KUBECONFIG", fmt.Sprintf("%s:%s", config1Path, config2Path))
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchContextTo("context-3")
		require.NoError(t, err)

		content, err := os.ReadFile(config1Path)
		require.NoError(t, err)
		require.Contains(t, string(content), "current-context: context-3")
	})
}

func TestSwitchNamespaceTo(t *testing.T) {
	t.Run("return error when namespace to switch to is empty", func(t *testing.T) {
		err := SwitchNamespaceTo("")

		require.Error(t, err)
		require.Contains(t, err.Error(), "namespace to switch to is required")
	})

	t.Run("return error when current context is not set", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchNamespaceTo("ns2")

		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to switch namespace to ns2 because current context is not set")
	})

	t.Run("set namespace of current context", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-3")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchNamespaceTo("ns2")
		require.NoError(t, err)

		content, err := os.ReadFile(destinationConfigPath)
		require.NoError(t, err)
		require.Contains(t, string(content), "namespace: ns2")
	})
}

func TestSwitchContextAndNamespace(t *testing.T) {
	t.Run("return error when context to switch is empty", func(t *testing.T) {
		err := SwitchContextAndNamespace("", "ns1")

		require.Error(t, err)
		require.Contains(t, err.Error(), "context to switch to is required")
	})

	t.Run("return error when context to switch to not exists in config files", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchContextAndNamespace("context-3", "ns1")

		require.Error(t, err)
		require.Contains(t, err.Error(), "context with name context-3 does not exist in kube config file(s)")
	})

	t.Run("set current context and namespace", func(t *testing.T) {
		destinationConfigPath := copyFileToTmp(t, "testdata/kubeconfig-1")
		defer os.Remove(destinationConfigPath)

		os.Setenv("KUBECONFIG", destinationConfigPath)
		defer os.Unsetenv("KUBECONFIG")

		err := SwitchContextAndNamespace("context-2", "ns2")
		require.NoError(t, err)

		content, err := os.ReadFile(destinationConfigPath)
		require.NoError(t, err)
		require.Contains(t, string(content), "current-context: context-2")
		require.Contains(t, string(content), "namespace: ns2")
	})
}

func copyFileToTmp(t *testing.T, sourcePath string) string {
	destinationPath := path.Join(os.TempDir(), fmt.Sprintf("kz-kube-config-%d", time.Now().UnixMilli()))
	data, err := os.ReadFile(sourcePath)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(destinationPath, data, 0644))
	return destinationPath
}
