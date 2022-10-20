package server

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type server struct {
	Port   string `toml:"port"`
	Remote string `toml:"remote"`
	Branch string `toml:"branch"`
	Path   string `toml:"path"`
}

func readConfig(filename string) (*server, error) {
	if _, err := os.Stat(filename); err != nil {
		log.Println("no server configuration found")
		os.Exit(1)
	}
	var config server

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	if config.Path == "" {
		log.Println("Error: path property is empty")
		os.Exit(1)
	}

	if config.Remote != "" {
		cloneRepo(config.Remote, config.Branch, config.Path)
	}

	return &config, nil
}

func cloneRepo(repo string, branch string, path string) {
	if branch == "" {
		branch = "main"
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           repo,
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil {
		log.Println("Error:", err)
	}
}
