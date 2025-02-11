package typex

import (
	"go/types"
	"reflect"

	"github.com/xoctopus/x/internal"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
)

func NewRType(t any) *RType {
	tt, ok := t.(reflect.Type)
	if !ok {
		tt = reflect.TypeOf(t)
	}
	must.NotNilWrap(tt, "invalid reflect.Type `t`")
	return &RType{t: tt}
}

// RType wraps reflect.Type and implements typex.Type
type RType struct {
	t reflect.Type
}

func (t *RType) Unwrap() any { return t.t }

func (t *RType) PkgPath() string { return t.t.PkgPath() }

func (t *RType) Name() string { return t.t.Name() }

func (t *RType) String() string { return internal.ReflectTypeID(t.t) }

func (t *RType) Kind() reflect.Kind { return t.t.Kind() }

func (t *RType) Implements(u Type) bool {
	switch x := u.(type) {
	case *RType:
		if u.Kind() == reflect.Interface {
			return t.t.Implements(x.t)
		}
		return false
	case *GType:
		if i, ok := x.t.Underlying().(*types.Interface); ok {
			return types.Implements(internal.NewTypesType(t.t), i)
		}
		return false
	default:
		return false
	}
}

func (t *RType) AssignableTo(u Type) bool {
	switch x := u.(type) {
	case *RType:
		return t.t.AssignableTo(x.t)
	default:
		return false
	}
}

func (t *RType) ConvertibleTo(u Type) bool {
	if x, ok := u.(*RType); ok {
		return t.t.ConvertibleTo(x.t)
	}
	return false
}

func (t *RType) Comparable() bool {
	return t.t.Comparable()
}

func (t *RType) Key() Type {
	if t.Kind() == reflect.Map {
		return NewRType(t.t.Key())
	}
	return nil
}

func (t *RType) Elem() Type {
	if reflectx.CanElem(t.t.Kind()) {
		return NewRType(t.t.Elem())
	}
	return nil
}

func (t *RType) Len() int {
	if t.t.Kind() == reflect.Array {
		return t.t.Len()
	}
	return 0
}

func (t *RType) NumField() int {
	if t.Kind() == reflect.Struct {
		return t.t.NumField()
	}
	return 0
}

func (t *RType) Field(i int) StructField {
	if t.Kind() == reflect.Struct && i < t.t.NumField() {
		return &RStructField{sf: t.t.Field(i)}
	}
	return nil
}

func (t *RType) FieldByName(name string) (StructField, bool) {
	if t.Kind() == reflect.Struct {
		if sf, ok := t.t.FieldByName(name); ok {
			return &RStructField{sf: sf}, true
		}
	}
	return nil, false
}

func (t *RType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	if t.Kind() == reflect.Struct {
		if sf, ok := t.t.FieldByNameFunc(match); ok {
			return &RStructField{sf: sf}, true
		}
	}
	return nil, false
}

func (t *RType) NumMethod() int {
	return t.t.NumMethod()
}

func (t *RType) Method(i int) Method {
	if i >= 0 && i < t.NumMethod() {
		return &RMethod{t.t.Method(i)}
	}
	return nil
}

func (t *RType) MethodByName(name string) (Method, bool) {
	if m, ok := t.t.MethodByName(name); ok {
		return &RMethod{m: m}, true
	}
	return nil, false
}

func (t *RType) IsVariadic() bool {
	if t.Kind() == reflect.Func {
		return t.t.IsVariadic()
	}
	return false
}

func (t *RType) NumIn() int {
	if t.Kind() == reflect.Func {
		return t.t.NumIn()
	}
	return 0
}

func (t *RType) In(i int) Type {
	if i < t.NumIn() {
		return NewRType(t.t.In(i))
	}
	return nil
}

func (t *RType) NumOut() int {
	if t.Kind() == reflect.Func {
		return t.t.NumOut()
	}
	return 0
}

func (t *RType) Out(i int) Type {
	if i < t.NumOut() {
		return NewRType(t.t.Out(i))
	}
	return nil
}

type RStructField struct {
	sf reflect.StructField
}

func (f *RStructField) PkgPath() string {
	return f.sf.PkgPath
}

func (f *RStructField) Name() string {
	return f.sf.Name
}

func (f *RStructField) Type() Type {
	return NewRType(f.sf.Type)
}

func (f *RStructField) Tag() reflect.StructTag {
	return f.sf.Tag
}

func (f *RStructField) Anonymous() bool {
	return f.sf.Anonymous
}

type RMethod struct {
	m reflect.Method
}

func (m *RMethod) PkgPath() string {
	return m.m.PkgPath
}

func (m *RMethod) Name() string {
	return m.m.Name
}

func (m *RMethod) Type() Type {
	return NewRType(m.m.Type)
}
