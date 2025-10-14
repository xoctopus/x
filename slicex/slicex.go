package slicex

func Unique[T comparable, E ~[]T](s E) E {
	if len(s) <= 1 {
		return s
	}
	seen := make(map[T]struct{}, len(s))

	r := make(E, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		r = append(r, v)
	}
	return r
}
