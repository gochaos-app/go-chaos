package do

import (
	"log"
	"os"

	"github.com/digitalocean/godo"
)

type dofn func(*godo.Client, string, string, int)

func DigitalOceanChaos(region string, service string, tag string, chaos string, number int) {
	// Check for environment variable
	token := os.Getenv("DIGITALOCEAN_TOKEN")
	if len(token) == 0 {
		log.Println("Cannot find DigitalOcean token variable, please set DIGITALOCEAN_TOKEN variable in system")
		return
	}
	//logs from token to digital ocean
	client := godo.NewFromToken(token)

	if number == 0 {
		log.Println("Will not destroy any droplet")
		return
	}

	awsMap := map[string]dofn{
		"droplet": DropletFn,
	}
	awsMap[service](client, tag, chaos, number)
}
