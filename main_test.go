package main

import (
	"fmt"
	"testing"
)

func TestFastFloatToString(t *testing.T) {
	data := []float64{0, -10.565884654, -1.0001, -0.999999, -0.32, 0, 0.00001, 0.5654, 1.0001, 3}
	for _, v := range data {
		got := fastFloatToString(v)
		expected := fmt.Sprintf("%.2f", v)
		if expected != got {
			t.Errorf("%v: expected=%v, got=%v", v, expected, got)
		}

	}
}
