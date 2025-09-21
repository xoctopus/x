package reflectx_test

import (
	"bytes"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/reflectx"
)

func TestParseFlags(t *testing.T) {
	t.Run("InvalidTag", func(t *testing.T) {
		cases := []struct {
			name string
			tag  reflect.StructTag
			err  error
		}{
			{
				name: "EmptyTag",
				tag:  reflect.StructTag("   "),
			},
			{
				name: "NoFlagKey",
				tag:  reflect.StructTag(`:`),
			},
			{
				name: "UnquotedFlagValue",
				tag:  reflect.StructTag(`any:"x`),
				err:  ErrInvalidFlagValue,
			},
			{
				name: "InvalidFlagName",
				tag:  reflect.StructTag(`json:"x y"`),
				err:  ErrInvalidFlagName,
			},
			{
				name: "EscapeFlagValue",
				tag:  reflect.StructTag(`escape_js\non:""`),
				err:  ErrInvalidFlagKey,
			},
			{
				name: "UnquotedOption",
				tag:  reflect.StructTag(`panic_unquoted:",a='b"`),
				err:  ErrInvalidOptionUnquoted,
			},
			{
				name: "InvalidOptionKey",
				tag:  reflect.StructTag(`panic_invalid_key:"key,'x\n\r'='any'"`),
				err:  ErrInvalidOptionKey,
			},
			{
				name: "InvalidOptionValue",
				tag:  reflect.StructTag(`panic_invalid_value:",x=a b c"`),
				err:  ErrInvalidOptionValue,
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				if c.err != nil {
					defer func() {
						err := recover()
						NewWithT(t).Expect(err).To(Equal(c.err))
					}()
				}
				tag := ParseTag(c.tag)
				NewWithT(t).Expect(tag).To(HaveLen(0))
			})
		}
	})

	t.Run("FlagDuplicated", func(t *testing.T) {
		tag := ParseTag(`json:"conflict" json:"ignored"`)
		NewWithT(t).Expect(tag).To(HaveLen(1))
		NewWithT(t).Expect(tag.Get("json").Value()).To(Equal(`"conflict"`))
	})

	t.Run("Success", func(t *testing.T) {
		cases := map[string]struct {
			tag      string
			key      string
			name     string
			quoted   string
			unquoted string
			value    string
			prettied string
			options  map[string]*Option
		}{
			"OptionValueContainsSpecialChar": {
				tag:      `tag:" a, x = '\"\\?#=%,;{}[] ' "`,
				key:      "tag",
				name:     "a",
				quoted:   `" a, x = '\"\\?#=%,;{}[] ' "`,
				unquoted: ` a, x = '"\?#=%,;{}[] ' `,
				value:    `"a,x='\"\\?#=%,;{}[] '"`,
				prettied: `tag:"a,x='\"\\?#=%,;{}[] '"`,
				options: map[string]*Option{
					"x": NewOption("x", `'"\?#=%,;{}[] '`, 0),
				},
			},
			"MultiOptions": {
				tag:      `tag:"b , x , y = '15.0,10.0' "`,
				key:      "tag",
				name:     "b",
				quoted:   `"b , x , y = '15.0,10.0' "`,
				unquoted: `b , x , y = '15.0,10.0' `,
				value:    `"b,x,y='15.0,10.0'"`,
				prettied: `tag:"b,x,y='15.0,10.0'"`,
				options: map[string]*Option{
					"x": NewOption("x", "", 0),
					"y": NewOption("y", "'15.0,10.0'", 1),
				},
			},
			"EmptyFlagValue": {
				tag:      `tag:"  , ,   "`,
				key:      "tag",
				name:     "",
				quoted:   `"  , ,   "`,
				unquoted: `  , ,   `,
				value:    `""`,
				prettied: `tag:""`,
				options:  map[string]*Option{},
			},
			"NeedHandleOptionKeyValueQuotes": {
				tag:      `tag:",'opt'=xyz"`,
				key:      "tag",
				name:     "",
				quoted:   `",'opt'=xyz"`,
				unquoted: `,'opt'=xyz`,
				value:    `",opt='xyz'"`,
				prettied: `tag:",opt='xyz'"`,
				options: map[string]*Option{
					"opt": NewOption("opt", "'xyz'", 0),
				},
			},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				tag := ParseTag(reflect.StructTag(c.tag))
				NewWithT(t).Expect(tag.Get("x")).To(BeNil())

				f := tag.Get("tag")
				NewWithT(t).Expect(f.Key()).To(Equal("tag"))
				NewWithT(t).Expect(f.Name()).To(Equal(c.name))
				NewWithT(t).Expect(f.QuotedValue()).To(Equal(c.quoted))
				NewWithT(t).Expect(f.UnquotedValue()).To(Equal(c.unquoted))
				NewWithT(t).Expect(f.Value()).To(Equal(c.value))
				NewWithT(t).Expect(f.String()).To(Equal(c.prettied))
				NewWithT(t).Expect(f.OptionLen()).To(Equal(len(c.options)))
				for k, v := range c.options {
					o := f.Option(k)
					NewWithT(t).Expect(o.Value()).To(Equal(v.Value()))
					NewWithT(t).Expect(bytes.Equal(o.RawValue(), v.RawValue())).To(BeTrue())
				}
			})
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		options := []*Option{
			NewOption("", "", 0),
			NewOption("", "has", 0),
		}
		for _, option := range options {
			NewWithT(t).Expect(option.Key()).To(Equal(""))
			NewWithT(t).Expect(option.Value()).To(Equal(""))
			NewWithT(t).Expect(option.IsZero()).To(BeTrue())
			NewWithT(t).Expect(option.String()).To(Equal(""))
		}
	})
}
