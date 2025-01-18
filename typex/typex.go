package typex

import (
	"go/types"
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

// types.Type implements in go/types
var (
	_ types.Type = (*types.Alias)(nil)
	_ types.Type = (*types.Array)(nil)
	_ types.Type = (*types.Basic)(nil)
	_ types.Type = (*types.Chan)(nil)
	_ types.Type = (*types.Interface)(nil)
	_ types.Type = (*types.Map)(nil)
	_ types.Type = (*types.Named)(nil)
	_ types.Type = (*types.Pointer)(nil)
	_ types.Type = (*types.Slice)(nil)
	_ types.Type = (*types.Struct)(nil)
	_ types.Type = (*types.Tuple)(nil)
	_ types.Type = (*types.TypeParam)(nil)
	_ types.Type = (*types.Signature)(nil)
	_ types.Type = (*types.Union)(nil)
)
