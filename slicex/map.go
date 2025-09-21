package slicex

func Map[E any, T any](list []E, f func(E) T) []T {
	out := make([]T, len(list))
	for i := range list {
		out[i] = f(list[i])
	}
	return out
}
