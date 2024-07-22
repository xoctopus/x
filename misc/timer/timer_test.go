package timer_test

import (
	"testing"
	"time"

	"github.com/xoctopus/x/misc/timer"
)

func TestSpan(t *testing.T) {
	cost := timer.Span()
	time.Sleep(time.Second)
	t.Log(cost())
}
