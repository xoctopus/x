package reflectx_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/reflectx"
)

func TestParseFlags(t *testing.T) {
	t.Run("InvalidTag", func(t *testing.T) {
		v := reflect.TypeOf(struct {
			EmptyTag          any `   `
			NoFlagKey         any `:`
			UnquotedFlagValue any `unquoted:"key`
			EscapeFlagKey     any `escape_js\non:""`
		}{})

		for i := range v.NumField() {
			f := v.Field(i)
			t.Run(f.Name, func(t *testing.T) {
				if strings.HasPrefix(string(f.Tag), "unquoted") {
					defer func() {
						e := recover()
						NewWithT(t).Expect(e).To(Equal(ErrInvalidFlagRaw))
					}()
				}
				if strings.HasPrefix(string(f.Tag), "escape") {
					defer func() {
						e := recover()
						NewWithT(t).Expect(e).To(Equal(ErrInvalidFlagKey))
					}()
				}
				tag := ParseTag(f.Tag)
				NewWithT(t).Expect(tag).To(HaveLen(0))
			})
		}
	})
	t.Run("FlagDuplicated", func(t *testing.T) {
		tag := ParseTag(`json:"conflict" json:"ignored"`)
		NewWithT(t).Expect(tag).To(HaveLen(1))
		NewWithT(t).Expect(tag.Get("json").Value()).To(Equal("conflict"))
	})

	t.Run("InvalidOption", func(t *testing.T) {
		v := reflect.TypeOf(struct {
			UnquotedOption   any `unquoted:",a='b"`
			InvalidOptionKey any `invalid:"key,'x\n\r'='any'"`
		}{})
		for i := range v.NumField() {
			f := v.Field(i)
			t.Run(f.Name, func(t *testing.T) {
				if strings.HasPrefix(string(f.Tag), "unquoted") {
					defer func() {
						e := recover()
						NewWithT(t).Expect(e).To(Equal(ErrInvalidOptionUnquoted))
					}()
				}
				if strings.HasPrefix(string(f.Tag), "invalid") {
					defer func() {
						e := recover()
						NewWithT(t).Expect(e).To(Equal(ErrInvalidOptionKey))
					}()
				}
				NewWithT(t).Expect(ParseTag(f.Tag)).To(HaveLen(0))
			})
		}
	})

	t.Run("Success", func(t *testing.T) {
		cases := []struct {
			tag     string
			name    string
			options map[string]*Option
			pretty  string
		}{
			{
				tag:    `tag:" a, x = '\"\\?#=%,;{}[] ' "`,
				pretty: `tag:"a,x='\"\\?#=%,;{}[] '"`,
				name:   "a",
				options: map[string]*Option{
					"x": NewOption("x", `'\"\\?#=%,;{}[] '`, 0),
				},
			},
			{
				tag:    `tag:"b , x , y = '15.0,10.0' "`,
				pretty: `tag:"b,x,y='15.0,10.0'"`,
				name:   "b",
				options: map[string]*Option{
					"y": NewOption("y", `'15.0,10.0'`, 0),
				},
			},
			{
				tag:    `tag:" c , "`,
				pretty: `tag:"c"`,
				name:   "c",
			},
			{
				tag:    `tag:"d"`,
				pretty: `tag:"d"`,
				name:   "d",
			},
			{
				tag:    `tag:", , y = 'abc' x = 1.1 "`,
				pretty: `tag:",y='abc',x='1.1'"`,
				options: map[string]*Option{
					"y": NewOption("y", `'abc'`, 0),
					"x": NewOption("x", `'1.1'`, 1),
				},
			},
			{
				tag:    `tag:", 'xyz' = 'abc'"`,
				pretty: `tag:""`,
			},
			{
				tag:    `tag:", = "`,
				pretty: `tag:""`,
			},
			{
				tag:    `tag:", y ='', , x"`,
				pretty: `tag:",y='',x"`,
				options: map[string]*Option{
					"x": NewOption("x", ``, 1),
					"y": NewOption("y", `''`, 0),
				},
			},
		}

		for _, c := range cases {
			flag := ParseTag(reflect.StructTag(c.tag)).Get("tag")

			NewWithT(t).Expect(flag.Key()).To(Equal("tag"))
			NewWithT(t).Expect(flag.Name()).To(Equal(c.name))
			NewWithT(t).Expect(flag.OptionLen()).To(Equal(len(c.options)))
			for k, v := range c.options {
				o := flag.Option(k)
				NewWithT(t).Expect(o.Value()).To(Equal(v))
			}
			raw, _ := strconv.Unquote(strings.TrimPrefix(c.tag, "tag:"))
			NewWithT(t).Expect(flag.Raw()).To(Equal(raw))
			pretty := strings.TrimPrefix(c.pretty, "tag:")
			NewWithT(t).Expect(flag.Value()).To(Equal(pretty))
			NewWithT(t).Expect(flag.String()).To(Equal(c.pretty))
		}
	})
}
