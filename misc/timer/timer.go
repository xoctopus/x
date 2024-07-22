package timer

import "time"

func Span() func() time.Duration {
	t := time.Now()
	return func() time.Duration {
		return time.Since(t)
	}
}
