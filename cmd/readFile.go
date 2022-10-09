package cmd

import (
	"io"
	"log"
	"net/http"
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

func ReadFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	tmpNameFile := "/tmp/config.hcl"
	tempFile, err := os.Create(tmpNameFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer tempFile.Close()
	_, err2 := tempFile.WriteString(string(body))
	if err2 != nil {
		log.Fatalln(err)
	}
	return tmpNameFile, nil
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
