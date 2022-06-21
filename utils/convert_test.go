package utils

import (
	"fmt"
	"math"
	"testing"
)

var mmToInchesTestData = []struct {
	mm     int
	inches float64
}{
	{0, 0.0},
	{1, 0.03937},
	{900, 35.4331},
}

func TestMmToInches(t *testing.T) {
	for _, d := range mmToInchesTestData {
		t.Run(fmt.Sprintf("convert %d mm to inches", d.mm), func(t *testing.T) {
			inches := MmToInches(d.mm)

			if math.Abs(inches-d.inches) > 0.0001 {
				t.Fatalf("%d mm should be %v inches (current: %v)", d.mm, d.inches, inches)
			}
		})
	}
}
