package main

import (
	"fmt"
	"sync"
	"time"
)

func doHOST(wg *sync.WaitGroup, q chan string) {

	defer wg.Done()
	for {
		host, ok := <-q
		if !ok {
			return
		}
		fmt.Println("ssh: ", host)
		time.Sleep(3 * time.Second)
	}
}
func main() {
	var wg sync.WaitGroup
	q := make(chan string, 10)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doHOST(&wg, q)
	}

	q <- "192.168.10.1"
	q <- "192.168.10.2"
	q <- "192.168.10.3"
	q <- "192.168.10.4"
	q <- "192.168.10.5"
	close(q)
	wg.Wait()
}
