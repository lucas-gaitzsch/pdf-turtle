package utils

import (
	"fmt"
	"testing"
)

var stringsToTrim = []struct {
	in  string
	out string
}{
	{" test ", "test"},
	{" test \n ", "test"},
	{" test \t ", "test"},
	{" \t\n test ", "test"},
	{" \t\n test \t test ", "test \t test"},
}

func TestTrim(t *testing.T) {
	for _, d := range stringsToTrim {
		t.Run(fmt.Sprintf("trim %s", d.in), func(t *testing.T) {
			trimmed := Trim(d.in)

			if trimmed != d.out {
				t.Fatalf("%s trimmed should be %s (current: %s)", d.in, d.out, trimmed)
			}
		})
	}
}
