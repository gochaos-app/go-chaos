package k8s

import (
	"context"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func podFn(namespace string, tag string, number int) {
	//Checking if go-chaos needs to do anything
	if number == 0 {
		log.Println("Will not destroy any Pod")
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

	//List pods
	podClient := clientset.CoreV1().Pods(namespace)

	list, err := podClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		log.Println("Error", err)
		return
	}
	var podList []string

	for _, pod := range list.Items {
		podList = append(podList, pod.Name)
	}

	if number > len(podList) {
		log.Println("Chaos not permitted", len(podList), "pods found with", key, value, "Number of pods to destroy is:", number)
		return
	}
	podList = podList[:number]
	for _, pod := range podList {
		err := podClient.Delete(context.TODO(), pod, metav1.DeleteOptions{})
		log.Println("Terminating pod:", pod)
		if err != nil {
			log.Println("Could not delete pod", err)
			return
		}
	}
}
