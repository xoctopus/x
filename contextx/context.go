package contextx

import (
	"context"

	"github.com/pkg/errors"
)

type WithContext = func(ctx context.Context) context.Context

func WithContextCompose(withs ...WithContext) WithContext {
	return func(ctx context.Context) context.Context {
		for i := range withs {
			ctx = withs[i](ctx)
		}
		return ctx
	}
}

type (
	Option[T any] func(*ctx[T])
	Valuer[T any] func() T
)

func With[T any](v T) Option[T] {
	return func(c *ctx[T]) {
		c.valuer = func() T {
			return v
		}
	}
}

func WithValuer[T any](valuer Valuer[T]) Option[T] {
	return func(c *ctx[T]) {
		c.valuer = valuer
	}
}

func New[T any](options ...Option[T]) Context[T] {
	c := &ctx[T]{}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	return c
}

func NewValue[T any](v T) Context[T] {
	c := &ctx[T]{}
	With(v)(c)
	return c
}

type Context[T any] interface {
	With(context.Context, T) context.Context
	From(context.Context) (T, bool)
	MustFrom(context.Context) T
}

type ctx[T any] struct {
	valuer Valuer[T]
}

func (c *ctx[T]) With(ctx context.Context, v T) context.Context {
	return WithValue(ctx, c, v)
}

func (c *ctx[T]) MustFrom(ctx context.Context) T {
	if v, ok := ctx.Value(c).(T); ok {
		return v
	}
	if c.valuer != nil {
		return c.valuer()
	}
	panic(errors.Errorf("%T not found in context", c))
}

func (c *ctx[T]) From(ctx context.Context) (T, bool) {
	if v, ok := ctx.Value(c).(T); ok {
		return v, ok
	}
	return *new(T), false
}
