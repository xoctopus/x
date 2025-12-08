package reflectx

import (
	"iter"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/stringsx"
)

// ParseTag parses a struct tag into a map of flag keys and values.
// Each value is further parsed into a flag name and its options.
// Control characters are allowed only in option values.
// Flag keys, flag names, and option names may contain only letters, digits, and underscores.
// Other characters in option values must be wrapped in single quotes.
func ParseTag(tag reflect.StructTag) Tag {
	flags := make(Tag)

	for i := 0; tag != ""; {
		// skip spaces
		i = 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// meet flag key
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		key := string(tag[:i])
		if !stringsx.ValidFlagName(key) {
			panic(codex.Errorf(ECODE__INVALID_FLAG_NAME, "name: %q", key))
		}

		tag = tag[i+1:]

		// meet flag value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			panic(codex.Errorf(ECODE__INVALID_FLAG_VALUE, "tag: %q", tag))
		}
		quoted := string(tag[:i+1])
		if _, ok := flags[key]; !ok {
			flags[key] = &Flag{
				key:     key,
				raw:     must.NoErrorV(strconv.Unquote(quoted)),
				options: make(map[string]*Option),
			}
			tag = tag[i+1:]
		}
	}

	for k := range flags {
		flags[k].parse()
	}

	return flags
}

type Tag map[string]*Flag

func (t Tag) Get(key string) *Flag {
	if f, ok := t[key]; ok {
		return f
	}
	return nil
}

// Flag parsed tag element
// eg: `db:"f_column,default='0',width=10,precision=4,primary"`
// the result is
//
//	{
//	  key:     "db",
//	  value:   "f_column,default='0',width=12,precision=4,primary",
//	  name:    "f_column"
//	  options: {
//	    "default":   '0',
//	    "width":     12,
//	    "precision": 4,
//	    "primary":   nil,
//	  }
//	}
type Flag struct {
	key      string
	name     string
	options  map[string]*Option
	raw      string
	value    string
	prettied string
}

func (f *Flag) parse() {
	val := strings.TrimSpace(f.raw)

	// scan value to parted by ','
	quoted := false
	parted := make([]string, 0)
	idx := 0
	for i, c := range []rune(val) {
		switch c {
		case '\'':
			quoted = !quoted
		case ',':
			if !quoted {
				goto FinishPart
			}
		}
		if i == len(val)-1 {
			if quoted {
				panic(codex.Errorf(ECODE__INVALID_OPTION_UNQUOTED, "raw: %q", val))
			}
			i++
			goto FinishPart
		}
		continue
	FinishPart:
		part := strings.TrimSpace(string([]rune(val)[idx:i]))
		parted = append(parted, part)
		idx = i + 1
	}

	// parse option part to Option
	quoted = false
	for index, part := range parted {
		part = strings.TrimSpace(part)
		if index == 0 {
			f.name = part
			if !stringsx.ValidFlagName(f.name) {
				panic(codex.Errorf(ECODE__INVALID_FLAG_NAME, "name: %q", f.name))
			}
			continue
		}
		if part == "" {
			continue
		}
		eq := false
		for i, c := range []rune(part) {
			switch c {
			case '\'':
				quoted = !quoted
			case '=':
				if !quoted {
					eq = true
					goto FinishOption
				}
			}
			if i == len(part)-1 {
				i++
				goto FinishOption
			}
			continue
		FinishOption:
			k, v := "", ""
			if eq {
				k = strings.TrimSpace(string([]rune(part)[:i]))
				v = strings.TrimSpace(string([]rune(part)[i+1:]))
			} else {
				k = part
			}
			k = unquote(k)
			if !stringsx.ValidFlagOptionKey(k) {
				panic(codex.Errorf(ECODE__INVALID_OPTION_KEY, "key: %q", k))
			}
			// if !strings.HasPrefix(v, "'") && !strings.HasSuffix(v, "'") {
			// 	// if option value contains control characters or spaces,
			// 	// it MUST be quoted with `'`.
			// 	if !stringsx.ValidUnquotedOptionValue(v) {
			// 		panic(codex.Errorf(ECODE__INVALID_OPTION_VALUE, "value: %q", v))
			// 	}
			// }
			if opt := NewOption(k, v, index); !opt.IsZero() {
				if _, exists := f.options[opt.key]; !exists {
					f.options[opt.key] = opt
				}
			}
			break
		}
	}
}

func (f *Flag) Key() string {
	return f.key
}

func (f *Flag) Name() string {
	return f.name
}

func (f *Flag) HasOption(key string) bool {
	_, ok := f.options[key]
	return ok
}

func (f *Flag) Option(key string) *Option {
	return f.options[key]
}

func (f *Flag) OptionLen() int {
	return len(f.options)
}

func (f *Flag) Options() iter.Seq[*Option] {
	return func(yield func(*Option) bool) {
		for _, opt := range f.options {
			yield(opt)
		}
	}
}

func (f *Flag) Value() string {
	if f.value != "" {
		return f.value
	}

	options := slices.SortedFunc(
		maps.Values(f.options),
		func(a *Option, b *Option) int {
			return a.index - b.index
		},
	)
	parts := []string{f.name}
	for _, opt := range options {
		parts = append(parts, opt.String())
	}

	f.value = strings.Join(parts, ",")
	f.value = strconv.Quote(f.value)
	return f.value
}

func (f *Flag) String() string {
	if f.prettied == "" {
		f.prettied = f.key + ":" + f.Value()
	}
	return f.prettied
}

func NewOption(key string, val string, offset int) *Option {
	return &Option{
		key:      key,
		val:      val,
		quoted:   quote(val),
		unquoted: unquote(val),
		index:    offset,
	}
}

type Option struct {
	key      string
	val      string
	quoted   string
	unquoted string
	index    int
}

func (o *Option) IsZero() bool { return o.key == "" }

func (o *Option) String() string {
	if o.IsZero() {
		return ""
	}
	if len(o.val) > 0 {
		return o.key + "=" + o.val
	}
	return o.key
}

func (o *Option) Key() string {
	return o.key
}

// Value returns literal text of option.
// eg: `db:"f_id,auto_inc,default='100'` Value() returns '100', Unquoted returns 100
func (o *Option) Value() string {
	if !o.IsZero() {
		return o.val
	}
	return ""
}

func (o *Option) Quoted() string {
	if !o.IsZero() {
		return o.quoted
	}
	return ""
}

func (o *Option) Unquoted() string {
	if !o.IsZero() {
		return o.unquoted
	}
	return ""
}

func unquote(s string) string {
	if strings.HasPrefix(s, "'") {
		s = strings.TrimPrefix(s, "'")
	}
	if strings.HasSuffix(s, "'") {
		s = strings.TrimSuffix(s, "'")
	}
	return s
}

func quote(s string) string {
	if !strings.HasPrefix(s, "'") {
		s = "'" + s
	}
	if !strings.HasSuffix(s, "'") {
		s = s + "'"
	}
	return s
}
