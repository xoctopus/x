package exp

import (
	"reflect"
	"unsafe"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
)

func Clone[T any](src T) (dst T) {
	s := reflect.ValueOf(src)
	d := reflect.New(s.Type()).Elem()

	deepCopy(d, s, make(map[uintptr]reflect.Value))

	return d.Interface().(T)
}

// DeepCopy performs a deep copy from the source reflect.Value `src` to the
// destination reflect.Value `dst`.
//
// It recursively copies all fields, including unexported ones, and supports
// structs, slices, maps, pointers, arrays, and interfaces. It does not support
// functions or channels.
func DeepCopy(dst, src reflect.Value) {
	deepCopy(dst, src, make(map[uintptr]reflect.Value))
}

func deepCopy(dst, src reflect.Value, visited map[uintptr]reflect.Value) {
	must.BeTrueWrap(dst.CanSet(), "invalid dst value cannot set")
	must.BeTrueWrap(src.IsValid(), "invalid src value")
	must.BeTrueWrap(src.CanInterface(), "invalid src cannot read value")

	// same value
	if dst.CanAddr() && src.CanAddr() && dst.Addr() == src.Addr() {
		return
	}

	// nil value
	if reflectx.CanNilValue(src) && reflectx.CanNilValue(dst) && src.IsNil() {
		dst.Set(reflect.Zero(src.Type()))
		return
	}

	tdst, tsrc := dst.Type(), src.Type()
	must.BeTrueWrap(tsrc.AssignableTo(tdst), "src type cannot assign to dst")

	if src.IsZero() {
		dst.Set(reflect.Zero(src.Type()))
		return
	}

	// TODO try skip circle reference
	// ptr := src.UnsafeAddr()
	// if val, ok := visited[ptr]; ok {
	// 	dst.Set(val)
	// 	return
	// }
	// visited[ptr] = dst

	// reflect.DeepEqual does not compare Chan type, and if Func type is not nil
	// it returns false anyway
	must.BeTrueWrap(src.Kind() != reflect.Func, "func type cannot be copied")
	must.BeTrueWrap(src.Kind() != reflect.Chan, "chan type cannot be copied")

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
		for i := range src.NumField() {
			if !tsrc.Field(i).IsExported() {
				must.BeTrueWrap(
					dst.Field(i).CanAddr() && src.Field(i).CanAddr(),
					"struct `%s` unexported field `%s` cannot be copied",
					tsrc.String(), tsrc.Field(i).Name,
				)
				deepCopy(HackField(dst, i), HackField(src, i), visited)
				continue
			}
			deepCopy(dst.Field(i), src.Field(i), visited)
		}
	default: // basic kinds
		dst.Set(src)
	}
}

// HackFieldByName returns an addressable `reflect.Value` for the specified field
// name of a struct, even if the field is unexported (private).
//
// This function uses unsafe.Pointer to bypass Go's visibility restrictions,
// allowing direct access to unexported struct fields. It should be used with
// caution, as it breaks encapsulation and may not be compatible with future
// versions of Go.
//
// Example:
//
//	type Foo struct { bar int }
//	f := Foo{bar: 42}
//	v := HackFieldByName(f, "bar")
//	v.SetInt(100)  // Successfully modifies the unexported field
func HackFieldByName(v any, name string) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	rv = reflectx.Indirect(rv)
	must.BeTrueWrap(rv.Kind() == reflect.Struct, "not a struct value")

	rv = rv.FieldByName(name)
	must.BeTrueWrap(rv.IsValid(), "field `%s` not found", name)
	must.BeTrueWrap(rv.CanAddr(), "field `%s` cannot addr", name)

	return reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
}

// HackField returns an addressable reflect.Value for the struct field at the
// given index, even if the field is unexported (private).
//
// Like HackFieldByName, this function uses unsafe.Pointer to bypass Go's
// visibility restrictions. It should be used with caution, as it breaks
// encapsulation and may be incompatible with future Go versions.
func HackField(v any, i int) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	rv = reflectx.Indirect(rv)
	must.BeTrueWrap(rv.Kind() == reflect.Struct, "not a struct value")

	rv = rv.Field(i)
	must.BeTrueWrap(rv.CanAddr(), "field `%d` cannot addr", i)

	return reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
}
