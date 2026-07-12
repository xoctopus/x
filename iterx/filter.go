package iterx

import "iter"

// Filter returns a filtered iteration from seq by filter
func Filter[V any](seq iter.Seq[V], filter func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for e := range seq {
			if !filter(e) {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Map returns a mapped sequence
func Map[I, O any](seq iter.Seq[I], m func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for e := range seq {
			if !yield(m(e)) {
				return
			}
		}
	}
}

// MapFilter returns a mapped and filtered sequence
func MapFilter[I, O any](seq iter.Seq[I], filter func(I) (O, bool)) iter.Seq[O] {
	return func(yield func(O) bool) {
		for e := range seq {
			o, ok := filter(e)
			if !ok {
				continue
			}
			if !yield(o) {
				return
			}
		}
	}
}
