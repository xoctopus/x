package retry_test

import (
	"testing"

	g "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/retry"
)

func TestRetry_Do(t *testing.T) {
	r := &retry.Retry{}
	r.SetDefault()
	g.NewWithT(t).Expect(r.Interval).To(g.Equal(retry.Default.Interval))
	g.NewWithT(t).Expect(r.Repeats).To(g.Equal(retry.Default.Repeats))

	times := 0
	exec := func() error {
		times++
		if times == 3 {
			return nil
		}
		return errors.Errorf("times %d", times)
	}

	g.NewWithT(t).Expect(retry.Do(r, exec)).To(g.BeNil())

	times = 0
	r.Repeats = 0
	g.NewWithT(t).Expect(retry.Do(r, exec)).NotTo(g.BeNil())
}
