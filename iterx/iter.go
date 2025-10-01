package iterx

import "iter"

// Map returns a mapped iteration from seq by m
func Map[I any, O any](seq iter.Seq[I], m func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for e := range seq {
			if !yield(m(e)) {
				return
			}
		}
	}
}

// MapSlice returns a mapped iteration from values by m
func MapSlice[I any, E ~[]I, O any](values E, m func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for _, e := range values {
			if !yield(m(e)) {
				return
			}
		}
	}
}

// Filter returns a filtered iteration from seq by filter
func Filter[V any](seq iter.Seq[V], filter func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for e := range seq {
			if filter(e) && !yield(e) {
				return
			}
		}
	}
}

// FilterSlice returns a filtered iteration from values by filter
func FilterSlice[V any, E ~[]V](values E, filter func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, e := range values {
			if filter(e) && !yield(e) {
				return
			}
		}
	}
}

// Slice converts a slice to iteration
func Slice[T any, E ~[]T](values E) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// SliceSeq converts a slice to iteration with index yielded
func SliceSeq[T any, E ~[]T](values E) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range values {
			if !yield(i, v) {
				return
			}
		}
	}
}

// MapSeq converts a map to iteration with key and value
func MapSeq[K comparable, V any, M ~map[K]V](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func Values[T any](seq iter.Seq[T]) []T {
	values := make([]T, 0)
	for v := range seq {
		values = append(values, v)
	}
	return values
}
