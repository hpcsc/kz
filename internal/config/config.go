package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"slices"
)

type Config struct {
	Namespaces []string
}

func (c *Config) AddNamespaces(namespaces ...string) {
	for _, n := range namespaces {
		if !slices.Contains(c.Namespaces, n) {
			c.Namespaces = append(c.Namespaces, n)
		}
	}
}

func LoadFromDefaultLocation() (*Config, error) {
	location, err := defaultConfigLocation()
	if err != nil {
		return nil, err
	}

	return Load(location)
}

func Load(location string) (*Config, error) {
	content, err := os.ReadFile(location)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}

		return nil, fmt.Errorf("failed to load yaml config from location %s: %v", location, err)
	}

	var c Config
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml config from location %s: %v", location, err)
	}

	return &c, nil
}

func SaveToDefaultLocation(c *Config) error {
	location, err := defaultConfigLocation()
	if err != nil {
		return err
	}

	return Save(location, c)
}

func Save(location string, c *Config) error {
	marshalled, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config to yaml: %v", err)
	}

	if err := os.WriteFile(location, marshalled, 0644); err != nil {
		return fmt.Errorf("failed to write config to yaml file at %s: %v", location, err)
	}

	return nil
}

func defaultConfigLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home directory: %v", err)
	}

	return path.Join(home, ".kz.yml"), nil
}
