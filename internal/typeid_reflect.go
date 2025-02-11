package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ReflectTypeID(t reflect.Type, wraps ...bool) string {
	wrap := len(wraps) > 0 && wraps[0]
	if id := t.Name(); id != "" {
		if pkg := t.PkgPath(); pkg != "" {
			id = pkg + "." + id
		}
		if wrap {
			return WrapNamedTypeID(id)
		}
		return id
	}

	switch t.Kind() {
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), ReflectTypeID(t.Elem(), wrap))
	case reflect.Chan:
		return fmt.Sprintf("%s%s", ChanDirPrefix(t), ReflectTypeID(t.Elem(), wrap))
	case reflect.Func:
		b := strings.Builder{}
		b.WriteString("func(")
		for i := range t.NumIn() {
			if i > 0 {
				b.WriteString(", ")
			}
			if i == t.NumIn()-1 && t.IsVariadic() {
				b.WriteString("...")
				b.WriteString(ReflectTypeID(t.In(i).Elem(), wrap))
				break
			}
			b.WriteString(ReflectTypeID(t.In(i), wrap))
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
				b.WriteString(ReflectTypeID(t.Out(i), wrap))
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
			b.WriteString(m.Name)
			b.WriteString(ReflectTypeID(m.Type, wrap)[4:])
		}
		b.WriteString(" }")
		return b.String()
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", ReflectTypeID(t.Key(), wrap), ReflectTypeID(t.Elem(), wrap))
	case reflect.Pointer:
		return fmt.Sprintf("*%s", ReflectTypeID(t.Elem(), wrap))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", ReflectTypeID(t.Elem(), wrap))
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
			b.WriteString(ReflectTypeID(f.Type, wrap))
			if f.Tag != "" {
				b.WriteString(" ")
				b.WriteString(strconv.Quote(string(f.Tag)))
			}
		}
		b.WriteString(" }")
		return b.String()
	default:
		panic(errors.Errorf("unexpected kind: %s", t.Kind()))
	}
}
