package testx

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func Recover(v any) string {
	if v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func AssertRecoverEqual(t *testing.T, caught any, expect string) {
	NewWithT(t).Expect(Recover(caught)).To(Equal(expect))
}

func AssertRecoverContains(t *testing.T, caught any, expect string) {
	NewWithT(t).Expect(Recover(caught)).To(ContainSubstring(expect))
}
