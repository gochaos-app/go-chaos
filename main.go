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
			if _, err := os.Stat(filename); err == nil {
				log.Println("Cannot read file, or file does not exist", err)
				os.Exit(1)
			}
			fmt.Println("Destroy initiated...")
			cmd.ReadFile()
			return nil
		},
		Commands: []cli.Command{
			{
				Name:    "Destroy",
				Aliases: []string{"d"},
				Usage:   "Execute destroy with custom file name",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err == nil {
						log.Println("Cannot read file, or file does not exist", err)
						os.Exit(1)
					}
					fmt.Println("Destroy initiated")
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
