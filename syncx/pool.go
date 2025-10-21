package syncx

import (
	"reflect"
	"sync"
)

func NewPool[T any](newer func() T) Pool[T] {
	return Pool[T]{
		pool: sync.Pool{New: func() any { return newer() }},
	}
}

type Pool[T any] struct {
	pool sync.Pool
}

func (p *Pool[T]) Put(v T) {
	p.pool.Put(v)
}

func (p *Pool[T]) Get() T {
	v := p.pool.Get().(T)
	if r, ok := any(v).(interface{ Reset() }); ok {
		r.Reset()
	}

	if reflect.TypeFor[T]().Kind() == reflect.Slice {
		reflect.ValueOf(&v).Elem().SetLen(0)
	}
	return v
}
