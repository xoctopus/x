package slicex

import (
	"maps"
	"slices"
)

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

// UniqueKeys collects partials P from slices S by func keyOfE
func UniqueKeys[E any, S ~[]E, K comparable](s S, keyOfE func(E) K) []K {
	dict := make(map[K]struct{})

	for i := range s {
		dict[keyOfE(s[i])] = struct{}{}
	}

	return slices.Collect(maps.Keys(dict))
}

// UniqueValues collects unique []E from slices S by func keyOfE
func UniqueValues[E any, S ~[]E, K comparable](s S, keyOfE func(E) K) []E {
	dict := make(map[K]E)

	for i := range s {
		dict[keyOfE(s[i])] = s[i]
	}

	return slices.Collect(maps.Values(dict))
}

// Equivalent compares two slices has same elements without order
func Equivalent[T comparable, E ~[]T](x, y E) bool {
	if len(x) != len(y) {
		return false
	}
	if len(x) == 0 {
		return (x == nil) == (y == nil)
	}

	marks := make(map[T]int)
	for _, v := range x {
		marks[v]++
	}
	for _, v := range y {
		if marks[v] == 0 {
			return false
		}
		marks[v]--
		if marks[v] == 0 {
			delete(marks, v)
		}
	}
	return len(marks) == 0
}

// M converts ~[]FROM to []TO by mapping function m
func M[FROM any, TO any, E ~[]FROM](s E, m func(FROM) TO) []TO {
	results := make([]TO, 0, len(s))
	for i := range s {
		results = append(results, m(s[i]))
	}
	return results
}

// Map converts ~[]E to map[K]v by mapping function m
func Map[E any, K comparable, V any, S ~[]E, R map[K]V](s S, m func(E) (K, V)) R {
	r := make(R)
	for i := range s {
		k, v := m(s[i])
		r[k] = v
	}
	return r
}
