package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gochaos-app/go-chaos/cmd"
)

func getroot(w http.ResponseWriter, r *http.Request) {
	log.Println("got / request")
	io.WriteString(w, "Welcome to Go-Chaos!\n")
	stringExecution := `To start running chaos engineering experiments, 

* use http://go-chaos-url/experiment?config=path_to_chaos_experiment.hcl for single config experiments
`
	io.WriteString(w, stringExecution)
}

func ServerFn(filename string) {
	config, err := readConfig(filename)
	if err != nil {
		log.Println("Error:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", getroot)
	mux.HandleFunc("/experiment", func(w http.ResponseWriter, r *http.Request) {
		hasConfig := r.URL.Query().Has("config")
		configFile := r.URL.Query().Get("config")

		ConfigPath := config.Path + "/" + configFile

		if hasConfig {
			if _, err := os.Stat(ConfigPath); err != nil {
				log.Println("No Chaos experiment file found", ConfigPath)
				stringNotFound := "No Chaos experiment file found" + ConfigPath + "\n"
				io.WriteString(w, stringNotFound)
			} else {
				stringFound := "Chaos experiment found" + ConfigPath + "\n"
				io.WriteString(w, stringFound)
				cfg, err := cmd.LoadConfig(ConfigPath)
				if err != nil {
					stringError := "error executing chaos:" + err.Error() + "\n"
					io.WriteString(w, stringError)
				} else {
					cmd.ExecuteChaos(cfg)
				}

			}
		}
		log.Println("got /experiment request")
	})

	var port string

	if config.Port == "" {
		port = ":3333"
	} else {
		port = ":" + config.Port
	}

	err = http.ListenAndServe(port, mux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
	}

}
