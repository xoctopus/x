package typex

import (
	"go/types"
	"reflect"
)

func NewReflectType(t reflect.Type) *ReflectType {
	return &ReflectType{Type: t}
}

type ReflectType struct {
	reflect.Type
}

var _ Type = (*ReflectType)(nil)

func (t *ReflectType) Unwrap() any {
	return t.Type
}

func (t *ReflectType) String() string {
	return typename(t)
}

func (t *ReflectType) Implements(u Type) bool {
	switch x := u.(type) {
	case *ReflectType:
		return t.Type.Implements(x.Type)
	case *GoType:
		if t.PkgPath() == "" {
			return false
		}
		if v, ok := x.Type.(*types.Interface); ok {
			return types.Implements(NewGoTypeFromReflectType(t.Type), v)
		}
	}
	return false
}

func (t *ReflectType) AssignableTo(u Type) bool {
	if x, ok := u.(*ReflectType); ok {
		return t.Type.AssignableTo(x.Type)
	}
	return false
}

func (t *ReflectType) ConvertibleTo(u Type) bool {
	if x, ok := u.(*ReflectType); ok {
		return t.Type.ConvertibleTo(x.Type)
	}
	return false
}

func (t *ReflectType) Key() Type {
	return NewReflectType(t.Type.Key())
}

func (t *ReflectType) Elem() Type {
	return NewReflectType(t.Type.Elem())
}

func (t *ReflectType) Field(i int) StructField {
	return &ReflectStructField{
		StructField: t.Type.Field(i),
	}
}

func (t *ReflectType) FieldByName(name string) (StructField, bool) {
	if f, ok := t.Type.FieldByName(name); ok {
		return &ReflectStructField{
			StructField: f,
		}, true
	}
	return nil, false
}

func (t *ReflectType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	if f, ok := t.Type.FieldByNameFunc(match); ok {
		return &ReflectStructField{
			StructField: f,
		}, true
	}
	return nil, false
}

func (t *ReflectType) Method(i int) Method {
	return &ReflectMethod{
		Method: t.Type.Method(i),
	}
}

func (t *ReflectType) MethodByName(name string) (Method, bool) {
	if m, ok := t.Type.MethodByName(name); ok {
		return &ReflectMethod{m}, true
	}
	return nil, false
}

func (t *ReflectType) In(i int) Type {
	return NewReflectType(t.Type.In(i))
}

func (t *ReflectType) Out(i int) Type {
	return NewReflectType(t.Type.Out(i))
}

type ReflectMethod struct {
	Method reflect.Method
}

var _ Method = (*ReflectMethod)(nil)

func (m *ReflectMethod) PkgPath() string {
	return m.Method.PkgPath
}

func (m *ReflectMethod) Name() string {
	return m.Method.Name
}

func (m *ReflectMethod) Type() Type {
	return NewReflectType(m.Method.Type)
}

type ReflectStructField struct {
	StructField reflect.StructField
}

var _ StructField = (*ReflectStructField)(nil)

func (f *ReflectStructField) PkgPath() string {
	return f.StructField.PkgPath
}

func (f *ReflectStructField) Name() string {
	return f.StructField.Name
}

func (f *ReflectStructField) Tag() reflect.StructTag {
	return f.StructField.Tag
}

func (f *ReflectStructField) Type() Type {
	return NewReflectType(f.StructField.Type)
}

func (f *ReflectStructField) Anonymous() bool {
	return f.StructField.Anonymous
}
