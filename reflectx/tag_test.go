package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/reflectx"
)

func TestParseFlags(t *testing.T) {
	for name, c := range map[string]struct {
		tag          reflect.StructTag
		flags        Flags
		nameTagValue string
	}{
		"InvalidTag/Empty":     {`   `, Flags{}, ""},
		"InvalidTag/NoTagName": {` : `, Flags{}, ""},
		"InvalidTag/Unquoted":  {`tag:"key`, Flags{}, ""},
		"InvalidTag/Conflict":  {`tag:"1" tag:"2"`, Flags{}, ""},
		"ContainsEscapesInOption": {
			`name:"b ,             default='\"\\?#=%,;{}[]\" ' "`,
			Flags{"name": {Name: "b", Options: [][2]string{{"default", `"\?#=%,;{}[]" `}}}},
			`b,default='"\?#=%,;{}[]" '`,
		},
		"MultiOptions": {
			`name:"c , omitempty , default='15.0,10.0'"`,
			Flags{"name": {Name: "c", Options: [][2]string{{"default", "15.0,10.0"}, {"omitempty", ""}}}},
			`c,default='15.0,10.0',omitempty`,
		},
		"EmptyOptions": {
			`name:"  ,, default='abc'"`,
			Flags{"name": {Name: "", Options: [][2]string{{"default", "abc"}}}},
			`,default='abc'`,
		},
		"NoOptions": {
			`name:"e"`,
			Flags{"name": {Name: "e"}},
			`e`,
		},
		"NoOptionKey": {
			`name:"f , = no_name_skipped_option"`,
			Flags{"name": {Name: "f"}},
			`f`,
		},
		"NoOptionValue": {
			`name:" g , default="`,
			Flags{"name": {Name: "g", Options: [][2]string{{"default", ""}}}},
			`g,default`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			NewWithT(t).Expect(ParseFlags(c.tag)).To(Equal(c.flags))
			if flag := c.flags.Get("name"); flag != nil {
				NewWithT(t).Expect(flag.TagValue()).To(Equal(c.nameTagValue))
			}
		})
	}

	flags := Flags{"tagName": {}}
	NewWithT(t).Expect(flags.Get("tagName")).NotTo(BeNil())
	NewWithT(t).Expect(flags.Get("json")).To(BeNil())
}
