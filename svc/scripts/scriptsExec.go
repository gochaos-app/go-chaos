package scripts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecuteScript(region string, project string, script string, namespace string, tag string, chaos string, number int, dry bool) {
	if strings.Count(script, ":") != 1 {
		log.Println("script must be in the format 'execution binary: path to script' ")
		return
	}
	parts := strings.Split(script, ":")
	executor := parts[0]
	script = parts[1]
	if dry {
		log.Println("Dry run, not executing script")
	} else {
		log.Println("Executing script", script)
		os.Setenv("REGION", region)
		os.Setenv("PROJECT", project)
		os.Setenv("NAMESPACE", namespace)
		os.Setenv("TAG", tag)
		os.Setenv("CHAOS", chaos)
		os.Setenv("NUMBER", fmt.Sprintf("%d", number))

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
