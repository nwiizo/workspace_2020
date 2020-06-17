package main

import "github.com/schollz/progressbar/v3"
import "time"

func main() {
	bar := progressbar.New(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
}
