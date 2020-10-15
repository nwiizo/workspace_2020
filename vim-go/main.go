package main

import "fmt"

type Foos struct {
	Foo1 int
	Foo2 int
	Foo3 int
}

func Add(a int, b int) int {
	c := a + b
	return c
}

func main() {
	fmt.Println("vim-go")
	fmt.Println(Add(4, 3))
}
