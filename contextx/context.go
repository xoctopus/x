package contextx

import (
	"context"
	"sync"

	"github.com/xoctopus/x/misc/must"
)

type (
	Option[T any] func(*ctx[T])
	Valuer[T any] func() T

	// Deprecated: use Carrier instead
	WithContext func(context.Context) context.Context
)

func WithDefault[T any](v T) Option[T] {
	return func(c *ctx[T]) {
		c.defaulter = func() T { return v }
	}
}

func WithOnce[T any]() Option[T] {
	return func(c *ctx[T]) {
		c.once = true
	}
}

func NewT[T any](options ...Option[T]) Context[T] {
	c := &ctx[T]{}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	return c
}

func NewV[T any](v T) Context[T] {
	c := &ctx[T]{}
	WithDefault(v)(c)
	return c
}

type Context[T any] interface {
	With(context.Context, T) context.Context
	From(context.Context) (T, bool)
	MustFrom(context.Context) T
	Carry(v T) Carrier
}

type ctx[T any] struct {
	defaulter Valuer[T]
	once      bool
	with      func() context.Context
}

func (c *ctx[T]) With(ctx context.Context, v T) context.Context {
	if c.once {
		if c.with == nil {
			c.with = sync.OnceValue(func() context.Context {
				return WithValue(ctx, c, v)
			})
		}
		return c.with()
	}
	return WithValue(ctx, c, v)
}

func (c *ctx[T]) From(ctx context.Context) (T, bool) {
	if v, ok := ctx.Value(c).(T); ok {
		return v, ok
	}
	if c.defaulter != nil {
		return c.defaulter(), true
	}
	var zero T
	return zero, false
}

func (c *ctx[T]) MustFrom(ctx context.Context) T {
	v, ok := c.From(ctx)
	must.BeTrueF(ok, "%T not found in context", c)
	return v
}

func (c *ctx[T]) Carry(v T) Carrier {
	return func(ctx context.Context) context.Context {
		return c.With(ctx, v)
	}
}

type Carrier func(context.Context) context.Context

func Compose(carriers ...Carrier) Carrier {
	return func(ctx context.Context) context.Context {
		for _, carrier := range carriers {
			ctx = carrier(ctx)
		}
		return ctx
	}
}
