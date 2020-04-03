package main

import "fmt"

func main() {
	for i := 1; i <= 10; i++ {
		func() {
			fmt.Println(i)
		}()
	}

	values := []int{2, 3, 4}
	for val := range values {
		go func(val interface{}) {
			fmt.Println(val)
		}(val)
	}

	for i := range values {
		val := values[i]
		go func() {
			fmt.Println(val)
		}()
	}
}
