package vttgen

import (
	"testing"
)

func TestTimeToSeconds(t *testing.T) {
	time := "01:03:30"

	out := timeToSeconds(time)
	if out != 3810 {
		t.Error("Expected 3810, got ", out)
	}
}

func TestTbr(t *testing.T) {
	in := "25 tbr"

	out := tbr(in)
	if out != 25 {
		t.Error("Expected 25, got ", out)
	}
}

func TestDuration(t *testing.T) {
	in := "Duration: 00:01:06.20, start: 0.000000, bitrate: 1272 kb/s"

	out := duration(in)
	if out != 66 {
		t.Error("Expected 66, got ", out)
	}
}
