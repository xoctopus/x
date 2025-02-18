package resultx

func Succeed[T any](v T) Result[T] {
	return Result[T]{&res[T, struct{ error }]{v: v, e: struct{ error }{nil}}}
}

func Err[T any](e error) Result[T] {
	return Result[T]{&res[T, struct{ error }]{e: struct{ error }{e}}}
}

func OK[T any](v T) ResultB[T] {
	return ResultB[T]{&res[T, bool]{v: v, e: true}}
}

func NotOK[T any]() ResultB[T] {
	return ResultB[T]{&res[T, bool]{e: false}}
}

func WrapResult[T any](v T, e error) Result[T] {
	return Result[T]{&res[T, struct{ error }]{v: v, e: struct{ error }{e}}}
}

func Unwrap[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func WrapResultB[T any](v T, b bool) ResultB[T] {
	return ResultB[T]{&res[T, bool]{v: v, e: b}}
}

func UnwrapB[T any](v T, b bool) T {
	if !b {
		panic(b)
	}
	return v
}

type Result[T any] struct {
	*res[T, struct{ error }]
}

type ResultB[T any] struct {
	*res[T, bool]
}

type R interface {
	bool | struct{ error }
}

type res[T any, E R] struct {
	v T
	e E
}

func (r *res[T, E]) succeed() bool {
	succeed := false
	switch e := any(r.e).(type) {
	case struct{ error }:
		succeed = e.error == nil
	case bool:
		succeed = e
	}
	return succeed
}

func (r *res[T, E]) Unwrap() T {
	if r.succeed() {
		return r.v
	}
	panic(r.e)
}

func (r *res[T, E]) UnwrapOr(v T) T {
	if r.succeed() {
		return r.v
	}
	return v
}

/*
func (r *res[T]) AndThen[U any](next func(T) result[U]) result[U] {
	if r.Failed() {
		return Err[U](r.e)
	}
	return next(r.v)
}
*/
