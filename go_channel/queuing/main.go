package main

import (
	"fmt"
	"time"
)

func main() {
	q := make(chan string, 5)

	go func() {
		time.Sleep(3 * time.Second)
		q <- "zun"
	}()
	go func() {
		time.Sleep(3 * time.Second)
		q <- "doko"
	}()

	var cmds []string
	cmds = append(cmds, <-q)
wait_some:
	for {
		select {
		case cmd := <-q:
			cmds = append(cmds, cmd)
		case <-time.After(1 * time.Second):
			break wait_some
		}
	}
	for _, cmd := range cmds {
		fmt.Println(cmd)
	}
}
