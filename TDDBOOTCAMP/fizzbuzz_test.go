package main

import (
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	result1 := Fizzbuzz(1)
	if result1 != "1" {
		t.Error("www")
	}
	result2 := Fizzbuzz(2)
	if result2 == "2" {
		t.Log("www")
	}
}
