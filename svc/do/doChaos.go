package do

import (
	"log"
	"os"

	"github.com/digitalocean/godo"
)

type dofn func(*godo.Client, string, string, int, bool)

func DigitalOceanChaos(region string, service string, tag string, chaos string, number int, dry bool) {
	// Check for environment variable
	token := os.Getenv("DIGITALOCEAN_TOKEN")
	if len(token) == 0 {
		log.Println("Cannot find DigitalOcean token variable, please set DIGITALOCEAN_TOKEN variable in system")
		return
	}
	//logs from token to digital ocean
	client := godo.NewFromToken(token)

	doMap := map[string]dofn{
		"droplet":       DropletFn,
		"load_balancer": LoadBalancerFn,
	}
	if _, servExists := doMap[service]; servExists {
		doMap[service](client, tag, chaos, number, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
