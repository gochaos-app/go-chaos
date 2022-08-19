package scripts

import (
	"log"
	"os/exec"
)

func ExecuteScript(script string, executor string) {
	log.Println("Executing script:", script)
	comand, err := exec.Command(executor, script).Output()
	if err != nil {
		log.Println("error", err)
	}
	output := string(comand)
	log.Println(output)
}
