package main

import "testing"

func TestAverage(t *testing.T) {
	tests := []struct {
		name string
		input float64
		exresult float64
	}{
		{
			name "nomal vuluea",
			input []float64{7.0, 5.0, 9.0},
			exresult 7.0
		}
	}

	result := average([]float64{7.0, 5.0, 9.0})
	if result != 7 {
		t.Errorf("average was incorrect, got: %f, want: %d.", result, 7)
	}
}
