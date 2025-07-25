package reflectx

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

// ParseFlags parses struct tag annotations into a structured Flags map.
// It supports tag formats like: `key:"value,opt1,opt2=v2"`.
// Each key-value pair in the tag is processed into a Flag with Value and Options.
// Conflicting tags (i.e., duplicate keys) will be detected and removed.
// The function also handles quoted values and options with or without values.
func ParseFlags(tag reflect.StructTag) Flags {
	var (
		_tag   = tag
		_flags = map[string]string{}
	)

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
		name := string(tag[:i])
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
			break // not a quoted string
		}
		quoted := string(tag[:i+1])
		value, _ := strconv.Unquote(quoted)
		if _, ok := _flags[name]; ok {
			fmt.Printf("[WARN] tag `%s` conflict in [ %s ]\n", name, _tag)
			return nil
		}
		_flags[name] = value
		tag = tag[i+1:]
	}

	flags := Flags{}

	for name, value := range _flags {
		f := (*Flag)(nil)
		runes := []rune(value)
		stage := 0
		option := [2]string{}
		quoted := false
		prev := 0
		for curr := 0; curr < len(runes); curr++ {
			switch runes[curr] {
			case ',':
				if quoted {
					continue
				}
				goto FinishStage
			case '\\':
				curr++
				continue
			case '\'':
				quoted = !quoted
				if quoted {
					continue
				}
				curr++
				goto FinishStage
			case '=':
				if quoted {
					continue
				}
				goto FinishStage
			}

			if curr == len(runes)-1 {
				curr++
				goto FinishStage
			}
			continue

		FinishStage:
			{
				switch stage {
				case 0: // finish Flag.Value
					f = NewFlag(strings.TrimSpace(string(runes[prev:curr])))
					stage = 1
				case 1: // finish option[0](option name)
					option[0] = strings.TrimSpace(string(runes[prev:curr]))
					if len(option[0]) > 0 && option[0][0] == '\'' {
						if option[0][len(option[0])-1] != '\'' {
							fmt.Printf("[WARN] option key %s unquoted\n", option[0])
							return Flags{}
						}
						option[0] = option[0][1 : len(option[0])-1]
					}
					if curr < len(runes) {
						if runes[curr] == '=' {
							stage = 2
						} else {
							stage = 1
						}
					}
					if curr == len(runes)-1 || runes[curr] == ',' {
						f.AddOption(option[0], option[1])
						option = [2]string{}
					}
				default: // finish option[1] (option value)
					option[1] = strings.TrimSpace(string(runes[prev:curr]))
					if option[1][0] == '\'' {
						if option[1][len(option[1])-1] != '\'' {
							fmt.Printf("[WARN] option value %s unquoted\n", option[0])
							return Flags{}
						}
						option[1] = option[1][1 : len(option[1])-1]
					}
					f.AddOption(option[0], option[1])
					option = [2]string{}
					stage = 1
				}
				prev = curr + 1
			}
		}
		flags.Add(name, f)
	}

	return flags
}

func NewFlag(name string) *Flag {
	return &Flag{
		name:      name,
		options:   make(map[string]string),
		conflicts: make(map[string]struct{}),
	}
}

// Flag represents a single tag value with optional key-value options.
// For example, the tag `name:"field,opt1,opt2=v2"` will be parsed into:
// {name: "field", options: [][2]string{{"opt1": ""}, {"opt2": "v2"}}}
type Flag struct {
	name      string
	options   map[string]string
	conflicts map[string]struct{}
}

func (f *Flag) Name() string { return f.name }

func (f *Flag) Option(key string) string { return f.options[key] }

func (f *Flag) AddOption(k, v string) {
	if k == "" {
		return
	}
	if _, ok := f.conflicts[k]; ok {
		goto RemoveConflict
	}
	if _, ok := f.options[k]; ok {
		goto RemoveConflict
	}
	f.options[k] = v
	return
RemoveConflict:
	fmt.Printf("[WARN] option `%s` conflicts", k)
	f.conflicts[k] = struct{}{}
	delete(f.options, k)
	return
}

func (f *Flag) String() string {
	if len(f.options) == 0 {
		return f.name
	}

	keys := maps.Keys(f.options)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	options := make([]string, len(keys))
	for i, key := range keys {
		val := f.options[key]
		if val == "" {
			options[i] = key
		} else {
			options[i] = key + "='" + strings.Trim(strconv.Quote(val), "\"") + "'"
		}
	}
	return f.name + "," + strings.Join(options, ",")
}

type Flags map[string]*Flag

func (fs Flags) Get(key string) *Flag {
	if f, ok := fs[key]; ok {
		return f
	}
	return nil
}

func (fs Flags) Add(name string, f *Flag) (conflicted bool) {
	_, conflicted = fs[name]
	fs[name] = f
	return conflicted
}

func (fs Flags) Delete(name string) {
	delete(fs, name)
}

func (fs Flags) String() string {
	names := maps.Keys(fs)
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	s := ""
	for i, name := range names {
		s += name + `:"`
		s += fs[name].String()
		s += `"`
		if i != len(names)-1 {
			s += " "
		}
	}
	return s
}
