package k8s

import (
	"log"

	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sfn func(string, string, string, int, bool)

func K8sConfig() (*kubernetes.Clientset, error) {

	rules := clientcmd.NewDefaultClientConfigLoadingRules()

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalln("Error: ", err)
		return nil, err
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	/*	defaultCfg := filepath.Join(os.Getenv("HOME"), ".kube", "config")
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
		}*/
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
