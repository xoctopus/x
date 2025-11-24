package retry_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/retry"
	. "github.com/xoctopus/x/testx"
)

func TestRetry_Do(t *testing.T) {
	r := &retry.Retry{}
	r.SetDefault()
	Expect(t, r.Interval, Equal(3*time.Second))
	Expect(t, r.Repeats, Equal(3))

	times := 0
	exec := func() error {
		times++
		if times == 3 {
			return nil
		}
		return errors.Errorf("times %d", times)
	}

	Expect(t, retry.Do(r, exec), Succeed())

	times = 0
	r.Repeats = 0
	Expect(t, retry.Do(r, exec), Failed())
}
