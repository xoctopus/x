package must

import (
	"github.com/pkg/errors"

	"github.com/xoctopus/x/reflectx"
)

func NoError(err error) {
	if err != nil {
		panic(errors.New(err.Error()))
	}
}

func NoErrorWrap(err error, msg string, args ...any) {
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

func BeTrueWrap(ok bool, msg string, args ...any) {
	if !ok {
		panic(errors.Errorf("must be true: "+msg, args...))
	}
}

func BeTrueV[V any](v V, ok bool) V {
	if !ok {
		panic(errors.New("must be true"))
	}
	return v
}

func NotNilV[V any](v V) V {
	rv := reflectx.Indirect(v)
	if rv == reflectx.InvalidValue {
		panic("must not nil: invalid value")
	}
	return v
}
