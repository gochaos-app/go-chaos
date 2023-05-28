package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	time_test := fmt.Sprintf("Get request took %.3f seconds\n", seconds)
	test := fmt.Sprintf("Info: %s\n", string(body))

	file, _ := os.OpenFile(report, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("could not create file", err)
		return
	}
	defer file.Close()
	content := time_test + test
	_, err = file.WriteString(content)
	if err != nil {
		log.Println("Could not write to file")
		return
	}

}
