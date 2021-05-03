package Watch

import "time"

type Stopwatch struct {
	startTime time.Time
	endTime   time.Duration
	// ticker    *time.Ticker
	Running bool
}

func NewStopwatch() Stopwatch {
	stopwatch := Stopwatch{
		time.Now(),
		-1,
		false,
	}
	return stopwatch
}

// Start the
// stopwatch
func (v *Stopwatch) Start() {
	v.startTime = time.Now()
	v.Running = true
}

// Stop return time elapsed since start
func (v *Stopwatch) Stop() time.Duration {
	// now := <-v.ticker.C
	// _ = <-v.ticker.C
	// v.ticker.Stop()
	v.endTime = time.Since(v.startTime)
	v.Running = false
	return v.endTime
}

// Diff gives difference between start and now or end
// of measured portion of time
func (v *Stopwatch) Diff() time.Duration {
	if v.endTime == -1 && !v.Running {
		return 0
	} else if v.endTime == -1 && v.Running {
		return time.Since(v.startTime)
	} else if v.endTime > -1 {
		return v.endTime
	} else {
		return -1
	}
}
