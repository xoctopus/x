package contextx

import (
	"context"

	"github.com/xoctopus/x/misc/must"
)

func With[K comparable, T any](ctx context.Context, v T) context.Context {
	k := *new(K)
	return context.WithValue(ctx, k, v)
}

func From[K comparable, T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(*new(K)).(T)
	return v, ok
}

func Must[K comparable, T any](ctx context.Context) T {
	return must.BeTrueV(From[K, T](ctx))
}

func Carry[K comparable, T any](v T) Carrier {
	return func(ctx context.Context) context.Context {
		return With[K, T](ctx, v)
	}
}
