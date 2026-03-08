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
}

func TestSplitCamelCase(t *testing.T) {
	for _, word := range [][]string{
		{""},
		{"a", "a"},
		{"A", "A"},
		{"abc", "abc"},
		{"ABC", "ABC"},
		{"HTTPServer", "HTTP", "Server"},
		{"PDFLoader", "PDF", "Loader"},
		{"UserID", "User", "ID"},
		{"UserID100", "User", "ID100"},
		{"UserV2", "User", "V2"},
		{"User100Failed", "User100", "Failed"},
		{"userV2Model", "user", "V2", "Model"},
		{"你好World", "你好", "World"},
		{"JSON解析器", "JSON", "解析器"},
		{"ΓειάΣουWorld", "Γειά", "Σου", "World"},
		{"golang👍", "golang", "👍"},
		{string([]byte{0xff, 0xfe}), string([]byte{0xff, 0xfe})},
	} {
		t.Run(word[0], func(t *testing.T) {
			words := stringsx.SplitCamelCase(word[0])
			Expect(t, words, Equal(word[1:]))
		})
	}

}
