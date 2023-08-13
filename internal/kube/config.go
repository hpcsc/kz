package kube

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
)

func ContextsFromConfig() ([]string, error) {
	// loading rules able to handle KUBECONFIG variable if set
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil).
		RawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kube config: %v", err)
	}

	var contexts []string
	for c := range kubeConfig.Contexts {
		contexts = append(contexts, c)
	}

	return contexts, nil
}
