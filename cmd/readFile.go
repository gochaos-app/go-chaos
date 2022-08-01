package cmd

import (
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ValidateFile(filename string) error {
	var valconfig GenConfig
	err := hclsimple.DecodeFile(filename, nil, &valconfig)
	if err != nil {
		log.Fatalln("Failed to load config", err)

		os.Exit(1)
		return err
	}
	log.Println("File readeable, you are good to execute chaos")
	return nil
}

func LoadConfig(filename string) (*GenConfig, error) {
	var config GenConfig
	err := hclsimple.DecodeFile(filename, nil, &config)
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
		return nil, err
	}
	return &config, nil
}
