package testx

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

func _recover(v any) string {
	if v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func AssertRecoverEqual(t *testing.T, caught any, expect string) {
	if x := _recover(caught); len(x) > 0 {
		gomega.NewWithT(t).Expect(x).To(gomega.Equal(expect))
	}
}

func AssertRecoverContains(t *testing.T, caught any, expect string) {
	if x := _recover(caught); len(x) > 0 {
		gomega.NewWithT(t).Expect(x).To(gomega.ContainSubstring(expect))
	}
}
