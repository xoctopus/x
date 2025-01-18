package internal

import (
	"go/ast"
	"go/parser"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

const (
	pkgPrefix = "xoctopus__"
	pkgSuffix = "__xoctopus"
)

func WrapPkgPath(path string) string {
	must.BeTrueWrap(
		!strings.Contains(path, "_"),
		"unsupported pkg path which contains `_`",
	)
	if !strings.Contains(path, "/") && !strings.Contains(path, ".") {
		return path
	}
	path = strings.ReplaceAll(strings.ReplaceAll(path, "/", "_"), ".", "__")
	return pkgPrefix + path + pkgSuffix
}

func UnwrapPkgPath(path string) string {
	if strings.HasPrefix(path, pkgPrefix) && strings.HasSuffix(path, pkgSuffix) {
		path = strings.TrimPrefix(path, pkgPrefix)
		path = strings.TrimSuffix(path, pkgSuffix)
		path = strings.ReplaceAll(strings.ReplaceAll(path, "__", "."), "_", "/")
	}
	return path
}

func WrapNamedTypeID(id string) string {
	idx := strings.Index(id, "[")
	if idx == 0 {
		return id
	}
	raw, params := id, ""
	if idx > 0 {
		raw, params = id[0:idx], id[idx+1:len(id)-1]
	}
	path, typename := "", raw
	if dot := strings.LastIndex(raw, "."); dot > 0 {
		path, typename = raw[:dot], raw[dot+1:]
	}
	b := strings.Builder{}
	if path == "" {
		b.WriteString(typename)
	} else {
		b.WriteString(WrapPkgPath(path) + "." + typename)
	}
	if len(params) > 0 {
		param, embeds, seq := make([]rune, 0), 0, 0
		b.WriteString("[")
		for s, i := id[idx+1:len(id)-1], 0; i < len(s); i++ {
			switch s[i] {
			case '[':
				embeds++
			case ']':
				embeds--
				must.BeTrue(embeds >= 0)
			case ',':
				if embeds == 0 {
					goto FinishParam
				}
			}
			param = append(param, rune(s[i]))
			if i == len(s)-1 {
				goto FinishParam
			}
			continue
		FinishParam:
			if seq > 0 {
				b.WriteString(",")
			}
			b.WriteString(WrapNamedTypeID(string(param)))
			param = param[0:0]
			seq++
		}
		b.WriteString("]")
	}
	return b.String()
}

func NewWrappedType(wrapped string) *WrappedType {
	if kind, ok := gBasicKindNames[wrapped]; ok {
		return &WrappedType{kind: kind, typename: wrapped}
	}
	exp, err := parser.ParseExpr(wrapped)
	must.NoError(err)

	ident := func(src string, x ast.Node) string {
		return src[x.Pos()-1 : x.End()-1]
	}

	switch e := exp.(type) {
	case *ast.MapType:
		return &WrappedType{
			kind: reflect.Map,
			key:  NewWrappedType(ident(wrapped, e.Key)),
			elem: NewWrappedType(ident(wrapped, e.Value)),
		}
	case *ast.ArrayType:
		t := &WrappedType{
			kind: reflect.Slice,
			elem: NewWrappedType(ident(wrapped, e.Elt)),
		}
		if e.Len != nil {
			t.kind = reflect.Array
			t.len, _ = strconv.Atoi(ident(wrapped, e.Len))
		}
		return t
	case *ast.ChanType:
		return &WrappedType{
			kind: reflect.Chan,
			dir:  gChanDirAst2Types[e.Dir],
			elem: NewWrappedType(ident(wrapped, e.Value)),
		}
	case *ast.FuncType:
		t := &WrappedType{kind: reflect.Func}
		if e.Params != nil {
			t.params = make([]*WrappedType, len(e.Params.List))
			for i, p := range e.Params.List {
				param := ident(wrapped, p.Type)
				if strings.HasPrefix(param, "...") {
					t.variadic = true
					t.params[i] = NewWrappedType("[]" + param[3:])
					break
				}
				t.params[i] = NewWrappedType(param)
			}
		}
		if e.Results != nil {
			t.results = make([]*WrappedType, len(e.Results.List))
			for i, r := range e.Results.List {
				t.results[i] = NewWrappedType(ident(wrapped, r.Type))
			}
		}
		return t
	case *ast.StarExpr:
		return &WrappedType{
			kind: reflect.Pointer,
			elem: NewWrappedType(ident(wrapped, e.X)),
		}
	case *ast.StructType:
		t := &WrappedType{kind: reflect.Struct}
		if e.Fields != nil {
			t.fields = make([]*WrappedType, len(e.Fields.List))
			for i, f := range e.Fields.List {
				t.fields[i] = NewWrappedType(ident(wrapped, f.Type))
				t.fields[i].anonymous = len(f.Names) == 0
				if t.fields[i].anonymous {
					t.fields[i].name = t.fields[i].typename
				} else {
					t.fields[i].name = f.Names[0].Name
				}
				if f.Tag != nil {
					t.fields[i].tag = f.Tag.Value
				}
			}
		}
		return t
	case *ast.InterfaceType:
		t := &WrappedType{
			kind:    reflect.Interface,
			methods: make([]*WrappedType, len(e.Methods.List)),
		}
		for i, m := range e.Methods.List {
			t.methods[i] = NewWrappedType("func" + ident(wrapped, m.Type))
			t.methods[i].name = m.Names[0].Name
		}
		return t
	default:
		t := &WrappedType{}
		typename := func(expr ast.Expr) (string, string) {
			switch x := expr.(type) {
			case *ast.SelectorExpr:
				return ident(wrapped, x.X), ident(wrapped, x.Sel)
			case *ast.Ident:
				return "", x.Name
			default:
				panic(errors.Errorf("invalid typename expr: [%T] %s", expr, ident(wrapped, expr)))
			}
		}
		switch named := exp.(type) {
		case *ast.IndexExpr:
			t.pkg, t.typename = typename(named.X)
			t.args = []*WrappedType{NewWrappedType(ident(wrapped, named.Index))}
		case *ast.IndexListExpr:
			t.pkg, t.typename = typename(named.X)
			t.args = make([]*WrappedType, len(named.Indices))
			for i, index := range named.Indices {
				t.args[i] = NewWrappedType(ident(wrapped, index))
			}
		default: // *ast.SelectorExpr *ast.Ident
			t.pkg, t.typename = typename(named)
		}
		return t
	}
}

type WrappedType struct {
	kind      reflect.Kind
	dir       types.ChanDir
	key       *WrappedType
	elem      *WrappedType
	len       int
	pkg       string
	typename  string
	name      string
	args      []*WrappedType
	params    []*WrappedType
	results   []*WrappedType
	variadic  bool
	fields    []*WrappedType
	tag       string
	anonymous bool
	methods   []*WrappedType
}

func (t *WrappedType) String() string {
	switch t.kind {
	case
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.UnsafePointer:
		return t.typename
	case reflect.Array:
		return "[" + strconv.Itoa(t.len) + "]" + t.elem.String()
	case reflect.Chan:
		return ChanDirPrefix(t.dir) + t.elem.String()
	case reflect.Func:
		b := strings.Builder{}
		b.WriteString("func(")
		for i, p := range t.params {
			if i > 0 {
				b.WriteString(", ")
			}
			if i == len(t.params)-1 && t.variadic {
				b.WriteString("...")
				b.WriteString(p.elem.String())
			} else {
				b.WriteString(p.String())
			}
		}
		b.WriteString(")")
		if len(t.results) > 0 {
			b.WriteString(" ")
			if len(t.results) > 1 {
				b.WriteString("(")
			}
			for i, r := range t.results {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(r.String())
			}
			if len(t.results) > 1 {
				b.WriteString(")")
			}
		}
		return b.String()
	case reflect.Interface:
		if t.pkg == "" && t.typename == "error" {
			return t.typename
		}
		if len(t.methods) == 0 {
			return "interface {}"
		}
		b := strings.Builder{}
		b.WriteString("interface { ")
		for i, m := range t.methods {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(m.name)
			b.WriteString(strings.TrimPrefix(m.String(), "func"))
		}
		b.WriteString(" }")
		return b.String()
	case reflect.Map:
		return "map[" + t.key.String() + "]" + t.elem.String()
	case reflect.Pointer:
		return "*" + t.elem.String()
	case reflect.Slice:
		return "[]" + t.elem.String()
	case reflect.Struct:
		if len(t.fields) == 0 {
			return "struct {}"
		}
		b := strings.Builder{}
		b.WriteString("struct { ")
		for i, f := range t.fields {
			if i > 0 {
				b.WriteString("; ")
			}
			if !f.anonymous {
				b.WriteString(f.name)
				b.WriteString(" ")
			}
			b.WriteString(f.String())
			if f.tag != "" {
				b.WriteString(" ")
				b.WriteString(f.tag)
			}
		}
		b.WriteString(" }")
		return b.String()
	default:
		b := strings.Builder{}
		if t.pkg != "" {
			b.WriteString(t.pkg)
			b.WriteString(".")
		}
		b.WriteString(t.typename)
		if len(t.args) > 0 {
			b.WriteString("[")
			for i, p := range t.args {
				if i > 0 {
					b.WriteString(",")
				}
				b.WriteString(p.String())
			}
			b.WriteString("]")
		}
		return b.String()
	}
}

func (t *WrappedType) Params() *types.Tuple {
	tuple := make([]*types.Var, len(t.params))
	for i := range len(t.params) {
		v := t.params[i]
		tuple[i] = types.NewParam(0, NewPackage(v.pkg), "", NewTypesTypeByID(v.String()))
	}
	return types.NewTuple(tuple...)
}

func (t *WrappedType) Results() *types.Tuple {
	tuple := make([]*types.Var, len(t.results))
	for i := range len(t.results) {
		v := t.results[i]
		tuple[i] = types.NewParam(0, NewPackage(v.pkg), "", NewTypesTypeByID(v.String()))
	}
	return types.NewTuple(tuple...)
}

func (t *WrappedType) Fields() ([]*types.Var, []string) {
	fields, tags := make([]*types.Var, len(t.fields)), make([]string, len(t.fields))
	for i := range t.fields {
		f := t.fields[i]
		fields[i] = types.NewField(0, NewPackage(f.pkg), f.name, NewTypesTypeByID(f.String()), f.anonymous)
		tags[i] = f.tag
	}
	return fields, tags
}

func (t *WrappedType) Methods() []*types.Func {
	fns := make([]*types.Func, len(t.methods))
	for i := range len(t.methods) {
		m := t.methods[i]
		fns[i] = types.NewFunc(0, NewPackage(m.pkg), m.name, NewTypesTypeByID(m.String()).(*types.Signature))
	}
	return fns
}
