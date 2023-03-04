package main

import (
	"errors"
	"log"
	"os"

	"github.com/gochaos-app/go-chaos/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "go-chaos",
		Usage: "a terminal based app that injects chaos into your cloud infrastrucure",
		Commands: []cli.Command{
			{
				Name:  "destroy",
				Usage: "Execute destroy with custom file name",
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
					if cmd.ExecuteChaos(cfg, false) != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "plan",
				Usage: "Execute a dry run with custom file name",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					if _, err := os.Stat(filename); err != nil {
						err := errors.New("cannot read" + filename + " or doesn't exists")
						return err
					}
					log.Println("Go-Chaos dry run initiated")
					cfg, err := cmd.LoadConfig(filename)
					if err != nil {
						return err
					}
					if cmd.ExecuteChaos(cfg, true) != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "validate",
				Usage: "Validate chaos file",
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
				Name:  "target",
				Usage: "Execute chaos on a single target",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					target := c.Args().Get(1)
					cmd.ExecuteTarget(file, target)
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
