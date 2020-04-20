package main

import "fmt"

func average(values []float64) float64 {
	var total float64 = 0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func main() {
	fmt.Println(average([]float64{7.0, 5.0, 9.0}))
}
