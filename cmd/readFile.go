package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

//func to decode file, if it fails it throws an error
func ValidateFile(filename string) error {
	var valconfig GenConfig
	err := hclsimple.DecodeFile(filename, nil, &valconfig)
	if err != nil {
		log.Fatalln("Failed to load config", err)
		return err
	}
	log.Println("File readeable, you are good to execute chaos")
	return nil
}

// Func to load config into memory, returns genconfig
func LoadConfig(filename string) (*GenConfig, error) {
	var config GenConfig
	err := hclsimple.DecodeFile(filename, nil, &config)
	if err != nil {
		log.Println("failed to load configuration: ", err)
		return nil, err
	}
	return &config, nil
}

//Func to execute a single job target out of a specified file
func ExecuteTarget(file string, target string) (*GenConfig, error) {
	var config GenConfig

	targetVars := strings.Split(target, ".")
	if len(targetVars) < 3 {
		err := errors.New("target must contain 3 values, separated by dot, example, aws.ec2.app:prod")
		log.Println("error: ", err)
		return nil, err
	}
	cloud, service, tag := targetVars[0], targetVars[1], targetVars[2]
	cfg, err := LoadConfig(file)
	if err != nil {
		log.Println("failed to load configuration: ", err)
		return nil, err
	}

	for j := 0; j < len(cfg.Job); j++ {
		if cfg.Job[j].Cloud == cloud && cfg.Job[j].Service == service && cfg.Job[j].Chaos.Tag == tag {
			switchService(cfg.Job[j])
			break
		}
	}
	return &config, nil
}
