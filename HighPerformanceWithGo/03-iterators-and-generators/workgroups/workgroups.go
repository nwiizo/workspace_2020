package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Test struct {
	Url string
}

func ReadOnStruct(fileBuffer []byte) ([]Test, error) {
	data := make([]Test, 20)
	err := yaml.Unmarshal(fileBuffer, &data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func retrieve(url string, wg *sync.WaitGroup) {
	// WaitGroup Counter-- when Goroutine is finished
	defer wg.Done()
	start := time.Now()
	res, err := http.Get(url)
	end := time.Since(start)
	if err != nil {
		panic(err)
	}
	// Print the status code from the response
	fmt.Println(url, res.StatusCode, end)

}

func main() {
	buf, err := ioutil.ReadFile("./test.yaml")

	data, err := ReadOnStruct(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for i := range data {
		// Increment WaitGroup Counter when new Goroutine is called
		wg.Add(1)
		urls := data[i].Url
		go retrieve(urls, &wg)
	}
	// Wait for the collection of Goroutines to finish
	wg.Wait()
}
