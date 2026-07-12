package iterx

import "iter"

func Flatten[T any](seq iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seq(func(s iter.Seq[T]) bool {
			cont := true
			s(func(v T) bool {
				cont = yield(v)
				return cont
			})
			return cont
		})
	}
}
