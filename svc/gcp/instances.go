package gcp

import (
	"context"
	"log"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/gochaos-app/go-chaos/ops"
	"google.golang.org/api/iterator"

	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type chaosVMfn func([]string, string, string, int, *compute.InstancesClient)

func instanceFn(project string, region string, tag string, chaos string, number int, dry bool) {

	// Separate tag string into key value components
	parts := strings.Split(tag, ":")
	var key, value string
	key = parts[0]
	value = parts[1]

	var filters *string

	//status := "status "
	stringFilter := "(labels." + key + ":" + value + ") AND (status = RUNNING)"
	filters = &stringFilter

	// Get list of vm instances in provided region and project
	ctx := context.Background()
	instanceClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Println("Couldn't get instance list", err)
		return
	}
	defer instanceClient.Close()
	req := &computepb.AggregatedListInstancesRequest{
		Project: project,
		Filter:  filters,
	}
	//how fast can I write on this small keyboard
	it := instanceClient.AggregatedList(ctx, req)
	var vms []string
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				vms = append(vms, instance.GetName())
			}
		}
	}
	if len(vms) == 0 {
		log.Println("Could not find any gcp vm with", tag)
		return
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", vms)
		return
	}

	if len(vms) >= number {
		log.Println("GCP VM Chaos permitted...")
	} else {
		log.Println("Chaos not permitted", len(vms), "instances found with", key, value, " Number of instances to destroy is:", number)
		return
	}

	vmsMap := map[string]chaosVMfn{
		"terminate": terminateVMFn,
		"stop":      stopVMFn,
		"reset":     restartVMFn,
	}
	if _, servExists := vmsMap[chaos]; servExists {
		vmsMap[chaos](ops.RandomArray(vms), region, project, number, instanceClient)
	} else {
		log.Println("Chaos not found")
		return
	}
}

func terminateVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) {
	ctx := context.Background()
	if number <= 0 {
		log.Println("Will not destroy any VM")
		return
	}

	vmList = vmList[:number]
	for _, vm := range vmList {
		req := &computepb.DeleteInstanceRequest{
			Project:  project,
			Zone:     zone,
			Instance: vm,
		}
		op, err := cfg.Delete(ctx, req)
		if err != nil {
			log.Println("unable to delete instance: ", err)
			return
		}

		if err = op.Wait(ctx); err != nil {
			log.Println("unable to wait for the operation:", err)
			return
		}
		log.Println("Instance Deleted", vm)
	}
}

func stopVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) {
	ctx := context.Background()

	vmList = vmList[:number]
	for _, vm := range vmList {
		req := &computepb.StopInstanceRequest{
			Project:  project,
			Zone:     zone,
			Instance: vm,
		}
		op, err := cfg.Stop(ctx, req)
		if err != nil {
			log.Println("unable to stop instance: ", err)
			return
		}

		if err = op.Wait(ctx); err != nil {
			log.Println("unable to wait for the operation:", err)
			return
		}
		log.Println("Instance Stopped", vm)
	}
}

func restartVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) {
	ctx := context.Background()

	vmList = vmList[:number]
	for _, vm := range vmList {
		req := &computepb.ResetInstanceRequest{
			Project:  project,
			Zone:     zone,
			Instance: vm,
		}
		op, err := cfg.Reset(ctx, req)
		if err != nil {
			log.Println("unable to reset instance: ", err)
			return
		}

		if err = op.Wait(ctx); err != nil {
			log.Println("unable to wait for the operation:", err)
			return
		}
		log.Println("Instance Reset", vm)
	}
}
