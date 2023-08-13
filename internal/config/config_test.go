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
			Namespaces: []string{
				"ns1",
				"ns2",
			},
		})
		require.NoError(t, err)

		content, err := os.ReadFile(destinationFile)
		require.NoError(t, err)
		require.Contains(t, string(content), "ns1")
		require.Contains(t, string(content), "ns2")
	})
}
