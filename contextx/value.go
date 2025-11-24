package contextx

import (
	"context"
	"fmt"
	"reflect"
)

// WithValue like context.WithValue but faster🤏
func WithValue(parent context.Context, k, v any) context.Context {
	if parent == nil {
		panic("parent is nil")
	}
	if k == nil {
		panic("key is nil")
	}
	return &kv{parent, k, v}
}

// kv vs context.Context: when finding key
//
// context.Context
// +----------------+
// | ctx.Value(key) |
// +----------------+
//
//	|
//	V
//
// +---------------------+
// | Value method of kv? |  (kv from WithValue)
// +---------------------+
//
//	|
//	|-- key matches? --> return v
//	|-- key not match? --> ctx.parent.Value(key)
//	V
//
// traverse recursively up to root context.Background()
//
// kv
// +-------------------+
// | c.Value(key)      |
// +-------------------+
//
//	|
//	V
//
// +-------------------+
// | if c.k == key     |
// +-------------------+
//
//	|-- match --> return c.v
//	|-- no match --> c.Context.Value(key)
//	V
//
// traverse std.Context
type kv struct {
	context.Context
	k, v any
}

func (c *kv) String() string {
	return nameof(c.Context) +
		".WithValue(key:" + reflect.TypeOf(c.k).String() +
		", val:" + stringify(c.v) + ")"
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
