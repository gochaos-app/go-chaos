package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mental12345/chaosCLI/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "goChaos",
		Usage: "A terminal based app that injects chaos into your cloud infrastrucure",
		Action: func(c *cli.Context) error {
			filename := "config.hcl"
			if _, err := os.Stat(filename); err != nil {
				log.Println("Cannot read file, or file does not exist", err)
				os.Exit(1)
			}
			fmt.Println("Destroy initiated...")
			cfg, err := cmd.LoadConfig(filename)
			if err != nil {
				return err
			}
			if cmd.ExecuteChaos(cfg) != nil {
				return err
			}

			return nil
		},
		Commands: []cli.Command{
			{
				Name:    "destroy",
				Aliases: []string{"d"},
				Usage:   "Execute destroy with custom file name",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err != nil {
						log.Println("Cannot read file, or file does not exist", err)
						os.Exit(1)
					}
					fmt.Println("Destroy initiated")
					cmd.LoadConfig(filename)
					return nil
				},
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "validate file",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err != nil {
						log.Println("Cannot read file, or file does not exist", err)
						os.Exit(1)
					}
					fmt.Println("Validation initiated")
					cmd.ValidateFile(filename)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
