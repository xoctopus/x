package stringsx_test

import (
	"testing"

	. "github.com/xoctopus/x/stringsx"
	. "github.com/xoctopus/x/testx"
)

func TestAtoi(t *testing.T) {
	i, err := Atoi("0b11")
	Expect(t, err, Succeed())
	Expect(t, i, Equal(3))

	i, err = Atoi("011")
	Expect(t, err, Succeed())
	Expect(t, i, Equal(9))

	i, err = Atoi("0x11")
	Expect(t, err, Succeed())
	Expect(t, i, Equal(17))

	i, err = Atoi("11")
	Expect(t, err, Succeed())
	Expect(t, i, Equal(11))

	i, err = Atoi("0")
	Expect(t, err, Succeed())
	Expect(t, i, Equal(0))
}
