package internal_test

import (
	"go/types"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/internal"
)

func TestNamedBacktrace(t *testing.T) {
	tt1 := NewPackage("net").Scope().Lookup("Addr").Type().(*types.Named)
	tt2 := NewPackage("io").Scope().Lookup("Reader").Type().(*types.Named)

	names := NamedBacktrace{}
	NewWithT(t).Expect(names.Len()).To(Equal(0))

	top := names.Top()
	NewWithT(t).Expect(top).To(BeNil())

	pushed := names.Push(tt1)
	NewWithT(t).Expect(pushed).To(BeTrue())
	pushed = names.Push(tt1)
	NewWithT(t).Expect(pushed).To(BeFalse())

	top = names.Top()
	NewWithT(t).Expect(types.Identical(top, tt1)).To(BeTrue())

	pushed = names.Push(tt2)
	NewWithT(t).Expect(pushed).To(BeTrue())
	pushed = names.Push(tt1)
	NewWithT(t).Expect(pushed).To(BeFalse())

	top = names.Top()
	NewWithT(t).Expect(types.Identical(top, tt2)).To(BeTrue())

	NewWithT(t).Expect(names.Len()).To(Equal(2))
}

func TestChanDirPrefix(t *testing.T) {
	t.Run("ShouldPanic", func(t *testing.T) {
		defer func() {
			v := recover()
			t.Logf("%T %v", v, v)
		}()
		ChanDirPrefix(int(types.RecvOnly))
	})
}
