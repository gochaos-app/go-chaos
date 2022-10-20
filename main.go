package main

import (
	"log"
	"os"
	"strings"

	"github.com/mental12345/go-chaos/cmd"
	"github.com/mental12345/go-chaos/server"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "go-chaos",
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
					if strings.HasPrefix(filename, "http") {
						test, err := cmd.ReadFromURL(filename)
						if err != nil {
							return err
						}
						cfg, err := cmd.LoadConfig(test)
						if err != nil {
							return err
						}
						if cmd.ExecuteChaos(cfg) != nil {
							return err
						}
						os.Exit(0)
					}

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
					if strings.HasPrefix(filename, "http") {
						test, err := cmd.ReadFromURL(filename)
						if err != nil {
							return err
						}
						cmd.ValidateFile(test)
						os.Exit(0)
					}
					if _, err := os.Stat(filename); err != nil {
						log.Println("Cannot read file, or file does not exist", err)
						os.Exit(1)
					}
					log.Println("Validation initiated")
					cmd.ValidateFile(filename)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start go-chaos server",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					log.Println("Server initiated")
					server.ServerFn(filename)
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
