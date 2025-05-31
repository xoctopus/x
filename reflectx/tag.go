package reflectx

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ParseFlags parses struct tag annotations into a structured Flags map.
// It supports tag formats like: `key:"value,opt1,opt2=v2"`.
// Each key-value pair in the tag is processed into a Flag with Value and Options.
// Conflicting tags (i.e., duplicate keys) will be detected and removed.
// The function also handles quoted values and options with or without values.
func ParseFlags(tag reflect.StructTag) Flags {
	var (
		full      = tag
		flags     = map[string]string{}
		conflicts = map[string]struct{}{}
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
		if _, ok := flags[name]; ok {
			fmt.Printf("[WARN] tag `%s` conflict in [ %s ]\n", name, full)
			conflicts[name] = struct{}{}
			continue
		}
		flags[name] = value
		tag = tag[i+1:]
	}

	for name := range conflicts {
		delete(flags, name)
	}

	results := Flags{}

	for key, value := range flags {
		f := Flag{Tag: key}
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
					f.Value = strings.TrimSpace(string(runes[prev:curr]))
					stage = 1
				case 1: // finish option[0]
					option[0] = strings.TrimSpace(string(runes[prev:curr]))
					if runes[curr] == '=' && curr != len(runes)-1 {
						stage = 2
					} else {
						if option[0] != "" {
							option[0] = strings.TrimSpace(string(runes[prev:curr]))
							option[0] = strings.TrimPrefix(option[0], "'")
							option[0] = strings.TrimSuffix(option[0], "'")
							f.Options = append(f.Options, option)
							option = [2]string{}
							stage = 1
						}
					}
				default: // finish option[1]
					if option[0] != "" {
						option[1] = strings.TrimSpace(string(runes[prev:curr]))
						option[1] = strings.TrimPrefix(option[1], "'")
						option[1] = strings.TrimSuffix(option[1], "'")
						f.Options = append(f.Options, option)
						option = [2]string{}
						stage = 1
					}
				}
				prev = curr + 1
			}
		}
		results[f.Tag] = &f
	}

	for name := range results {
		sort.Slice(results[name].Options, func(i, j int) bool {
			return results[name].Options[i][0] < results[name].Options[j][0]
		})
	}

	return results
}

type Flag struct {
	Tag     string
	Value   string
	Options [][2]string
}

type Flags map[string]*Flag

func (fs Flags) Get(key string) *Flag {
	if f, ok := fs[key]; ok {
		return f
	}
	return nil
}
