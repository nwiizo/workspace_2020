package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start main\n")
	go sayHello()
	time.Sleep(time.Second * 1)
	fmt.Println("say hello from main goroutine")

}

func sayHello() {
	fmt.Println("say hello")
}
