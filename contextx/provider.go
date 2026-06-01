package contextx

import (
	"context"

	"github.com/xoctopus/x/misc/must"
)

func With[KT comparable, P any](ctx context.Context, v P) context.Context {
	return context.WithValue(ctx, *new(KT), v)
}

func From[KT comparable, P any](ctx context.Context) (P, bool) {
	v, ok := ctx.Value(*new(KT)).(P)
	return v, ok
}

func FromOr[KT comparable, P any](ctx context.Context, or P) P {
	if x, ok := From[KT, P](ctx); ok {
		return x
	}
	return or
}

func Must[KT comparable, P any](ctx context.Context) P {
	return must.BeTrueV(From[KT, P](ctx))
}

func Carry[KT comparable, P any](v P) Carrier {
	return func(ctx context.Context) context.Context {
		return With[KT, P](ctx, v)
	}
}

type Provider interface {
	WithContext(context.Context) context.Context
}
