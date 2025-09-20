package reflectx

import (
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

// ParseTag parses a struct tag into a map of flag keys and values.
// It scans the tag string for key-value pairs in the format `key:"value"`
// and parse the value to flag name and options
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

		// meet flag name
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

		// meet flag value and unquote it
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			panic(ErrInvalidFlagRaw)
		}
		quoted := string(tag[:i+1])
		raw, _ := strconv.Unquote(quoted)
		if _, ok := flags[key]; !ok {
			flags[key] = &Flag{key: key, raw: raw}
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
	key     string
	name    string
	options map[string]*Option
	raw     string
	pretty  string
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
		case '\\':
			if quoted {
				i++
				continue
			}
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
		parted = append(parted, string([]rune(val)[idx:i]))
		idx = i + 1
	}

	if len(parted) == 0 {
		f.raw = ""
		f.name = ""
		f.options = make(map[string]*Option)
		return
	}

	f.name = strings.TrimSpace(parted[0])
	if !validate(f.name) {
		panic(ErrInvalidFlagName)
	}

	// parse option part to Option
	quoted = false
	for index, part := range parted[1:] {
		part = strings.TrimSpace(part)
		eq := false
		for i, c := range []rune(part) {
			switch c {
			case '\'':
				quoted = !quoted
			case '\\':
				if quoted {
					i++
					continue
				}
			case '=':
				if !quoted {
					eq = true
					goto FinishOption
				}
			}
			if i == len(part)-1 {
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
			if !validate(opt.key) {
				panic(ErrInvalidOptionKey)
			}

			opt.key = unquote(opt.key)
			opt.val = quote(opt.val)

			if !opt.IsZero() {
				if _, exists := f.options[opt.key]; !exists {
					f.options[opt.key] = opt
				}
			}
		}
	}
}

func (f *Flag) Key() string { return f.key }

func (f *Flag) Name() string { return f.name }

func (f *Flag) WithRaw(raw string) *Flag {
	f2 := &Flag{
		key:     f.key,
		raw:     raw,
		options: make(map[string]*Option),
	}
	f2.parse()
	return f2
}

func (f *Flag) Option(key string) *Option {
	return f.options[key]
}

func (f *Flag) OptionLen() int { return len(f.options) }

func (f *Flag) Raw() string { return f.raw }

func (f *Flag) Value() string {
	if f.pretty != "" {
		return f.pretty
	}

	options := maps.Values(f.options)
	sort.Slice(options, func(i, j int) bool {
		return options[i].index < options[j].index
	})

	parts := []string{f.name}
	for _, opt := range options {
		parts = append(parts, opt.String())
	}

	f.pretty = strings.Join(parts, ",")
	return f.pretty
}

func (f *Flag) String() string {
	return f.key + ":" + strconv.Quote(f.Value())
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
	return o.key + "='" + o.val + "'"
}

func (o *Option) Key() string { return o.key }

func (o *Option) Value() string { return o.val }

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
		if !(c >= 'a' && c < 'z' ||
			c > 'A' && c < 'Z' ||
			c > '0' && c < '9' ||
			c == '_') {
			return false
		}
	}
	return true
}

var (
	ErrInvalidFlagRaw        = errors.New("")
	ErrInvalidFlagKey        = errors.New("")
	ErrInvalidFlagName       = errors.New("")
	ErrInvalidOptionKey      = errors.New("")
	ErrInvalidOptionUnquoted = errors.New("")
)
