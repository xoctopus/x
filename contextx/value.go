package contextx

import (
	"context"
	"fmt"
	"reflect"
)

// WithValue like context.WithValue but faster
func WithValue(parent context.Context, k, v any) context.Context {
	if parent == nil {
		panic("parent is nil")
	}
	if k == nil {
		panic("key is nil")
	}
	return &kv{parent, k, v}
}

type kv struct {
	context.Context
	k, v any
}

func (c *kv) String() string {
	return nameof(c.Context) +
		".WithValue(type " + reflect.TypeOf(c.k).String() +
		", val" + stringify(c.v) + ")"
}

func (c *kv) Value(k any) any {
	if c.k == k {
		return c.v
	}
	return c.Context.Value(k)
}

func nameof(c context.Context) string {
	if str, ok := c.(fmt.Stringer); ok {
		return str.String()
	}
	return reflect.TypeOf(c).String()
}

func stringify(v any) string {
	switch s := v.(type) {
	case fmt.Stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}
