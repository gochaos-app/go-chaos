package do

import (
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/gochaos-app/go-chaos/config"
)

type dofn func(*godo.Client, string, string, int, bool)

func DigitalOceanChaos(cfg config.JobConfig, dry bool) {
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
	if _, servExists := doMap[cfg.Service]; servExists {
		doMap[cfg.Service](client, cfg.Chaos.Tag, cfg.Chaos.Chaos, cfg.Chaos.Count, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
