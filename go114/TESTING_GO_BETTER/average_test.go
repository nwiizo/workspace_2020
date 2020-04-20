package main

import "testing"

func TestAverage(t *testing.T) {
	result := average([]float64{7.0, 5.0, 9.0})
	if result != 7 {
		t.Errorf("average was incorrect, got: %f, want: %d.", result, 7)
	}
}

func TestAverage2(t *testing.T) {
	result := average([]float64{7.5, 5.5, 9.5})
	if result != 7.5 {
		t.Errorf("average was incorrect, got: %f, want: %f.", result, 7.5)
	}
}
