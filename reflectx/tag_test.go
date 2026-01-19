package reflectx_test

import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/reflectx"
	. "github.com/xoctopus/x/testx"
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
				err:  codex.New(reflectx.ECODE__INVALID_FLAG_VALUE),
			},
			{
				name: "InvalidFlagName",
				tag:  reflect.StructTag(`json:"x y"`),
				err:  codex.New(reflectx.ECODE__INVALID_FLAG_NAME),
			},
			{
				name: "EscapeFlagValue",
				tag:  reflect.StructTag(`escape_js\non:""`),
				err:  codex.New(reflectx.ECODE__INVALID_FLAG_NAME),
			},
			{
				name: "UnquotedOption",
				tag:  reflect.StructTag(`panic_unquoted:",a='b"`),
				err:  codex.New(reflectx.ECODE__INVALID_OPTION_UNQUOTED),
			},
			{
				name: "InvalidOptionKey",
				tag:  reflect.StructTag(`panic_invalid_key:"key,'x\n\r'='any'"`),
				err:  codex.New(reflectx.ECODE__INVALID_OPTION_KEY),
			},
			// {
			// 	name: "InvalidOptionValue",
			// 	tag:  reflect.StructTag(`panic_invalid_value:",x=a b c"`),
			// 	err:  codex.New(ECODE__INVALID_OPTION_VALUE),
			// },
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				if c.err != nil {
					ExpectPanic(t, func() { reflectx.ParseTag(c.tag) }, IsError(c.err))
				} else {
					Expect(t, reflectx.ParseTag(c.tag), HaveLen[reflectx.Tag](0))
				}
			})
		}
	})

	t.Run("FlagDuplicated", func(t *testing.T) {
		tag := reflectx.ParseTag(`json:"conflict" json:"ignored"`)
		Expect(t, tag, HaveLen[reflectx.Tag](1))
		Expect(t, tag.Get("json").Value(), Equal(`"conflict"`))
	})

	t.Run("Success", func(t *testing.T) {
		cases := map[string]struct {
			tag      string
			key      string
			name     string
			value    string
			prettied string
			options  map[string]*reflectx.Option
		}{
			"OptionValueContainsSpecialChar": {
				tag:      `tag:" a, x = '\"\\?#=%,;{}[] ' "`,
				key:      "tag",
				name:     "a",
				value:    `"a,x='\"\\?#=%,;{}[] '"`,
				prettied: `tag:"a,x='\"\\?#=%,;{}[] '"`,
				options: map[string]*reflectx.Option{
					"x": reflectx.NewOption("x", `'"\?#=%,;{}[] '`, 0),
				},
			},
			"MultiOptions": {
				tag:      `tag:"b , x , y = '15.0,10.0' "`,
				key:      "tag",
				name:     "b",
				value:    `"b,x,y='15.0,10.0'"`,
				prettied: `tag:"b,x,y='15.0,10.0'"`,
				options: map[string]*reflectx.Option{
					"x": reflectx.NewOption("x", "", 0),
					"y": reflectx.NewOption("y", "'15.0,10.0'", 1),
				},
			},
			"EmptyFlagValue": {
				tag:      `tag:"  , ,   "`,
				key:      "tag",
				name:     "",
				value:    `""`,
				prettied: `tag:""`,
				options:  map[string]*reflectx.Option{},
			},
			"NeedHandleOptionKeyValueQuotes": {
				tag:      `tag:",'opt'=xyz"`,
				key:      "tag",
				name:     "",
				value:    `",opt=xyz"`,
				prettied: `tag:",opt=xyz"`,
				options: map[string]*reflectx.Option{
					"opt": reflectx.NewOption("opt", "xyz", 0),
				},
			},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				tag := reflectx.ParseTag(reflect.StructTag(c.tag))
				Expect(t, tag.Get("x"), BeNil[*reflectx.Flag]())

				f := tag.Get("tag")
				Expect(t, f.Key(), Equal("tag"))
				Expect(t, f.Name(), Equal(c.name))
				Expect(t, f.Value(), Equal(c.value))
				Expect(t, f.String(), Equal(c.prettied))
				Expect(t, f.OptionLen(), Equal(len(c.options)))
				for o := range f.Options() {
					k := o.Key()
					Expect(t, f.HasOption(k), BeTrue())
					Expect(t, f.Option(k), Equal(o))
					Expect(t, o.Value(), Equal(c.options[k].Value()))
					Expect(t, o.Quoted(), Equal(c.options[k].Quoted()))
					Expect(t, o.Unquoted(), Equal(c.options[k].Unquoted()))
				}
			})
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		options := []*reflectx.Option{
			reflectx.NewOption("", "", 0),
			reflectx.NewOption("", "has", 0),
		}
		for _, option := range options {
			Expect(t, option.IsZero(), BeTrue())
			Expect(t, option.Key(), Equal(""))
			Expect(t, option.Value(), Equal(""))
			Expect(t, option.Quoted(), Equal(""))
			Expect(t, option.Unquoted(), Equal(""))
			Expect(t, option.String(), Equal(""))
		}
	})

	t.Run("External", func(t *testing.T) {
		f := reflectx.ParseTag(`db:"f_col,default=CURRENT_TIMESTAMP(3),onupdate=CURRENT_TIMESTAMP(3)"`).Get("db")
		Expect(t, f.Key(), Equal("db"))
		Expect(t, f.Name(), Equal("f_col"))
		Expect(t, f.Option("default").String(), Equal("default=CURRENT_TIMESTAMP(3)"))
		Expect(t, f.Option("default").Value(), Equal("CURRENT_TIMESTAMP(3)"))
	})
}
