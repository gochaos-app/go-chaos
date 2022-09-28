package k8s

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sfn func(string, string, string, int)

func K8sConfig() (*kubernetes.Clientset, error) {
	defaultCfg := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	var config *rest.Config

	config, err := clientcmd.BuildConfigFromFlags("", defaultCfg)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil, err
	}
	return clientset, nil
}

func KubernetesChaos(namespace string, service string, tag string, chaos string, number int) {
	k8sMap := map[string]k8sfn{
		"pod":        podFn,
		"deployment": deploymentFn,
	}
	if _, servExists := k8sMap[service]; servExists {
		k8sMap[service](namespace, tag, chaos, number)
	} else {
		log.Println("Service not found")
		return
	}

}
