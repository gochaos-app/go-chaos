package do

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/mental12345/go-chaos/ops"
)

type chaosDropletFn func([]int, int, *godo.Client)

func DropletFn(config *godo.Client, tag string, chaos string, number int) {
	droplets, _, err := config.Droplets.ListByTag(context.TODO(), tag, &godo.ListOptions{})
	if err != nil {
		log.Println("error listing droplets:", err)
		return
	}

	var DropletsInstances []int

	for _, droplet := range droplets {
		DropletsInstances = append(DropletsInstances, droplet.ID)
	}
	if len(DropletsInstances) >= number {
		log.Println("Droplet Chaos permitted...")
	} else {
		log.Println("Chaos not permitted", len(DropletsInstances), "droplets found with", tag, "Number of droplets to destroy is:", number)
		return
	}
	DropletMap := map[string]chaosDropletFn{
		"terminate": terminateDropletFn,
		"stop":      stopDropletFn,
		"poweroff":  poweroffDropletFn,
		"reboot":    rebootDropletFn,
	}

	if _, servExists := DropletMap[chaos]; servExists {
		DropletMap[chaos](ops.Random(DropletsInstances), number, config)
	} else {
		log.Println("Chaos not found")
		return
	}

}

func terminateDropletFn(dropletIDs []int, number int, client *godo.Client) {
	dropletIDs = dropletIDs[:number]
	for _, droplet := range dropletIDs {
		log.Println("Deleting droplet: ", droplet)
		_, err := client.Droplets.Delete(context.TODO(), droplet)
		if err != nil {
			log.Println("Error deleting droplet", droplet, ": ", err)
		}
	}
}

func stopDropletFn(dropletIDs []int, number int, client *godo.Client) {
	dropletIDs = dropletIDs[:number]
	for _, droplet := range dropletIDs {
		log.Println("shutting down droplet: ", droplet)
		_, _, err := client.DropletActions.Shutdown(context.TODO(), droplet)
		if err != nil {
			log.Println("Error Shutting down droplet", droplet, ": ", err)
		}
	}
}

func poweroffDropletFn(dropletIDs []int, number int, client *godo.Client) {
	dropletIDs = dropletIDs[:number]
	for _, droplet := range dropletIDs {
		log.Println("Powering off droplet: ", droplet)
		_, _, err := client.DropletActions.PowerOff(context.TODO(), droplet)
		if err != nil {
			log.Println("Error power off droplet", droplet, ": ", err)
		}
	}
}

func rebootDropletFn(dropletIDs []int, number int, client *godo.Client) {
	dropletIDs = dropletIDs[:number]
	for _, droplet := range dropletIDs {
		log.Println("Rebooting droplet: ", droplet)
		_, _, err := client.DropletActions.Reboot(context.TODO(), droplet)
		if err != nil {
			log.Println("Error rebooting droplet", droplet, ": ", err)
		}
	}
}
