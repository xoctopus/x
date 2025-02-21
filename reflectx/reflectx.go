package reflectx

import (
	"bytes"
	"reflect"
)

var (
	InvalidValue = reflect.Value{}
	InvalidType  = reflect.TypeOf(nil)
)

// Indirect deref all level pointer references
func Indirect(v any) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if !rv.IsValid() {
		return InvalidValue
	}

	if rv.Kind() == reflect.Pointer {
		return Indirect(rv.Elem())
	}

	if rv.Kind() == reflect.Interface {
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

// IsZero check if input v is zero
func IsZero(v any) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	kind := rv.Kind()
	if !rv.IsValid() {
		return true
	}
	if kind == reflect.Pointer || kind == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		return IsZero(rv.Elem())
	}

	if rv.CanInterface() {
		if checker, ok := rv.Interface().(interface{ IsZero() bool }); ok {
			return checker.IsZero()
		}
	}

	// check if a CanLen value's length is 0(not include Array type)
	switch kind {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		return rv.Len() == 0
	default:
		return rv.IsZero()
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

func Clone[T any](src T) (dst T) {
	s := reflect.ValueOf(src)
	d := reflect.New(s.Type()).Elem()

	deepCopy(d, s, make(map[uintptr]reflect.Value))

	return d.Interface().(T)
}

func DeepCopy(v, x reflect.Value) {
	deepCopy(v, x, make(map[uintptr]reflect.Value))
}

func deepCopy(dst, src reflect.Value, visited map[uintptr]reflect.Value) {
	mustBeTrue(dst.CanSet(), "invalid dst value cannot set")
	mustBeTrue(src.IsValid(), "invalid src value")

	if dst.CanAddr() && src.CanAddr() && dst.Addr() == src.Addr() {
		return
	}

	tdst, tsrc := dst.Type(), src.Type()
	mustBeTrue(tsrc.AssignableTo(tdst), "src type cannot assign to dst")

	if src.IsZero() {
		dst.Set(reflect.Zero(src.Type()))
		return
	}

	// reflect.DeepEqual does not compare Chan type, and if Func type is not nil
	// it returns false anyway
	mustBeTrue(src.Kind() != reflect.Func, "func type cannot be copied")
	mustBeTrue(src.Kind() != reflect.Chan, "chan type cannot be copied")

	// if src.CanAddr() {
	// 	unsafe.Pointer()
	// 	if vv, ok := visited[src.Addr().Pointer()]; ok {
	// 		dst.Set(vv)
	// 		return
	// 	}
	// 	visited[src.Pointer()] = dst
	// }

	// dst.Set(reflect.New(src.Type()).Elem())
	switch src.Kind() {
	case reflect.Array:
		dst.Set(reflect.New(src.Type()).Elem())
		for i := range src.Len() {
			deepCopy(dst.Index(i), src.Index(i), visited)
		}
	case reflect.Interface:
		val := reflect.New(src.Elem().Type()).Elem()
		deepCopy(val, src.Elem(), visited)
		dst.Set(val)
	case reflect.Map:
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		for _, key := range src.MapKeys() {
			val := reflect.New(src.MapIndex(key).Type()).Elem()
			deepCopy(val, src.MapIndex(key), visited)
			dst.SetMapIndex(key, val)
		}
	case reflect.Pointer:
		val := reflect.New(src.Elem().Type()).Elem()
		deepCopy(val, src.Elem(), visited)
		dst.Set(val.Addr())
	case reflect.Slice:
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := range src.Len() {
			val := reflect.New(src.Index(i).Type()).Elem()
			deepCopy(val, src.Index(i), visited)
			dst.Index(i).Set(val)
		}
	case reflect.Struct:
		typ := src.Type()
		for i := range src.NumField() {
			mustBeTrue(
				typ.Field(i).IsExported(),
				"struct type `%s` field `%s` cannot be copied",
				typ.String(), typ.Field(i).Name,
			)
		}
		for i := range src.NumField() {
			deepCopy(dst.Field(i), src.Field(i), visited)
		}
	default: // basic kinds
		dst.Set(src)
	}
}
