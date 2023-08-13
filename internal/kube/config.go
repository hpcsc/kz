package kube

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func ContextsFromConfig() ([]string, error) {
	kubeConfig, err := config()
	if err != nil {
		return nil, err
	}

	var contexts []string
	for c := range kubeConfig.Contexts {
		contexts = append(contexts, c)
	}

	return contexts, nil
}

func SwitchContextTo(ctx string, destinationConfigPath string) error {
	kubeConfig, err := clientcmd.LoadFromFile(destinationConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load kube config from %s: %v", destinationConfigPath, err)
	}

	kubeConfig.CurrentContext = ctx

	if err := clientcmd.WriteToFile(*kubeConfig, destinationConfigPath); err != nil {
		return fmt.Errorf("failed to write updated kube config to %s: %v", destinationConfigPath, err)
	}

	return nil
}

func SwitchNamespaceTo(namespace string, destinationConfigPath string) error {
	kubeConfig, err := clientcmd.LoadFromFile(destinationConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load kube config from %s: %v", destinationConfigPath, err)
	}

	if len(kubeConfig.CurrentContext) == 0 {
		return fmt.Errorf("unable to switch namespace to %s because current context is not set", namespace)
	}

	kubeConfig.Contexts[kubeConfig.CurrentContext].Namespace = namespace

	if err := clientcmd.WriteToFile(*kubeConfig, destinationConfigPath); err != nil {
		return fmt.Errorf("failed to write updated kube config to %s: %v", destinationConfigPath, err)
	}

	return nil
}

func config() (api.Config, error) {
	// loading rules able to handle KUBECONFIG variable if set
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil).
		RawConfig()
	if err != nil {
		return api.Config{}, fmt.Errorf("failed to load kube config: %v", err)
	}

	return kubeConfig, nil
}
