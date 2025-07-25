package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/xoctopus/x/reflectx"
)

func TestParseFlags(t *testing.T) {
	t.Run("InvalidTag", func(t *testing.T) {
		for name, tag := range map[string]reflect.StructTag{
			"Empty":               `   `,
			"NoTagName":           ` : `,
			"Unquoted":            `tag:"key`,
			"Conflict":            `tag:"1" tag:"2"`,
			"OptionKeyUnquoted":   `tag:",'a=b"`,
			"OptionValueUnquoted": `tag:",a='b"`,
		} {
			t.Run(name, func(t *testing.T) {
				NewWithT(t).Expect(len(ParseFlags(tag))).To(Equal(0))
			})
		}
	})

	build := func(name, value string, options ...[2]string) Flags {
		fs := make(Flags)
		f := NewFlag(value)
		for _, opt := range options {
			f.AddOption(opt[0], opt[1])
		}
		fs.Add(name, f)
		return fs
	}

	for name, c := range map[string]struct {
		tag   reflect.StructTag
		flags Flags
		str   string
	}{
		"ContainsEscape": {
			tag:   `name:" b , default = '\"\\?#=%,;{}[] ' "`,
			flags: build("name", "b", [2]string{"default", `"\?#=%,;{}[] `}),
			str:   `name:"b,default='\"\\?#=%,;{}[] '"`,
		},
		"ContainsMultiFlags": {
			tag:   `name:" c , omitempty, default = '15.0,10.0' "`,
			flags: build("name", "c", [2]string{"default", "15.0,10.0"}, [2]string{"omitempty"}),
			str:   `name:"c,default='15.0,10.0',omitempty"`,
		},
		"NoFlagName": {
			tag:   `name:"  ,  , default='abc'"`,
			flags: build("name", "", [2]string{"default", "abc"}),
			str:   `name:",default='abc'"`,
		},
		"NoOptions": {
			tag:   `name:" x ,"`,
			flags: build("name", "x"),
			str:   `name:"x"`,
		},
		"NoOptionKey": {
			tag:   `name:", = 'abc'"`,
			flags: build("name", "", [2]string{"", "abc"}),
			str:   `name:""`,
		},
		"NoOptionValue": {
			tag:   `name:",default="`,
			flags: build("name", "", [2]string{"default"}),
			str:   `name:",default"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			flags := ParseFlags(c.tag)
			NewWithT(t).Expect(flags).To(Equal(c.flags))
			NewWithT(t).Expect(flags.String()).To(Equal(c.str))
		})
	}

	// 	// "EmptyOptions": {
	// 	// 	`name:"  ,, default='abc'"`,
	// 	// 	Flags{"name": {Name: "", Options: [][2]string{{"default", "abc"}}}},
	// 	// 	`,default='abc'`,
	// 	// },
	// 	// "NoOptions": {
	// 	// 	`name:"e"`,
	// 	// 	Flags{"name": {Name: "e"}},
	// 	// 	`e`,
	// 	// },
	// 	// "NoOptionKey": {
	// 	// 	`name:"f , = no_name_skipped_option"`,
	// 	// 	Flags{"name": {Name: "f"}},
	// 	// 	`f`,
	// 	// },
	// 	// "NoOptionValue": {
	// 	// 	`name:" g , default="`,
	// 	// 	Flags{"name": {Name: "g", Options: [][2]string{{"default", ""}}}},
	// 	// 	`g,default`,
	// 	// },
	// 	// "NoNameHasOption": {
	// 	// 	`name:",required"`,
	// 	// 	Flags{"name": {Name: "", Options: [][2]string{{"required", ""}}}},
	// 	// 	`,required`,
	// 	// },
	// } {
	// 	t.Run(name, func(t *testing.T) {
	// 		NewWithT(t).Expect(ParseFlags(c.tag)).To(Equal(c.flags))
	// 		if flag := c.flags.Get("name"); flag != nil {
	// 			NewWithT(t).Expect(flag.TagValue()).To(Equal(c.nameTagValue))
	// 		}
	// 	})
	// }

}
