package k8s

import (
	"context"
	"log"
	"strings"

	"github.com/gochaos-app/go-chaos/ops"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type chaosPodfn func([]string, string, string, *kubernetes.Clientset)

func podFn(namespace string, tag string, chaos string, number int, dry bool) {

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
		log.Println("Error:", err)
		return
	}

	var podList []string

	for _, pod := range list.Items {
		podList = append(podList, pod.Name)
	}

	if len(podList) == 0 {
		log.Println("Could not find any pods with", label)
		return
	}

	podList = podList[:number]

	if dry {
		log.Println("Dry mode")
		log.Println("Will apply chaos on pods", podList)
		return
	}

	//Checking if go-chaos needs to do anything
	if number <= 0 {
		log.Println("Will not destroy any Pod")
		return
	}

	if number > len(podList) {
		log.Println("Chaos not permitted", len(podList), "pods found with", key, value, "Number of pods to destroy is:", number)
		return
	}

	podsMap := map[string]chaosPodfn{
		"terminate":    terminatePodFn,
		"terminateAll": terminateAllFn,
	}
	if _, servExists := podsMap[chaos]; servExists {
		podsMap[chaos](ops.RandomArray(podList), namespace, label, clientset)
	} else {
		log.Println("Chaos not found")
		return
	}

}

func terminatePodFn(podList []string, namespace string, tags string, client *kubernetes.Clientset) {
	pods := client.CoreV1().Pods(namespace)
	for _, pod := range podList {
		err := pods.Delete(context.TODO(), pod, metav1.DeleteOptions{})
		log.Println("Terminating pod:", pod)
		if err != nil {
			log.Println("Could not delete pod", err)
			return
		}
	}
}

func terminateAllFn(podList []string, namespace string, tags string, client *kubernetes.Clientset) {
	pods := client.CoreV1().Pods(namespace)
	log.Println("Terminating collection with labels:", tags)
	err := pods.DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: tags,
	})
	if err != nil {
		log.Println("Could not delete pod collection", err)
		return
	}
}
