package do

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/mental12345/chaosctl/ops"
)

type chaosLoadBalancerFn func(string, string, int, *godo.Client)

func LoadBalancerFn(config *godo.Client, tag string, chaos string, number int) {
	loadBalancerList, _, err := config.LoadBalancers.List(context.TODO(), &godo.ListOptions{})

	if err != nil {
		log.Println("error listing Load Balancers:", err)
		return
	}
	var lbName, lbID string
	for _, lb := range loadBalancerList {
		if lb.Name == tag {
			lbName = lb.Name
			lbID = lb.ID
		}
	}
	if lbName == "" || lbID == "" {
		log.Println("Couldn't find the load balancer")
	}
	LbMap := map[string]chaosLoadBalancerFn{
		"remove": removeDropletsFn,
	}
	LbMap[chaos](lbID, lbName, number, config)
}

func removeDropletsFn(id string, name string, number int, client *godo.Client) {

	LoadBalancers, _, err := client.LoadBalancers.Get(context.TODO(), id)
	if err != nil {
		log.Println("error listing droplets:", err)
		return
	}
	if number == 0 {
		log.Println("Will not execute removal of droplets")
	}
	if len(LoadBalancers.DropletIDs) < number {
		log.Println("Cannot remove from load balancer, existing droplets in load balancer:", len(LoadBalancers.DropletIDs))
		return
	}
	dropletsIDs := ops.Random(LoadBalancers.DropletIDs)
	dropletsIDs = dropletsIDs[:number]
	response, err := client.LoadBalancers.RemoveDroplets(context.TODO(), id, dropletsIDs...)
	if err != nil {
		log.Println("could not remove droplets from load balancer")
		return
	}
	log.Println(response)

}
