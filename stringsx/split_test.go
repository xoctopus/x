package stringsx_test

import (
	"testing"

	"github.com/xoctopus/x/stringsx"
	. "github.com/xoctopus/x/testx"
)

func Test_SplitToWords(t *testing.T) {
	words := []string{"I", "Am", "A", "10", "Years", "Senior"}
	cases := []struct {
		phrase string
		words  []string
	}{
		{"IAmA10YearsSenior", words},
		{"I Am A 10 Years Senior", words},
		{". I_ Am_A_10_Years____Senior__", words},
		{"I-~~ Am\nA\t10 Years *** Senior", words},
		{"lowercase", []string{"lowercase"}},
		{"Class", []string{"Class"}},
		{"MyClass", []string{"My", "Class"}},
		{"HTML", []string{"HTML"}},
		{"QOSType", []string{"QOS", "Type"}},
		{"\xF1\xF4\x8F\xBF\xBF", []string{"\xF1\xF4\x8F\xBF\xBF"}},
	}
	for _, c := range cases {
		Expect(t, stringsx.SplitToWords(c.phrase), Equal(c.words))
	}

	t.Log(stringsx.CheckLetterType('1'))
	t.Log(stringsx.CheckLetterType('a'))
	t.Log(stringsx.CheckLetterType('A'))
	t.Log(stringsx.CheckLetterType('+'))
	t.Log(stringsx.CheckLetterType('🚀'))
}
