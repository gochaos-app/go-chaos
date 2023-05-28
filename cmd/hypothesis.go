package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

func Ping(url string, report string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()

	time_test := fmt.Sprintf("\nTOOK %.3f seconds\n", seconds)
	request := fmt.Sprintf("REQUEST:\n%s", string(reqDump))
	response := fmt.Sprintf("RESPONSE\n%s", string(respDump))

	file, _ := os.OpenFile(report, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("could not create file", err)
		return
	}
	defer file.Close()
	content := time_test + request + response
	_, err = file.WriteString(content)
	if err != nil {
		log.Println("Could not write to file")
		return
	}

}
