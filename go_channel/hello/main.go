package main

import (
	"fmt"
	"time"
)

func main() {
	go sayHello()
	time.Sleep(time.Second * 1)
	fmt.Println("say hello from main goroutine")

}

func sayHello() {
	fmt.Println("say hello")
}
