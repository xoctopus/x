package typex

import (
	"bytes"
	"fmt"
	"reflect"
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

func Typename(t Type) string {
	if t == nil || t.Unwrap() == nil {
		return "nil"
	}

	buf := bytes.NewBuffer(nil)
	for t.Kind() == reflect.Pointer {
		buf.WriteByte('*')
		t = t.Elem()
	}

	if name := t.Name(); name != "" {
		if pkg := t.PkgPath(); pkg != "" {
			buf.WriteString(pkg)
			buf.WriteRune('.')
		}
		buf.WriteString(name)
		return buf.String()
	}
	buf.WriteString(t.String())
	return buf.String()
}

func Deref(t Type) Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func typename(t Type) string {
	if pkg := t.PkgPath(); pkg != "" {
		return pkg + "." + t.Name()
	}

	kind := t.Kind()
	if _, ok := ReflectKindToTypesKind[kind]; ok {
		return kind.String()
	}

	switch kind {
	case reflect.Slice:
		return "[]" + typename(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), typename(t.Elem()))
	case reflect.Chan:
		return "chan " + typename(t.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", typename(t.Key()), typename(t.Elem()))
	case reflect.Pointer:
		return "*" + typename(t.Elem())
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
			buf.WriteString(typename(f.Type()))
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
				buf.WriteString(NewPackage(pkg).Name())
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
				buf.WriteString(typename(p.Elem()))
			} else {
				buf.WriteString(typename(p))
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
			buf.WriteString(typename(t.Out(i)))
		}
		if t.NumOut() > 1 {
			buf.WriteRune(')')
		}
		return buf.String()
	default:
		return t.Name()
	}
}
