package k8s

import (
	"context"
	"log"
	"strings"

	"github.com/mental12345/go-chaos/ops"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type chaosDaemonfn func([]string, string, string, int, *kubernetes.Clientset)

func daemonFn(namespace string, tag string, chaos string, number int) {

	//Separating tags and adding "="
	var key, value string
	parts := strings.Split(tag, ":")
	key = parts[0]
	value = parts[1]
	label := key + "=" + value

	//Logging k8s
	clientset, _ := K8sConfig()

	daemonsClient := clientset.AppsV1().DaemonSets(namespace)
	//List daemon
	list, err := daemonsClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		log.Println("Error", err)
		return
	}
	var daemonList []string

	for _, daemonSet := range list.Items {
		daemonList = append(daemonList, daemonSet.Name)
	}

	daemonsMap := map[string]chaosDaemonfn{
		"terminate": terminateDaemonFn,
	}
	if _, servExists := daemonsMap[chaos]; servExists {
		daemonsMap[chaos](ops.Random(daemonList), namespace, label, number, clientset)
	} else {
		log.Println("Chaos not found")
		return
	}
}

func terminateDaemonFn(daemonList []string, namespace string, tags string, number int, client *kubernetes.Clientset) {

	if number == 0 {
		log.Println("Will not perform chaos on any daemonSet")
		return
	}

	if number > len(daemonList) {
		log.Println("Chaos not permitted", len(daemonList), "daemonsets found.", "Number of daemonSets to destroy is:", number)
		return
	}

	daemonList = daemonList[:number]
	daemonsClient := client.AppsV1().DaemonSets(namespace)

	for _, dplmnt := range daemonList {
		err := daemonsClient.Delete(context.TODO(), dplmnt, metav1.DeleteOptions{})
		log.Println("Terminating DaemonSet:", dplmnt)
		if err != nil {
			log.Println("Could not delete DaemonSet", err)
			return
		}
	}
}
