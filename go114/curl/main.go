package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	flag.Parse()
	Url := "https://" + flag.Arg(0)
	fmt.Println(Url)
	req, err := http.NewRequest("HEAD", Url, nil)
	if err != nil {
		// handle err
	}
	for {
		resp, err := client.Do(req)
		if err != nil {
			// handle err
		}
		fmt.Println(*resp)
		resp.Body.Close()
		time.Sleep(time.Second * 1)
		defer resp.Body.Close()
	}
}
