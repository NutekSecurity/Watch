package Watch

import "time"

func UniversalTime() time.Time {
	return time.Now().UTC()
}

func LocalTime() time.Time {
	return time.Now().Local()
}

func UnixTime() int {
	return int(time.Now().UnixNano()) / 1e+6
}
