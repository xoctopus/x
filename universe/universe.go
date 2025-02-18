package universe

import (
	"go/ast"
	"go/parser"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/resultx"
)

type Universe interface {
	UniverseKind() Kind
	String() string
}

func NewUniverse(t reflect.Type) (result resultx.Result[Universe]) {
	if t == nil {
		return resultx.Succeed[Universe](nil)
	}

	if u, ok := gUniverse.Load(t); ok {
		return resultx.Succeed(u.(Universe))
	}

	defer func() {
		if result.Succeed() {
			gUniverse.Store(t, result.Unwrap())
		}
	}()

	if id := t.Name(); id != "" {
		if path := t.PkgPath(); path != "" {
			id = path + "." + id
		}
		return NewUniverseByID(id)
	}

	u := &Unnamed{kind: t.Kind()}
	switch t.Kind() {
	case reflect.Array:
		u.len, u.elem = t.Len(), NewUniverse(t.Elem()).Unwrap()
	case reflect.Chan:
		u.dir, u.elem = ChanDir(t.ChanDir()).Unwrap(), NewUniverse(t.Elem()).Unwrap()
	case reflect.Func:
		u.variadic = t.IsVariadic()
		u.ins = make([]Universe, t.NumIn())
		for i := range u.ins {
			u.ins[i] = NewUniverse(t.In(i)).Unwrap()
		}
		u.outs = make([]Universe, t.NumOut())
		for i := range u.outs {
			u.outs[i] = NewUniverse(t.Out(i)).Unwrap()
		}
	case reflect.Interface:
		u.methods = make([]*Unnamed, t.NumMethod())
		for i := range u.methods {
			m := t.Method(i)
			u.methods[i] = &Unnamed{
				kind:     reflect.Func,
				name:     m.Name,
				ins:      make([]Universe, m.Type.NumIn()),
				outs:     make([]Universe, m.Type.NumOut()),
				variadic: m.Type.IsVariadic(),
			}
			for j := range u.methods[i].ins {
				u.methods[i].ins[j] = NewUniverse(m.Type.In(j)).Unwrap()
			}
			for j := range u.methods[i].outs {
				u.methods[i].outs[j] = NewUniverse(m.Type.Out(j)).Unwrap()
			}
		}
		sort.Slice(u.methods, func(i, j int) bool {
			return u.methods[i].name < u.methods[j].name
		})
	case reflect.Map:
		u.key, u.elem = NewUniverse(t.Key()).Unwrap(), NewUniverse(t.Elem()).Unwrap()
	case reflect.Pointer:
		u.elem = NewUniverse(t.Elem()).Unwrap()
	case reflect.Slice:
		u.elem = NewUniverse(t.Elem()).Unwrap()
	case reflect.Struct:
		u.fields = make([]*Field, t.NumField())
		for i := range t.NumField() {
			f := t.Field(i)
			u.fields[i] = &Field{
				name:     f.Name,
				typ:      NewUniverse(f.Type).Unwrap(),
				tag:      string(f.Tag),
				embedded: f.Anonymous,
			}
		}
	default:
		return resultx.Succeed[Universe](nil)
	}
	return resultx.Succeed[Universe](u)
}

func NewUniverseByID(id string) (result resultx.Result[Universe]) {
	id = WrapID(id)

	if kind, ok := gBasicName2Kind[id]; ok {
		return resultx.Succeed[Universe](&Named{
			name: kind.String(),
		})
	}

	ident := func(code string, x ast.Node) string {
		return code[x.Pos()-1 : x.End()-1]
	}

	expr := resultx.Unwrap(parser.ParseExpr(id))
	switch e := expr.(type) {
	case *ast.ArrayType:
		if e.Len != nil {
			return resultx.Succeed[Universe](&Unnamed{
				kind: reflect.Array,
				elem: NewUniverseByID(ident(id, e.Elt)).Unwrap(),
				len:  resultx.Unwrap(strconv.Atoi(ident(id, e.Len))),
			})
		}
		return resultx.Succeed[Universe](&Unnamed{
			kind: reflect.Slice,
			elem: NewUniverseByID(ident(id, e.Elt)).Unwrap(),
		})
	case *ast.ChanType:
		return resultx.Succeed[Universe](&Unnamed{
			kind: reflect.Chan,
			dir:  ChanDir(e.Dir).Unwrap(),
			elem: NewUniverseByID(ident(id, e.Value)).Unwrap(),
		})
	case *ast.FuncType:
		u := &Unnamed{kind: reflect.Func}
		if e.Params != nil {
			u.ins = make([]Universe, len(e.Params.List))
			for i, p := range e.Params.List {
				param := ident(id, p.Type)
				if i == len(e.Params.List)-1 && strings.HasPrefix(param, "...") {
					u.variadic = true
					u.ins[i] = NewUniverseByID("[]" + param[3:]).Unwrap()
					break
				}
				u.ins[i] = NewUniverseByID(param).Unwrap()
			}
		}
		if e.Results != nil {
			u.outs = make([]Universe, len(e.Results.List))
			for i, r := range e.Results.List {
				u.outs[i] = NewUniverseByID(ident(id, r.Type)).Unwrap()
			}
		}
		return resultx.Succeed[Universe](u)
	case *ast.InterfaceType:
		u := &Unnamed{
			kind:    reflect.Interface,
			methods: make([]*Unnamed, len(e.Methods.List)),
		}
		for i, m := range e.Methods.List {
			u.methods[i] = NewUniverseByID("func" + ident(id, m.Type)).Unwrap().(*Unnamed)
			u.methods[i].name = m.Names[0].Name
		}
		return resultx.Succeed[Universe](u)
	case *ast.MapType:
		return resultx.Succeed[Universe](&Unnamed{
			kind: reflect.Map,
			key:  NewUniverseByID(ident(id, e.Key)).Unwrap(),
			elem: NewUniverseByID(ident(id, e.Value)).Unwrap(),
		})
	case *ast.StarExpr:
		return resultx.Succeed[Universe](&Unnamed{
			kind: reflect.Pointer,
			elem: NewUniverseByID(ident(id, e.X)).Unwrap(),
		})
	case *ast.StructType:
		u := &Unnamed{kind: reflect.Struct}
		if e.Fields != nil {
			u.fields = make([]*Field, len(e.Fields.List))
			for i, f := range e.Fields.List {
				u.fields[i] = &Field{typ: NewUniverseByID(ident(id, f.Type)).Unwrap()}
				if u.fields[i].embedded = len(f.Names) == 0; u.fields[i].embedded {
					u.fields[i].name = u.fields[i].typ.(*Named).name
				} else {
					u.fields[i].name = f.Names[0].Name
				}
				if f.Tag != nil {
					u.fields[i].tag = resultx.Unwrap(strconv.Unquote(f.Tag.Value))
				}
			}
		}
		return resultx.Succeed[Universe](u)
	case *ast.Ident:
		return resultx.Succeed[Universe](&Named{name: e.Name})
	case *ast.SelectorExpr:
		return resultx.Succeed[Universe](&Named{
			pkg:  NewPackage(ident(id, e.X)),
			name: ident(id, e.Sel),
		})
	case *ast.IndexExpr:
		n := NewUniverseByID(ident(id, e.X)).Unwrap().(*Named)
		n.args = []Universe{NewUniverseByID(ident(id, e.Index)).Unwrap()}
		return resultx.Succeed[Universe](n)
	case *ast.IndexListExpr:
		n := NewUniverseByID(ident(id, e.X)).Unwrap().(*Named)
		n.args = make([]Universe, len(e.Indices))
		for i, index := range e.Indices {
			n.args[i] = NewUniverseByID(ident(id, index)).Unwrap()
		}
		return resultx.Succeed[Universe](n)
	default:
		panic(errors.Errorf("invalid id: [%T] %s", e, ident(id, e)))
	}
}

type Named struct {
	pkg  Package
	name string
	args []Universe
}

func (t *Named) UniverseKind() Kind {
	return TypeName
}

func (t *Named) String() string {
	b := strings.Builder{}
	if t.pkg != nil {
		b.WriteString(t.pkg.Path())
		b.WriteString(".")
	}
	b.WriteString(t.name)
	if len(t.args) > 0 {
		b.WriteString("[")
		for i, arg := range t.args {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(arg.String())
		}
		b.WriteString("]")
	}
	return b.String()
}

type Field struct {
	name     string
	typ      Universe
	tag      string
	embedded bool
}

func (v *Field) String() string {
	b := strings.Builder{}
	if !v.embedded {
		b.WriteString(v.name)
		b.WriteString(" ")
	}
	b.WriteString(v.typ.String())
	if len(v.tag) > 0 {
		b.WriteString(" ")
		b.WriteString(strconv.Quote(v.tag))
	}
	return b.String()
}

type Unnamed struct {
	kind     reflect.Kind
	name     string
	key      Universe
	elem     Universe
	len      int
	dir      types.ChanDir
	ins      []Universe
	outs     []Universe
	variadic bool
	fields   []*Field
	methods  []*Unnamed
}

func (t *Unnamed) UniverseKind() Kind {
	return TypeLit
}

func (t *Unnamed) String() string {
	switch t.kind {
	case reflect.Array:
		return "[" + strconv.Itoa(t.len) + "]" + t.elem.String()
	case reflect.Chan:
		return gChanDirPrefix[t.dir] + t.elem.String()
	case reflect.Func:
		b := strings.Builder{}

		if t.name == "" {
			b.WriteString("func(")
		} else {
			b.WriteString(t.name + "(")
		}
		for i := range t.ins {
			if i > 0 {
				b.WriteString(", ")
			}
			if i == len(t.ins)-1 && t.variadic {
				b.WriteString("..." + t.ins[i].String()[2:])
				break
			}
			b.WriteString(t.ins[i].String())
		}
		b.WriteString(")")

		if len(t.outs) == 0 {
			return b.String()
		}
		b.WriteString(" ")
		if len(t.outs) > 1 {
			b.WriteString("(")
		}
		for i, v := range t.outs {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(v.String())
		}
		if len(t.outs) > 1 {
			b.WriteString(")")
		}
		return b.String()
	case reflect.Interface:
		if len(t.methods) == 0 {
			return "interface {}"
		}
		b := strings.Builder{}
		b.WriteString("interface { ")
		for i, m := range t.methods {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(m.String())
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
			b.WriteString(f.String())
		}
		b.WriteString(" }")
		return b.String()
	default:
		panic(errors.Errorf("unexpected Unnamed kind: %s", t.kind))
	}
}
