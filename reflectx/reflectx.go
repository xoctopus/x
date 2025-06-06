package reflectx

import (
	"bytes"
	"reflect"
)

var (
	InvalidValue = reflect.Value{}
	InvalidType  = reflect.TypeOf(nil)
)

type ZeroChecker interface {
	IsZero() bool
}

var TypeZeroChecker = reflect.TypeFor[ZeroChecker]()

// Indirect deref all level pointer references
func Indirect(v any) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if !rv.IsValid() {
		return InvalidValue
	}

	if (rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer) &&
		rv.Type().Name() == "" {
		return Indirect(rv.Elem())
	}

	return rv
}

// IndirectNew returns the indirect value of v
// this function is safe and WILL NOT trigger panic. if the input is invalid,
// InvalidValue returns. validation of return is recommended.
func IndirectNew(v any) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if !rv.IsValid() {
		return InvalidValue
	}

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() && rv.CanSet() {
			rv.Set(New(rv.Type()))
		}
		return IndirectNew(rv.Elem())
	}

	if rv.Kind() == reflect.Interface {
		return IndirectNew(rv.Elem())
	}
	return rv
}

// Deref returns the basic type of t
func Deref(t reflect.Type) reflect.Type {
	if t == InvalidType {
		return InvalidType
	}
	if kind := t.Kind(); kind == reflect.Pointer || kind == reflect.Interface {
		return Deref(t.Elem())
	}
	return t
}

// New a `reflect.Value` with reflect.Type
// not like reflect.New, but new all level pointer ref
func New(t reflect.Type) reflect.Value {
	v := reflect.New(t).Elem()
	if t.Kind() == reflect.Pointer {
		v.Set(New(t.Elem()).Addr())
	}
	return v
}

// NewElem new the indirect type of t
func NewElem(t reflect.Type) reflect.Value {
	ptrLevel := 0
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		ptrLevel++
	}

	rv := reflect.New(t)

	for i := 0; i < ptrLevel-1; i++ {
		t = reflect.PointerTo(t)
		nextrv := reflect.New(t)
		nextrv.Elem().Set(rv)
		rv = nextrv
	}

	return rv.Elem()
}

// IsZero checks whether the given value is zero or its underlying value is zero
//
// If the value implements the ZeroChecker interface, IsZero will use its IsZero
// method to determine zero-ness. Special handling is provided for slices, maps,
// strings, and channels, which are considered zero if their length is zero.
func IsZero(v any) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if !rv.IsValid() {
		return true
	}

	for kind := rv.Kind(); kind == reflect.Pointer || kind == reflect.Interface; {
		if !rv.IsValid() || rv.IsNil() {
			return true
		}
		rv = rv.Elem()
		kind = rv.Kind()
	}

	if rv.Type().Implements(TypeZeroChecker) {
		return rv.Interface().(ZeroChecker).IsZero()
	}

	if rv.CanAddr() && reflect.PointerTo(rv.Type()).Implements(TypeZeroChecker) {
		return rv.Addr().Interface().(ZeroChecker).IsZero()
	}

	if !rv.IsValid() || rv.IsZero() {
		return true
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		return rv.Len() == 0
	default:
		return false
	}
}

// Typename returns the full type name of rt
func Typename(rt reflect.Type) string {
	buf := bytes.NewBuffer(nil)
	for rt.Kind() == reflect.Ptr {
		buf.WriteByte('*')
		rt = rt.Elem()
	}

	if name := rt.Name(); name != "" {
		if pkg := rt.PkgPath(); pkg != "" {
			buf.WriteString(pkg)
			buf.WriteByte('.')
		}
		buf.WriteString(name)
		return buf.String()
	}

	buf.WriteString(rt.String())
	return buf.String()
}

func IsBytes(v any) bool {
	if _, ok := v.([]byte); ok {
		return true
	}
	t := typeof(v)
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 && t.Elem().PkgPath() == ""
}

func IsInteger(v any) bool {
	k := typeof(v).Kind()
	return k >= reflect.Int && k <= reflect.Uint64
}

func IsFloat(v any) bool {
	k := typeof(v).Kind()
	return k == reflect.Float64 || k == reflect.Float32
}

func IsNumeric(v any) bool {
	k := typeof(v).Kind()
	return k >= reflect.Int && k <= reflect.Complex128
}

func typeof(v any) reflect.Type {
	switch x := v.(type) {
	case reflect.Type:
		return x
	case interface{ Type() reflect.Type }:
		return x.Type()
	default:
		return reflect.TypeOf(v)
	}
}

func CanElem(k reflect.Kind) bool {
	return k == reflect.Chan || k == reflect.Pointer || k == reflect.Map ||
		k == reflect.Slice || k == reflect.Array
}
