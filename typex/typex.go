package typex

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"

	"github.com/xoctopus/x/typex/internal"
)

type Type interface {
	// Unwrap to types.Type or reflect.Type
	Unwrap() any

	PkgPath() string
	Name() string
	String() string
	Kind() reflect.Kind

	Implements(Type) bool
	AssignableTo(Type) bool
	ConvertibleTo(Type) bool
	Comparable() bool

	Key() Type
	Elem() Type
	Len() int

	NumField() int
	Field(int) StructField
	FieldByName(string) (StructField, bool)
	FieldByNameFunc(func(string) bool) (StructField, bool)

	NumMethod() int
	Method(int) Method
	MethodByName(string) (Method, bool)

	IsVariadic() bool
	NumIn() int
	In(int) Type
	NumOut() int
	Out(int) Type
}

type Method interface {
	PkgPath() string
	Name() string
	Type() Type
}

type StructField interface {
	PkgPath() string
	Name() string
	Type() Type
	Tag() reflect.StructTag
	Anonymous() bool
}

func NewT(t Type) *T {
	return &T{t}
}

type T struct {
	t Type
}

func (v *T) String() string {
	t := v.t
	if pkg := t.PkgPath(); pkg != "" {
		return pkg + "." + t.Name()
	}

	k := t.Kind()
	if _, ok := internal.KindsR2G[k]; ok {
		return k.String()
	}

	switch k {
	case reflect.Slice:
		return "[]" + NewT(t.Elem()).String()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), NewT(t.Elem()).String())
	case reflect.Chan:
		prefix := ""
		switch x := t.Unwrap().(type) {
		case reflect.Type:
			switch x.ChanDir() {
			case reflect.SendDir:
				prefix = "chan<- "
			case reflect.RecvDir:
				prefix = "<-chan "
			default:
				prefix = "chan "
			}
		case types.Type:
			switch x.(*types.Chan).Dir() {
			case types.SendOnly:
				prefix = "chan<- "
			case types.RecvOnly:
				prefix = "<-chan "
			default:
				prefix = "chan "
			}
		}
		return prefix + NewT(t.Elem()).String()
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", NewT(t.Key()).String(), NewT(t.Elem()).String())
	case reflect.Pointer:
		return "*" + NewT(t.Elem()).String()
	case reflect.Struct:
		buf := bytes.NewBuffer(nil)
		buf.WriteString("struct {")
		for i := 0; i < t.NumField(); i++ {
			buf.WriteRune(' ')
			f := t.Field(i)
			if !f.Anonymous() {
				buf.WriteString(f.Name())
				buf.WriteRune(' ')
			}
			buf.WriteString(NewT(f.Type()).String())
			if tag := f.Tag(); tag != "" {
				buf.WriteRune(' ')
				buf.WriteString("`" + string(tag) + "`")
			}
			if i == t.NumField()-1 {
				buf.WriteRune(' ')
			} else {
				buf.WriteRune(';')
			}
		}
		buf.WriteString("}")
		return buf.String()
	case reflect.Interface:
		if t.Name() == "error" {
			return "error"
		}
		buf := bytes.NewBuffer(nil)
		buf.WriteString("interface {")
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			buf.WriteRune(' ')
			if pkg := m.PkgPath(); pkg != "" {
				buf.WriteString(internal.NewPackage(pkg).Name())
				buf.WriteRune('.')
			}
			buf.WriteString(m.Name())
			buf.WriteString(m.Type().String()[4:])
			if i == t.NumMethod()-1 {
				buf.WriteRune(' ')
			} else {
				buf.WriteRune(';')
			}
		}
		buf.WriteString("}")
		return buf.String()
	case reflect.Func:
		buf := bytes.NewBuffer(nil)
		buf.WriteString("func(")
		for i := 0; i < t.NumIn(); i++ {
			p := t.In(i)
			if i == t.NumIn()-1 && t.IsVariadic() {
				buf.WriteString("...")
				buf.WriteString(NewT(p.Elem()).String())
			} else {
				buf.WriteString(NewT(p).String())
			}
			if i < t.NumIn()-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteRune(')')
		if t.NumOut() > 0 {
			buf.WriteRune(' ')
		}
		if t.NumOut() > 1 {
			buf.WriteRune('(')
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(NewT(t.Out(i)).String())
		}
		if t.NumOut() > 1 {
			buf.WriteRune(')')
		}
		return buf.String()
	default:
		return t.Name()
	}
}

func (v *T) ID() string {
	if pkg := v.t.PkgPath(); pkg != "" {
		return pkg + "." + v.t.Name()
	}
	return v.t.Name()
}
