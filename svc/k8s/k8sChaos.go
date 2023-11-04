package k8s

import (
	"log"

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

func KubernetesChaos(namespace string, service string, tag string, chaos string, number int, dry bool) {
	k8sMap := map[string]k8sfn{
		"pod":        podFn,
		"deployment": deploymentFn,
		"daemonSet":  daemonFn,
	}
	if _, servExists := k8sMap[service]; servExists {
		k8sMap[service](namespace, tag, chaos, number, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
