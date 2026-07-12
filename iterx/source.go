package iterx

import (
	"context"
	"iter"
	"slices"
)

func Of[V any, E ~[]V](values E) iter.Seq[V] {
	return slices.Values(values)
}

func SliceOf[T any, E ~[]T](values E) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range values {
			if !yield(i, v) {
				return
			}
		}
	}
}

// MapOf converts a map to iteration with key and value
func MapOf[K comparable, V any, M ~map[K]V](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func Recv[V any](c <-chan V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range c {
			if !yield(v) {
				return
			}
		}
	}
}

func RecvContext[V any](ctx context.Context, c <-chan V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok || !yield(v) {
					return
				}
			}
		}
	}
}
