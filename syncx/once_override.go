package syncx

import "sync/atomic"

func NewOnceOverride[T any](defaultValue T) *OnceOverride[T] {
	return &OnceOverride[T]{v: defaultValue}
}

// OnceOverride holds a value that has an initial default and can be overridden
// once. After the first successful Set, further modifications are not allowed.
type OnceOverride[T any] struct {
	v    T
	done atomic.Bool
}

func (o *OnceOverride[T]) Set(x T) bool {
	if o.done.CompareAndSwap(false, true) {
		o.v = x
		return true
	}
	return false
}

func (o *OnceOverride[T]) Value() T {
	return o.v
}
