package docx_test

import (
	"testing"

	"github.com/xoctopus/x/docx"
	. "github.com/xoctopus/x/testx"
)

type T struct{}

func (*T) DocOf(names ...string) ([]string, bool) {
	if len(names) == 0 {
		return []string{"can be doc"}, true
	}
	name := names[0]
	if name == "a" {
		return []string{"field a"}, true
	}
	return []string{}, false
}

func TestDocOf(t *testing.T) {
	d, ok := docx.Of(&T{}, "", "a")
	Expect(t, d, Equal([]string{"field a"}))
	Expect(t, ok, BeTrue())

	d, ok = docx.Of(T{}, "", "a")
	Expect(t, d, Equal([]string{}))
	Expect(t, ok, BeFalse())

	d, ok = docx.Of(&T{}, "p ", "a")
	Expect(t, d, Equal([]string{"p field a"}))
	Expect(t, ok, BeTrue())
}
