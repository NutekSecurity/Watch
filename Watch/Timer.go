package Watch

import (
	"time"
)

type Timer struct {
	EndTime time.Time
	Running bool
	Ended   bool
}

func NewTimer() Timer {
	return Timer{
		EndTime: time.Now(),
		Ended:   false,
		Running: false,
	}
}

func (t *Timer) SetEndTime(hours int, minutes int, seconds int) {
	t.EndTime = time.Now()
	duration := time.Duration(hours*int(time.Hour) + minutes*int(time.Minute) + seconds*int(time.Second))
	t.EndTime = t.EndTime.Add(duration)
	t.Running = true
}

func (t *Timer) Stop() {
	t = &Timer{
		EndTime: time.Now(),
		Ended:   false,
		Running: false,
	}
}

func (t *Timer) Countdown() string {
	countdown := time.Until(t.EndTime)
	if countdown <= 0 {
		t.Running = true
		t.Ended = true
	}
	return countdown.String()
	// hours := countdown / 1e9 * 3600
	// countdown = countdown - 1e9 * hours * 3600
	// minutes := countdown / 1e9 * 60
	// countdown = countdown - 1e9 * minutes * 60
	// seconds = countdown / 1e9
	// return hours, minutes, seconds
}
