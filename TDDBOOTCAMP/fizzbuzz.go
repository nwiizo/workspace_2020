package main

import (
	"fmt"
	"strconv"
)

func main() {
	for n := 1; n < 100+1; n++ {
		fmt.Println(Fizzbuzz(n))
	}
}

func Fizzbuzz(i int) string {
	if i%5 == 0 && i%3 == 0 {
		return "Fizzbuzz"
	}
	if i%3 == 0 {
		return "Fizz"
	}

	if i%5 == 0 {
		return "buzz"
	}
	s := strconv.Itoa(i)
	return s
}
