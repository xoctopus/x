package contextx

import (
	"context"

	"github.com/pkg/errors"
)

type (
	Option[T any] func(*ctx[T])
	Valuer[T any] func() T
	WithContext   func(context.Context) context.Context
)

func WithDefault[T any](v T) Option[T] {
	return func(c *ctx[T]) {
		c.defaulter = func() T { return v }
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
	WithCompose(v T) WithContext
}

type ctx[T any] struct {
	defaulter Valuer[T]
}

func (c *ctx[T]) With(ctx context.Context, v T) context.Context {
	return WithValue(ctx, c, v)
}

func (c *ctx[T]) From(ctx context.Context) (T, bool) {
	if v, ok := ctx.Value(c).(T); ok {
		return v, ok
	}
	if c.defaulter != nil {
		return c.defaulter(), true
	}
	return *new(T), false
}

func (c *ctx[T]) MustFrom(ctx context.Context) T {
	if v, ok := ctx.Value(c).(T); ok {
		return v
	}
	if c.defaulter != nil {
		return c.defaulter()
	}
	panic(errors.Errorf("%T not found in context", c))
}

func (c *ctx[T]) WithCompose(v T) WithContext {
	return func(ctx context.Context) context.Context {
		return WithValue(ctx, c, v)
	}
}

func Compose(withs ...WithContext) WithContext {
	return func(ctx context.Context) context.Context {
		for _, with := range withs {
			ctx = with(ctx)
		}
		return ctx
	}
}
