package main

import "fmt"

func main() {
	var out []*int
	for i := 0; i < 3; i++ {
		out = append(out, &i)
	}
	fmt.Println("Values:", *out[0], *out[1], *out[2])
	fmt.Println("Addresses:", out[0], out[1], out[2])

	var out1 []*int
	for i1 := 0; i1 < 3; i1++ {
		i1 := i1
		out1 = append(out1, &i1)
	}
	fmt.Println("Values:", *out1[0], *out1[1], *out1[2])
	fmt.Println("Addresses:", out1[0], out1[1], out1[2])
}
