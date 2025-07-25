package reflectx

import (
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/xoctopus/x/misc/must"
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
		if !validate(key) {
			panic(ErrInvalidFlagKey)
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
			panic(ErrInvalidFlagValue)
		}
		quoted := string(tag[:i+1])
		if _, ok := flags[key]; !ok {
			flags[key] = &Flag{
				key:      key,
				quoted:   quoted,
				unquoted: must.NoErrorV(strconv.Unquote(quoted)),
				options:  make(map[string]*Option),
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
	quoted   string
	unquoted string
	value    string
	prettied string
}

func (f *Flag) parse() {
	val := strings.TrimSpace(f.unquoted)

	// scan value to parted by ','
	quoted := false
	parted := make([]string, 0)
	idx := 0
	for i, c := range []rune(val) {
		_ = val[i : i+1]
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
				panic(ErrInvalidOptionUnquoted)
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
			if !validate(f.name) {
				panic(ErrInvalidFlagName)
			}
			continue
		}
		if part == "" {
			continue
		}
		eq := false
		for i, c := range []rune(part) {
			_ = part[i : i+1]
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
			opt := &Option{index: index}
			if eq {
				opt.key = strings.TrimSpace(string([]rune(part)[:i]))
				opt.val = strings.TrimSpace(string([]rune(part)[i+1:]))
			} else {
				opt.key = part
			}
			opt.key = unquote(opt.key)
			if !validate(opt.key) {
				panic(ErrInvalidOptionKey)
			}
			if len(opt.val) > 2 && opt.val[0] != '\'' && opt.val[len(opt.val)-1] != '\'' {
				if !validate(opt.val) {
					panic(ErrInvalidOptionValue)
				}
			}
			opt.val = quote(opt.val)
			if !opt.IsZero() {
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

func (f *Flag) Option(key string) *Option {
	return f.options[key]
}

func (f *Flag) OptionLen() int {
	return len(f.options)
}

func (f *Flag) QuotedValue() string {
	return f.quoted
}

func (f *Flag) UnquotedValue() string {
	return f.unquoted
}

func (f *Flag) Value() string {
	if f.value != "" {
		return f.value
	}

	options := maps.Values(f.options)
	sort.Slice(options, func(i, j int) bool {
		return options[i].index < options[j].index
	})

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
	return &Option{key: key, val: val, index: offset}
}

type Option struct {
	key   string
	val   string
	index int
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

func (o *Option) Key() string { return o.key }

func (o *Option) Value() string {
	if o.IsZero() {
		return ""
	}
	return o.val
}

func (o *Option) RawValue() []byte {
	v := o.Value()
	if v != "" {
		v = unquote(v)
	}
	return []byte(v)
}

func unquote(s string) string {
	if len(s) > 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}

func quote(s string) string {
	if len(s) > 2 && s[0] != '\'' && s[len(s)-1] != '\'' {
		return `'` + s + `'`
	}
	return s
}

func validate(key string) bool {
	for _, c := range key {
		if !(c >= 'a' && c <= 'z' ||
			c >= 'A' && c <= 'Z' ||
			c >= '0' && c <= '9' ||
			c == '_' || c == '-') {
			return false
		}
	}
	return true
}

var (
	ErrInvalidFlagKey        = errors.New("invalid flag key")
	ErrInvalidFlagValue      = errors.New("invalid flag value")
	ErrInvalidFlagName       = errors.New("invalid flag name")
	ErrInvalidOptionKey      = errors.New("invalid option key")
	ErrInvalidOptionValue    = errors.New("invalid option value")
	ErrInvalidOptionUnquoted = errors.New("invalid option unquoted")
)
