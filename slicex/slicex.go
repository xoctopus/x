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

func UniqueM[SE any, RE comparable, SS ~[]SE](ss SS, m func(e SE) RE) []RE {
	dict := make(map[RE]struct{})

	for i := range ss {
		dict[m(ss[i])] = struct{}{}
	}

	return slices.Collect(maps.Keys(dict))
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
