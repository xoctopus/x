package reflectx

import (
	"reflect"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

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
	must.BeTrueWrap(dst.CanSet(), "invalid dst value cannot set")
	must.BeTrueWrap(src.IsValid(), "invalid src value")
	must.BeTrueWrap(src.CanInterface(), "invalid src cannot read value")

	if dst.CanAddr() && src.CanAddr() && dst.Addr() == src.Addr() {
		return
	}

	tdst, tsrc := dst.Type(), src.Type()
	must.BeTrueWrap(tsrc.AssignableTo(tdst), "src type cannot assign to dst")

	if src.IsZero() {
		dst.Set(reflect.Zero(src.Type()))
		return
	}

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
		typ := src.Type()
		for i := range src.NumField() {
			// must.BeTrueWrap(
			// 		typ.Field(i).IsExported(),
			// 		"struct type `%s` field `%s` cannot be copied",
			// 		typ.String(), typ.Field(i).Name,
			// 	)
			if !typ.Field(i).IsExported() {
				if !dst.Field(i).CanAddr() || !src.Field(i).CanAddr() {
					panic(errors.Errorf("cannot be copied, unexported field `%s` of `%s` cannot addr", typ.Field(i).Name, typ))
				}
				ft := dst.Field(i).Type()
				dstv := reflect.NewAt(ft, unsafe.Pointer(dst.Field(i).UnsafeAddr())).Elem()
				srcv := reflect.NewAt(ft, unsafe.Pointer(src.Field(i).UnsafeAddr())).Elem()
				deepCopy(dstv, srcv, visited)
				continue
			}
			deepCopy(dst.Field(i), src.Field(i), visited)
		}
		// for i := range src.NumField() {
		// 	deepCopy(dst.Field(i), src.Field(i), visited)
		// }
	default: // basic kinds
		dst.Set(src)
	}
}

func HackFieldByName(v any, name string) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	rv = Indirect(rv)
	must.BeTrueWrap(rv.Kind() == reflect.Struct, "not a struct value")

	rv = rv.FieldByName(name)
	must.BeTrueWrap(rv.IsValid(), "field `%s` not found", name)
	must.BeTrueWrap(rv.CanAddr(), "field `%s` cannot addr", name)

	return reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
}

func HackField(v any, i int) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	rv = Indirect(rv)
	must.BeTrueWrap(rv.Kind() == reflect.Struct, "not a struct value")

	rv = rv.Field(i)
	must.BeTrueWrap(rv.CanAddr(), "field `%d` cannot addr", i)

	return reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
}
