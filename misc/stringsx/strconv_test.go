package stringsx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/misc/stringsx"
)

func TestAtoi(t *testing.T) {
	i, err := Atoi("0b11")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(i).To(Equal(3))

	i, err = Atoi("011")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(i).To(Equal(9))

	i, err = Atoi("0x11")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(i).To(Equal(17))

	i, err = Atoi("11")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(i).To(Equal(11))
}
