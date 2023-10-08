package gcp

import (
	"context"
	"errors"
	"log"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/gochaos-app/go-chaos/ops"
	"google.golang.org/api/iterator"

	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type chaosVMfn func([]string, string, string, int, *compute.InstancesClient) error

func instanceFn(project string, region string, tag string, chaos string, number int, dry bool) error {

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
		err := errors.New("Couldn't get instance list")
		return err
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
			return err
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				vms = append(vms, instance.GetName())
			}
		}
	}

	// FIXME: Implement a better way to handle below logic
	if len(vms) == 0 {
		err := errors.New("Could not find any gcp vm with")
		return err
	}
	if dry {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", vms)
		return nil
	}

	if len(vms) >= number {
		log.Println("GCP VM Chaos permitted...")
	} else {
		err := errors.New("Chaos not permitted: instances found is smaller than the number of instances to destroy")
		return err
	}

	vmsMap := map[string]chaosVMfn{
		"terminate": terminateVMFn,
		"stop":      stopVMFn,
		"reset":     restartVMFn,
	}
	if _, servExists := vmsMap[chaos]; servExists {
		vmsMap[chaos](ops.RandomArray(vms), region, project, number, instanceClient)
	} else {
		err := errors.New("Chaos not found")
		return err
	}

	return nil
}

func terminateVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) error {
	ctx := context.Background()
	if number <= 0 {
		err := errors.New("Will not destroy any VM")
		return err
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
			err := errors.New("unable to delete instance: ")
			return err
		}

		if err = op.Wait(ctx); err != nil {
			err := errors.New("unable to wait for the operation:")
			return err
		}
		log.Println("Instance Deleted", vm)
	}

	return nil
}

func stopVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) error {
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
			return err
		}

		if err = op.Wait(ctx); err != nil {
			return err
		}
		log.Println("Instance Stopped", vm)
	}

	return nil
}

func restartVMFn(vmList []string, zone string, project string, number int, cfg *compute.InstancesClient) error {
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
			return err
		}

		if err = op.Wait(ctx); err != nil {
			return err
		}
		log.Println("Instance Reset", vm)
	}

	return nil
}
