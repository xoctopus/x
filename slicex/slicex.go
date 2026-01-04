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

// UniqueKeys collects partials P from slices S by func m
func UniqueKeys[E any, S ~[]E, P comparable](s S, m func(E) P) []P {
	dict := make(map[P]struct{})

	for i := range s {
		dict[m(s[i])] = struct{}{}
	}

	return slices.Collect(maps.Keys(dict))
}

// UniqueValues collects unique []E from slices S by func m
func UniqueValues[E any, S ~[]E, K comparable](s S, m func(E) K) []E {
	dict := make(map[K]E)

	for i := range s {
		dict[m(s[i])] = s[i]
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
