package do

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/gochaos-app/go-chaos/ops"
)

type chaosLoadBalancerFn func(string, string, int, *godo.Client)

func LoadBalancerFn(config *godo.Client, tag string, chaos string, number int, dry bool) {
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
		return
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on LoadBalancer ", lbName)
		return
	}
	if number <= 0 {
		log.Println("Will not destroy any load balancer resource")
		return
	}

	LbMap := map[string]chaosLoadBalancerFn{
		"removeDroplets": removeDropletsFn,
		"removeRules":    removeRulesFn,
	}

	if _, servExists := LbMap[chaos]; servExists {
		LbMap[chaos](lbID, lbName, number, config)
	} else {
		log.Println("Chaos not found")
		return
	}

}

func removeRulesFn(id string, name string, number int, client *godo.Client) {
	LoadBalancers, _, err := client.LoadBalancers.Get(context.TODO(), id)
	if err != nil {
		log.Println("error getting load balancer", err)
	}
	if number == 0 {
		log.Println("Will not execute removal of rules")
		return
	}
	if len(LoadBalancers.ForwardingRules) < number {
		log.Println("Cannot remove from load balancer, existing rules in load balancer:", len(LoadBalancers.ForwardingRules))
		return
	}

	rulesIDs := LoadBalancers.ForwardingRules
	rulesIDs = rulesIDs[:number]
	_, err = client.LoadBalancers.RemoveForwardingRules(context.TODO(), id, rulesIDs...)
	if err != nil {
		log.Println("could not remove forwarding rules from load balancer: ", err)
		return
	}
}

func removeDropletsFn(id string, name string, number int, client *godo.Client) {

	LoadBalancers, _, err := client.LoadBalancers.Get(context.TODO(), id)
	if err != nil {
		log.Println("error getting load balancer:", err)
		return
	}
	if number == 0 {
		log.Println("Will not execute removal of droplets")
		return
	}
	if len(LoadBalancers.DropletIDs) < number {
		log.Println("Cannot remove from load balancer, existing droplets in load balancer:", len(LoadBalancers.DropletIDs))
		return
	}
	dropletsIDs := ops.RandomArray(LoadBalancers.DropletIDs)
	dropletsIDs = dropletsIDs[:number]
	_, err = client.LoadBalancers.RemoveDroplets(context.TODO(), id, dropletsIDs...)
	if err != nil {
		log.Println("could not remove droplets from load balancer: ", err)
		return
	}
}
