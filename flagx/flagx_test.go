package flagx_test

import (
	"testing"

	"github.com/xoctopus/x/flagx"
	. "github.com/xoctopus/x/testx"
)

func TestFlag(t *testing.T) {
	f := flagx.NewFlag[uint8]()
	Expect(t, f.Is(0b0001), BeFalse())
	f.With(0b1111)
	Expect(t, f.Value(), Equal[uint8](0b1111))
	Expect(t, f.Is(0b1111), BeTrue())
	Expect(t, f.Is(0b0111), BeTrue())
	Expect(t, f.Is(0b0011), BeTrue())
	Expect(t, f.Is(0b0001), BeTrue())

	f.Trim(0b1001)
	Expect(t, f.Value(), Equal[uint8](0b0110))
	Expect(t, f.Is(0b0000), BeTrue())
	Expect(t, f.Is(0b0110), BeTrue())
	Expect(t, f.Is(0b0010), BeTrue())
	Expect(t, f.Is(0b0100), BeTrue())

	Expect(t, f.Is(0b1100), BeFalse())
	Expect(t, f.Is(0b0001), BeFalse())
}
