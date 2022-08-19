package scripts

import (
	"log"
	"os/exec"
)

func ExecuteScript(script string) {
	log.Println("Executing script:", script)
	comand, err := exec.Command("/bin/sh", script).Output()
	if err != nil {
		log.Println("error", err)
	}
	output := string(comand)
	log.Println(output)
}
