package restc

import "time"

func minDuration(duration1, duration2 time.Duration) time.Duration {
	if duration1 < duration2 {
		return duration1
	}
	return duration2
}
