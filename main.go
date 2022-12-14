package main

import (
	"errors"
	"log"
	"os"

	"github.com/gochaos-app/go-chaos/cmd"
	"github.com/gochaos-app/go-chaos/server"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "go-chaos",
		Usage: "a terminal based app that injects chaos into your cloud infrastrucure",
		Commands: []cli.Command{
			{
				Name:    "destroy",
				Aliases: []string{"d"},
				Usage:   "Execute destroy with custom file name",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err != nil {
						err := errors.New("cannot read" + filename + " or doesn't exists")
						return err
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
				Usage:   "Validate chaos file",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err != nil {
						err := errors.New("cannot read" + filename + " or doesn't exists")
						return err
					}
					log.Println("Validation initiated")
					cmd.ValidateFile(filename)
					return nil
				},
			},
			{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "Execute chaos on a single target",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					target := c.Args().Get(1)
					cmd.ExecuteTarget(file, target)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start go-chaos server",
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
