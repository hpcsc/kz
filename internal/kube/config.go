package kube

import (
	"errors"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func ContextsFromConfig() ([]string, error) {
	ca := clientcmd.NewDefaultPathOptions()
	cfg, err := ca.GetStartingConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get starting config: %v", err)
	}

	var contexts []string
	for c := range cfg.Contexts {
		contexts = append(contexts, c)
	}

	return contexts, nil
}

func SwitchContextTo(ctx string) error {
	if len(ctx) == 0 {
		return errors.New("context to switch to is required")
	}

	ca := clientcmd.NewDefaultPathOptions()
	cfg, err := ca.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("failed to get starting config: %v", err)
	}

	if !contextExists(cfg.Contexts, ctx) {
		return fmt.Errorf("context with name %s does not exist in kube config file(s)", ctx)
	}

	cfg.CurrentContext = ctx

	if err := clientcmd.ModifyConfig(ca, *cfg, true); err != nil {
		return fmt.Errorf("failed to modify config: %v", err)
	}

	return nil
}

func SwitchNamespaceTo(namespace string) error {
	if len(namespace) == 0 {
		return errors.New("namespace to switch to is required")
	}

	ca := clientcmd.NewDefaultPathOptions()
	cfg, err := ca.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("failed to get starting config: %v", err)
	}

	if len(cfg.CurrentContext) == 0 {
		return fmt.Errorf("unable to switch namespace to %s because current context is not set", namespace)
	}

	cfg.Contexts[cfg.CurrentContext].Namespace = namespace

	if err := clientcmd.ModifyConfig(ca, *cfg, true); err != nil {
		return fmt.Errorf("failed to modify config: %v", err)
	}

	return nil
}

func SwitchContextAndNamespace(ctx string, namespace string) error {
	if len(ctx) == 0 {
		return errors.New("context to switch to is required")
	}

	ca := clientcmd.NewDefaultPathOptions()
	cfg, err := ca.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("failed to get starting config: %v", err)
	}

	if !contextExists(cfg.Contexts, ctx) {
		return fmt.Errorf("context with name %s does not exist in kube config file(s)", ctx)
	}

	cfg.CurrentContext = ctx
	cfg.Contexts[ctx].Namespace = namespace

	if err := clientcmd.ModifyConfig(ca, *cfg, true); err != nil {
		return fmt.Errorf("failed to modify config: %v", err)
	}

	return nil
}

func contextExists(contexts map[string]*api.Context, ctx string) bool {
	for c := range contexts {
		if c == ctx {
			return true
		}
	}

	return false
}
