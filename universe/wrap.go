package universe

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

var (
	brackets = map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}
	keywords = map[string]any{
		"[":           struct{}{},
		"map[":        struct{}{},
		"func(":       struct{}{},
		"interface {": struct{}{},
		"struct {":    struct{}{},
		"chan ":       struct{}{},
		"chan<- ":     struct{}{},
		"<-chan ":     struct{}{},
		"*":           struct{}{},
	}
)

// bracketed returns the sub string of `id` bracketed by `identifier` and the indexes
// of brackets
func bracketed(id string, identifier0 rune) (string, int, int) {
	identifier1, ok := brackets[identifier0]
	must.BeTrueWrap(ok, "invalid bracket identifier: %v", identifier0)

	l, r, embeds, quoted := -1, -1, 0, false
End:
	for i := 0; i < len(id); i++ {
		c := rune(id[i])
		switch c {
		case identifier0:
			if quoted {
				continue
			}
			if embeds == 0 {
				l = i
			}
			embeds++
		case identifier1:
			if quoted {
				continue
			}
			embeds--
			if embeds == 0 {
				must.BeTrue(l >= 0)
				r = i
				break End
			}
		case '"':
			quoted = !quoted
		case '\\':
			i++
		}
	}
	if l < 0 && r < 0 {
		return "", l, r
	}
	return id[l+1 : r], l, r
}

// quoted returns the quoted sub string and quoter indexes.
func quoted(id string) (string, int, int) {
	quoted, l, r := false, -1, -1
	for i, c := range id {
		if i < len(id)-1 && id[i+1] == '\\' {
			i += 2
			continue
		}
		if c == '"' {
			if !quoted && l < 0 {
				quoted = true
				l = i
				continue
			}
			if quoted && l >= 0 {
				r = i
				break
			}
		}
	}
	return id[l+1 : r], l, r
}

// separate separates id by `sep`, unlike strings.Split, it will ignore `sep` within
// code blocks. eg:
// id = `struct { ... }`; sep = ' ', returns `struct`, `{ ... }`
// id = path/to/pkg.Typename[T1,T2]; sep = ',', returns `path/to/pkg.Typename[T1,T2]`
func separate(id string, sep rune) []string {
	must.BeTrue(sep == ',' || sep == ';' || sep == ' ')

	if len(id) == 0 {
		return nil
	}

	var (
		parts  = make([]string, 0)
		part   []rune
		embeds [3]int
		quoted bool
	)

	for i := 0; i < len(id); i++ {
		c := rune(id[i])
		switch c {
		case '(':
			embeds[0]++
		case ')':
			embeds[0]--
			must.BeTrue(embeds[0] >= 0)
		case '[':
			embeds[1]++
		case ']':
			embeds[1]--
			must.BeTrue(embeds[1] >= 0)
		case '{':
			embeds[2]++
		case '}':
			embeds[2]--
			must.BeTrue(embeds[2] >= 0)
		case '"':
			quoted = !quoted
		case sep:
			if embeds[0] == 0 && embeds[1] == 0 && embeds[2] == 0 && !quoted {
				goto FinishPart
			}
		}
		part = append(part, c)
		if c == '\\' {
			part = append(part, rune(id[i+1]))
			i++
		}
		if i == len(id)-1 {
			goto FinishPart
		}
		continue
	FinishPart:
		p := strings.TrimSpace(string(part))
		must.BeTrue(len(p) > 0)
		parts = append(parts, p)
		part = part[0:0]
	}

	return parts
}

// reverse return reversed string
func reverse(id string) string {
	_id := []rune(id)
	slices.Reverse(_id)
	return string(_id)
}

// field parses the struct field info, returns name, type and tag of field
func field(id string) (name string, typ string, tag string) {
	if id[len(id)-1] == '"' {
		id = reverse(id)

		_tag, ql, qr := quoted(id)
		must.BeTrue(ql >= 0 && qr > 0)
		_tag = "\"" + reverse(_tag) + "\""
		_tag, err := strconv.Unquote(_tag)
		must.NoError(err)
		tag = _tag

		id = strings.TrimSpace(reverse(id[qr+1:]))
	}

	switch parts := separate(id, ' '); len(parts) {
	case 1:
		typ = parts[0]
	case 2:
		name, typ = parts[0], parts[1]
	default:
		if len(parts) == 3 && parts[1] == "=" {
			typ = parts[2]
		} else {
			name, typ = parts[0], strings.Join(parts[1:], " ")
		}
	}
	return
}

// WrapID wraps the type id. if the id is a named type wrap the package path to
// help to parse to an ast.Expr. eg:
// github.com/path/to/pkg.TypeName => xoctopus__github__com_path_to_pkg__xoctopus.TypeName
func WrapID(id string) (wrapped string) {
	if _id, ok := gWrapID.Load(id); ok {
		return _id.(string)
	}

	defer func(id string) {
		gWrapID.Store(wrapped, wrapped)
		gWrapID.Store(id, wrapped)
	}(id)

	if id == "error" {
		return id
	}

	// basics
	if kind, ok := gBasicName2Kind[id]; ok {
		return kind.String()
	}

	// slice: []elem / array: [len]elem
	if strings.HasPrefix(id, "[") {
		sub, _, r := bracketed(id, '[')
		return "[" + sub + "]" + WrapID(id[r+1:])
	}

	// map: map[key]elem
	if strings.HasPrefix(id, "map[") {
		sub, _, r := bracketed(id, '[')
		return "map[" + WrapID(sub) + "]" + WrapID(id[r+1:])
	}

	if strings.HasPrefix(id, "chan ") {
		return "chan " + WrapID(id[5:])
	}

	if strings.HasPrefix(id, "chan<- ") {
		return "chan<- " + WrapID(id[7:])
	}

	if strings.HasPrefix(id, "<-chan ") {
		return "<-chan " + WrapID(id[7:])
	}

	// struct: struct { fields... }
	if strings.HasPrefix(id, "struct {") {
		sub, _, _ := bracketed(id, '{')
		if len(sub) == 0 {
			return "struct {}"
		}

		b := strings.Builder{}
		for i, p := range separate(sub, ';') {
			if i > 0 {
				b.WriteString("; ")
			}
			name, typ, tag := field(p)
			if len(name) > 0 {
				b.WriteString(name)
				b.WriteString(" ")
			}
			b.WriteString(WrapID(typ))
			if len(tag) > 0 {
				b.WriteString(" ")
				b.WriteString(strconv.Quote(tag))
			}
		}
		return "struct { " + b.String() + " }"
	}

	// interface: interface { methods... }
	if strings.HasPrefix(id, "interface {") {
		sub, _, _ := bracketed(id, '{')
		if len(sub) == 0 {
			return "interface {}"
		}

		b := strings.Builder{}
		for i, p := range separate(sub, ';') {
			if i > 0 {
				b.WriteString("; ")
			}
			idx := strings.Index(p, "(")
			typ := WrapID("func" + p[idx:])
			b.WriteString(p[0:idx] + typ[4:])
		}
		return "interface { " + b.String() + " }"
	}

	// func: func(params...) results
	if strings.HasPrefix(id, "func(") {
		params, _, pr := bracketed(id, '(')

		b := strings.Builder{}
		b.WriteString("func(")
		for i, p := range separate(params, ',') {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(WrapID(p))
		}
		b.WriteString(")")

		if pr == len(id)-1 {
			return b.String()
		}

		b.WriteString(" ")
		id = strings.TrimSpace(id[pr+1:])
		if id[0] == '(' {
			b.WriteString("(")
			results, _, _ := bracketed(id, '(')
			for i, r := range separate(results, ',') {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(WrapID(r))
			}
			b.WriteString(")")
		} else {
			b.WriteString(id)
		}
		return b.String()
	}

	// pointer: *elem
	if strings.HasPrefix(id, "*") {
		return "*" + WrapID(id[1:])
	}

	// variadic: ...elem
	if strings.HasPrefix(id, "...") {
		return "..." + WrapID(id[3:])
	}

	// named: package_path.typename[type arguments...]
	path, name, args := "", "", ""
	if _args, l, r := bracketed(id, '['); l > 0 && r > 0 {
		args = _args
		id = id[0:l]
	}
	dot := strings.LastIndex(id, ".")
	must.BeTrue(dot != -1)
	path, name = id[0:dot], id[dot+1:]
	b := strings.Builder{}
	b.WriteString(wrap(path))
	b.WriteString(".")
	b.WriteString(name)
	if len(args) > 0 {
		b.WriteString("[")
		for i, a := range separate(args, ',') {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(WrapID(a))
		}
		b.WriteString("]")
	}
	return b.String()
}

func WrapIDByT(t reflect.Type) (wrapped string) {
	if t == nil {
		return ""
	}

	if _id, ok := gWrapID.Load(t); ok {
		return _id.(string)
	}

	defer func() {
		gWrapID.Store(t, wrapped)
	}()

	if id := t.Name(); id != "" {
		if t.PkgPath() != "" {
			id = t.PkgPath() + "." + id
		}
		return WrapID(id)
	}

	switch t.Kind() {
	case reflect.Array:
		return "[" + strconv.Itoa(t.Len()) + "]" + WrapIDByT(t.Elem())
	case reflect.Chan:
		return gChanDirPrefix[t.ChanDir()] + WrapIDByT(t.Elem())
	case reflect.Func:
		b := strings.Builder{}
		b.WriteString("func(")
		for i := range t.NumIn() {
			if i > 0 {
				b.WriteString(", ")
			}
			if i == t.NumIn()-1 && t.IsVariadic() {
				b.WriteString("...")
				b.WriteString(WrapIDByT(t.In(i).Elem()))
				break
			}
			b.WriteString(WrapIDByT(t.In(i)))
		}
		b.WriteString(")")
		if t.NumOut() > 0 {
			b.WriteString(" ")
			if t.NumOut() > 1 {
				b.WriteString("(")
			}
			for i := range t.NumOut() {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(WrapIDByT(t.Out(i)))
			}
			if t.NumOut() > 1 {
				b.WriteString(")")
			}
		}
		return b.String()
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return "interface {}"
		}
		b := strings.Builder{}
		b.WriteString("interface { ")
		for i := range t.NumMethod() {
			if i > 0 {
				b.WriteString("; ")
			}
			m := t.Method(i)
			f := WrapIDByT(m.Type)
			b.WriteString(m.Name)
			b.WriteString(f[4:])
		}
		b.WriteString(" }")
		return b.String()
	case reflect.Map:
		return "map[" + WrapIDByT(t.Key()) + "]" + WrapIDByT(t.Elem())
	case reflect.Pointer:
		return "*" + WrapIDByT(t.Elem())
	case reflect.Slice:
		return "[]" + WrapIDByT(t.Elem())
	case reflect.Struct:
		if t.NumField() == 0 {
			return "struct {}"
		}
		b := strings.Builder{}
		b.WriteString("struct { ")
		for i := range t.NumField() {
			if i > 0 {
				b.WriteString("; ")
			}
			f := t.Field(i)
			if !f.Anonymous {
				b.WriteString(f.Name)
				b.WriteString(" ")
			}
			b.WriteString(WrapIDByT(f.Type))
			if len(f.Tag) > 0 {
				b.WriteString(" ")
				b.WriteString(strconv.Quote(string(f.Tag)))
			}
		}
		b.WriteString(" }")
		return b.String()
	default:
		return ""
	}
}
