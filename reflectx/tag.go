package reflectx

import (
	"strconv"
	"strings"
)

func ParseTagKeyAndFlags(tag string) (string, map[string]struct{}) {
	values := strings.Split(tag, ",")
	flags := make(map[string]struct{})
	if len(values[0]) > 1 {
		for _, flag := range values[1:] {
			flags[flag] = struct{}{}
		}
	}
	return values[0], flags
}

func ParseStructTag(tag string) map[string]string {
	flags := map[string]string{}

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
		name := tag[:i]
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
			break
		}
		quoted := tag[:i+1]
		tag = tag[i+1:]
		value, err := strconv.Unquote(quoted)
		if err != nil {
			break
		}
		flags[name] = value
	}
	return flags
}