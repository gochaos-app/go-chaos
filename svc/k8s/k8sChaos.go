package k8s

import (
	"log"

	"github.com/gochaos-app/go-chaos/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sfn func(string, string, string, int, bool)

func K8sConfig() (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("chaos from outside cluster")
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
		config, err = kubeconfig.ClientConfig()
		if err != nil {
			log.Println("Error getting config", err)
			return nil, err
		}
	} else {
		log.Println("chaos from inside cluster")
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	return clientset, nil
}

func KubernetesChaos(cfg config.JobConfig, dry bool) {
	k8sMap := map[string]k8sfn{
		"pod":        podFn,
		"deployment": deploymentFn,
		"daemonSet":  daemonFn,
	}
	if _, servExists := k8sMap[cfg.Service]; servExists {
		k8sMap[cfg.Service](cfg.Namespace, cfg.Chaos.Tag, cfg.Chaos.Chaos, cfg.Chaos.Count, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
