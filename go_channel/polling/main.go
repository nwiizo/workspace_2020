package main

import (
	"fmt"
	"time"
)

func main() {

	q := make(chan struct{}, 2000)

	go func() {
		// 重たい処理
		time.Sleep(5 * time.Second)
		q <- struct{}{}
	}()

	for {
		fmt.Println(len(q))
		if len(q) > 0 {
			break
		}
		// q に溜まるまで他の事をしたい
		time.Sleep(1 * time.Second)
		fmt.Println("Do it!!!")
	}
}
