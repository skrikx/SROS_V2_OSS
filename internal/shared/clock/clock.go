package clock

import "time"

// Clock is the minimum neutral time provider used by downstream workflows.
type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
