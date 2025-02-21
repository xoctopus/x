package reflectx

import (
	"reflect"
	"strconv"
	"strings"
)

// ParseStructTag parse struct tag to tag key/value map
// eg: `tagKey:"tagName,tagFlag1,tagFlag2=v"` will parsed to
// map[string]string{"cmd": "tagName,tagFlag1,tagFlag2=v"}
func ParseStructTag(tag reflect.StructTag) map[string]string {
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
		value, err := strconv.Unquote(quoted)
		mustBeTrue(err == nil, "invalid quoted value")
		flags[name] = value
		tag = tag[i+1:]
	}
	return flags
}

// ParseTagValue parse tag value to name and flags
// eg: "tagName,tagFlag1,tagFlag2=v" will be parsed to
// name: tagName, flags: map[string]{"tagFlag1", "tagFlag2=v"}
func ParseTagValue(tagValue string) (string, map[string]struct{}) {
	values := strings.Split(tagValue, ",")
	flags := make(map[string]struct{})
	for _, flag := range values[1:] {
		flags[flag] = struct{}{}
	}
	return values[0], flags
}
