package module

import (
	"log"
	"os"
	"plugin"
)

func ModulesChaos(region string, service string, project string, namespace string, tag string, chaos string, number int, dry bool) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic ocurred", err)
		}
	}()

	if _, err := os.Stat(service); err == nil {
		log.Println("Plugin dir found")
	} else {
		log.Println("Plugin dir not found")
		return
	}
	// Load module
	module := service + "/module.so"
	if _, err := os.Stat(module); err == nil {
		log.Println("Module found")
	} else {
		log.Println("Module not found")
		return
	}
	plug, err := plugin.Open(module)
	if err != nil {
		log.Println("Error opening the module", err)
		return
	}
	if dry {
		log.Println("Dry mode")
		log.Println("Will not execute chaos actions")
		return
	}
	symChaos, err := plug.Lookup("ChaosFunc")
	if err != nil {
		log.Println("No Chaos symbol")
		return
	}

	symChaos.(func(
		string,
		string,
		string,
		string,
		string,
		string,
		int))(region, service, project, namespace, tag, chaos, number)

}
