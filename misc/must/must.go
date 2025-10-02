package must

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func NoError(err error) {
	if err != nil {
		panic(errors.New(err.Error()))
	}
}

func NoErrorF(err error, msg string, args ...any) {
	if err != nil {
		panic(errors.Wrapf(err, msg, args...))
	}
}

func NoErrorV[V any](v V, err error) V {
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return v
}

func BeTrue(ok bool) {
	if !ok {
		panic(errors.New("must be true"))
	}
}

func BeTrueF(b bool, msg string, args ...any) {
	if !b {
		panic(errors.Errorf("must be true: "+msg, args...))
	}
}

func BeTrueWrap(b bool, w error) {
	BeTrueF(b, "[err: %v]", w)
}

func BeTrueWrapF(b bool, w error, msg string, args ...any) {
	BeTrueF(b, msg+fmt.Sprintf(" [err: %v]", w), args...)
}

func BeTrueV[V any](v V, ok bool) V {
	if !ok {
		panic(errors.New("must be true"))
	}
	return v
}

func NotNilV[V any](v V) V {
	NotNilF(any(v), "")
	return v
}

func NotNilF(v any, msg string, args ...any) {
	format := ""
	rv := reflect.ValueOf(v)
	switch kind := rv.Kind(); kind {
	default:
		return
	case reflect.Invalid:
		format = "must not nil, but got invalid value."
		goto Panic
	case reflect.Chan, reflect.Func, reflect.Pointer, reflect.UnsafePointer,
		reflect.Slice, reflect.Map:
		if rv.IsNil() {
			format = "must not nil for type `%s`."
			args = append([]any{rv.Type()}, args...)
			goto Panic
		}
		return
	}
Panic:
	if msg != "" {
		format = format + " " + msg
	}
	err := errors.Errorf(format, args...)
	panic(err)
}

func NotNilWrap(v any, w error) {
	NotNilF(v, "[err: %v]", w)
}

func NotNilWrapF(v any, w error, msg string, args ...any) {
	NotNilF(v, msg+fmt.Sprintf(" [err: %v]", w), args...)
}

func IdenticalTypes(t1, t2 any) bool {
	rt1, ok := t1.(reflect.Type)
	if !ok {
		rt1 = reflect.TypeOf(t1)
	}
	rt2, ok := t2.(reflect.Type)
	if !ok {
		rt2 = reflect.TypeOf(t2)
	}
	return rt1 == rt2
}
