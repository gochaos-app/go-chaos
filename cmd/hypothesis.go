package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func Ping(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)

}

func GenerateReport() {
	fmt.Println("testing")
}
