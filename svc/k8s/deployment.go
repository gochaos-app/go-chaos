package k8s

import (
	"context"
	"log"
	"strings"

	"github.com/gochaos-app/go-chaos/ops"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type chaosDeploymentfn func([]string, string, string, int, *kubernetes.Clientset)

func deploymentFn(namespace string, tag string, chaos string, number int) {

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

	deploymentsMap := map[string]chaosDeploymentfn{
		"terminate": terminateDeploymentFn,
		"update":    updateDeploymentFn,
	}
	if _, servExists := deploymentsMap[chaos]; servExists {
		deploymentsMap[chaos](ops.Random(deploymentList), namespace, label, number, clientset)
	} else {
		log.Println("Chaos not found")
		return
	}
}

func terminateDeploymentFn(deploymentList []string, namespace string, tags string, number int, client *kubernetes.Clientset) {

	if number <= 0 {
		log.Println("Will not perform chaos on any deployment")
		return
	}

	if number > len(deploymentList) {
		log.Println("Chaos not permitted", len(deploymentList), "deployments found.", "Number of deployment to destroy is:", number)
		return
	}

	deploymentList = deploymentList[:number]
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

func updateDeploymentFn(deploymentList []string, namespace string, tags string, number int, client *kubernetes.Clientset) {
	if len(deploymentList) > 1 {
		log.Println("Chaos not permitted, when updating only one deployment with specified labels should exist, deployments found:", len(deploymentList))
		return
	}
	if number < 0 {
		log.Println("Cannot update to negative number")
		return
	}
	deploymentsClient := client.AppsV1().Deployments(namespace)
	deployment := deploymentList[0]

	scale, err := deploymentsClient.GetScale(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		log.Println("error:", err)
	}
	sc := *scale
	sc.Spec.Replicas = int32(number)
	update, err := deploymentsClient.UpdateScale(context.TODO(), deployment, &sc, metav1.UpdateOptions{})
	if err != nil {
		log.Println("error:", err)
	}
	log.Println("Updating:", deployment, "to:", update.Spec.Replicas)

}
