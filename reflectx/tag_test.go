package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/reflectx"
)

func TestParseFlags(t *testing.T) {
	for name, c := range map[string]struct {
		tag   reflect.StructTag
		flags Flags
	}{
		"InvalidTag/Empty":     {`   `, Flags{}},
		"InvalidTag/NoTagName": {` : `, Flags{}},
		"InvalidTag/Unquoted":  {`tag:"key`, Flags{}},
		"InvalidTag/Conflict":  {`tag:"1" tag:"2"`, Flags{}},
		"ContainsEscapesInOption": {
			`name:"b ,             default='\"\\?#=%,;{}[]\" ' "`,
			Flags{
				"name": {
					Tag:     "name",
					Value:   "b",
					Options: [][2]string{{"default", `"\?#=%,;{}[]" `}},
				},
			},
		},
		"MultiOptions": {
			`name:"c , omitempty , default='15.0,10.0'"`,
			Flags{
				"name": {
					Tag:   "name",
					Value: "c",
					Options: [][2]string{
						{"default", "15.0,10.0"},
						{"omitempty", ""},
					},
				},
			},
		},
		"EmptyOptions": {
			`name:"  ,, default='abc'"`,
			Flags{
				"name": {
					Tag:   "name",
					Value: "",
					Options: [][2]string{
						{"default", "abc"},
					},
				},
			},
		},
		"NoOptions": {
			`name:"e"`,
			Flags{"name": {Tag: "name", Value: "e"}},
		},
		"NoOptionKey": {
			`name:"f , = no_name_skipped_option"`,
			Flags{"name": {Tag: "name", Value: "f"}},
		},
		"NoOptionValue": {
			`name:" g , option_key="`,
			Flags{"name": {Tag: "name", Value: "g", Options: [][2]string{{"option_key", ""}}}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			NewWithT(t).Expect(ParseFlags(c.tag)).To(Equal(c.flags))
		})
	}

	flags := Flags{"tagName": {}}
	NewWithT(t).Expect(flags.Get("tagName")).NotTo(BeNil())
	NewWithT(t).Expect(flags.Get("json")).To(BeNil())
}
