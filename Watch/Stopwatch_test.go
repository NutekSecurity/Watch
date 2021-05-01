package Watch

import (
	"testing"
	"time"
)

func TestStopwatch(t *testing.T) {
	stopwatch := NewStopwatch()
	if stopwatch.Running {
		t.Error("Should not run after initialization")
	}
	stopwatch.Start()
	if !stopwatch.Running {
		t.Error("Not \"Running\"")
	}
	time.Sleep(1 * time.Second)
	diff1 := stopwatch.Stop()
	if stopwatch.Running {
		t.Error("Didn't stopped")
	}
	diff2 := stopwatch.Diff()
	if diff1 != diff2 {
		t.Error("Difference does not match")
	}
	if stopwatch.endTime < 1e3 {
		t.Error("Should change more", stopwatch.endTime)
	}
	if stopwatch.endTime > 1e9+10e6 {
		t.Error("Should change less", stopwatch.endTime)
	}
}
