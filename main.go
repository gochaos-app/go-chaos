package main

import (
	"log"
	"os"

	"github.com/mental12345/chaosctl/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "chaosctl",
		Usage: "A terminal based app that injects chaos into your cloud infrastrucure",
		Action: func(c *cli.Context) error {
			filename := "config.hcl"
			if _, err := os.Stat(filename); err != nil {
				log.Println("Cannot read file, or file does not exist", err)
				os.Exit(1)
			}
			log.Println("Validation initiated")
			cmd.ValidateFile(filename)
			log.Println("Destroy initiated...")
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
					log.Println("Destroy initiated")
					cfg, err := cmd.LoadConfig(filename)
					if err != nil {
						return err
					}
					if cmd.ExecuteChaos(cfg) != nil {
						return err
					}
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
					log.Println("Validation initiated")
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
