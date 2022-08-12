package k8s

import (
	"context"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type chaosDeploymentfn func([]string, string, string, *kubernetes.Clientset)

func deploymentFn(namespace string, tag string, chaos string, number int) {
	//Checking if go-chaos needs to do anything
	if number == 0 {
		log.Println("Will not destroy any deployment")
		return
	}

	//Separating tags and adding "="
	var key, value string
	parts := strings.Split(tag, ":")
	key = parts[0]
	value = parts[1]
	label := key + "=" + value

	//Logging k8s
	clientset, _ := K8sConfig()

	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	//List deployment
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		log.Println("Error", err)
		return
	}
	var deploymentList []string

	for _, deployment := range list.Items {
		deploymentList = append(deploymentList, deployment.Name)
	}

	if number > len(deploymentList) {
		log.Println("Chaos not permitted", len(deploymentList), "deployments found with", key, value, "Number of deployment to destroy is:", number)
		return
	}

	deploymentList = deploymentList[:number]

	deploymentsMap := map[string]chaosDeploymentfn{
		"terminate": terminateDeploymentFn,
	}

	deploymentsMap[chaos](deploymentList, namespace, label, clientset)

}

func terminateDeploymentFn(deploymentList []string, namespace string, tags string, client *kubernetes.Clientset) {
	deploymentsClient := client.AppsV1().Deployments(namespace)

	for _, dplmnt := range deploymentList {
		err := deploymentsClient.Delete(context.TODO(), dplmnt, metav1.DeleteOptions{})
		log.Println("Terminating Deployment:", dplmnt)
		if err != nil {
			log.Println("Could not delete Deployment", err)
			return
		}
	}
}
