//go:build unit

package config

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	t.Run("do not add duplicate namespaces", func(t *testing.T) {
		c := Config{
			Namespaces: []string{
				"ns1",
				"ns2",
			},
		}

		c.AddNamespaces("ns2", "ns3")

		require.Equal(t, []string{
			"ns1",
			"ns2",
			"ns3",
		}, c.Namespaces)
	})

	t.Run("able to delete multiple namespaces", func(t *testing.T) {
		c := Config{
			Namespaces: []string{
				"ns1",
				"ns2",
				"ns3",
				"ns4",
			},
		}

		c.DeleteNamespaces("ns1", "ns3")

		require.Equal(t, []string{
			"ns2",
			"ns4",
		}, c.Namespaces)
	})

	t.Run("return all contexts that partially match given query", func(t *testing.T) {
		c := Config{
			Contexts: []string{
				"context1",
				"context2",
				"context3",
				"context4",
			},
		}

		contexts := c.ContextsMatching("2")

		require.Equal(t, []string{"context2"}, contexts)
	})
}

func TestLoad(t *testing.T) {
	t.Run("return empty config when not found", func(t *testing.T) {
		c, err := Load("not-existing")

		require.NoError(t, err)
		require.Equal(t, &Config{}, c)
	})

	t.Run("return config when found", func(t *testing.T) {
		c, err := Load("testdata/config.yaml")

		require.NoError(t, err)
		require.Equal(t, &Config{
			Contexts: []string{
				"context1",
				"context2",
			},
			Namespaces: []string{
				"ns1",
				"ns2",
			},
		}, c)
	})
}

func TestSave(t *testing.T) {
	t.Run("save config to filesystem", func(t *testing.T) {
		destinationFile := path.Join(os.TempDir(), fmt.Sprintf("kz-config-%d.yaml", time.Now().UnixMilli()))
		defer os.Remove(destinationFile)

		err := Save(destinationFile, &Config{
			Contexts: []string{
				"context1",
				"context2",
			},
			Namespaces: []string{
				"ns1",
				"ns2",
			},
		})
		require.NoError(t, err)

		contentAsByte, err := os.ReadFile(destinationFile)
		content := string(contentAsByte)
		require.NoError(t, err)
		require.Contains(t, content, "context1")
		require.Contains(t, content, "context2")
		require.Contains(t, content, "ns1")
		require.Contains(t, content, "ns2")
	})
}
