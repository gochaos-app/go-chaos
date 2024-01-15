package scripts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gochaos-app/go-chaos/config"
)

func ExecuteScript(cfg config.JobConfig, dry bool) {
	if strings.Count(cfg.Service, ":") != 1 {
		log.Println("script must be in the format 'execution binary: path to script' ")
		return
	}
	parts := strings.Split(cfg.Service, ":")
	executor := parts[0]
	script := parts[1]
	if dry {
		log.Println("Dry run, not executing script")
	} else {
		log.Println("Executing script", script)
		os.Setenv("REGION", cfg.Region)
		os.Setenv("PROJECT", cfg.Project)
		os.Setenv("NAMESPACE", cfg.Namespace)
		os.Setenv("TAG", cfg.Chaos.Tag)
		os.Setenv("CHAOS", cfg.Chaos.Chaos)
		os.Setenv("NUMBER", fmt.Sprintf("%d", cfg.Chaos.Count))

		cmd, err := exec.Command(executor, script).Output()

		if err != nil {
			log.Fatal(err)
		}

		output := string(cmd)
		log.Println(output)

	}
}

/*
func ExecuteScript(script string, executor string) {
	log.Println("Executing script:", script)
	comand, err := exec.Command(executor, script).Output()
	if err != nil {
		log.Println("error", err)
	}
	output := string(comand)
	log.Println(output)
}
*/
