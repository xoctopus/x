package stack

import (
	"sync"

	"github.com/xoctopus/x/container/list"
)

type Stack[T any] interface {
	// Len returns size of Stack
	Len() int
	// Push pushes element to head of Stack
	Push(v T)
	// Pop pops element from head of Stack
	Pop() (T, bool)
	// Top returns top element of Stack
	Top() (T, bool)
	// Clear releases all the entries, resulting in an empty Stack.
	Clear()
	// Range calls f sequentially for each element present in the Stack.
	// If f returns false, range stops the iteration.
	Range(func(T) bool)
}

func NewStack[T any]() Stack[T] {
	return &stack[T]{
		List: list.New[T](),
	}
}

func NewSafeStack[T any]() Stack[T] {
	return &stack[T]{
		mtx:  &sync.RWMutex{},
		List: list.New[T](),
	}
}

type stack[T any] struct {
	mtx *sync.RWMutex
	list.List[T]
}

func (s *stack[T]) Len() int {
	if s.mtx != nil {
		s.mtx.RLock()
		defer s.mtx.RUnlock()
	}
	return s.List.Len()
}

func (s *stack[T]) Push(v T) {
	if s.mtx != nil {
		s.mtx.Lock()
		defer s.mtx.Unlock()
	}
	s.PushFront(v)
}

func (s *stack[T]) Top() (T, bool) {
	if s.mtx != nil {
		s.mtx.RLock()
		defer s.mtx.RUnlock()
	}
	if e := s.List.Front(); e != nil {
		return e.Value, true
	}
	return *new(T), false
}

func (s *stack[T]) Pop() (T, bool) {
	if s.mtx != nil {
		s.mtx.Lock()
		defer s.mtx.Unlock()
	}
	if e := s.List.Front(); e != nil {
		s.List.Remove(e)
		return e.Value, true
	}
	return *new(T), false
}

func (s *stack[T]) Clear() {
	if s.mtx != nil {
		s.mtx.Lock()
		defer s.mtx.Unlock()
	}
	s.List.Clear()
}

func (s *stack[T]) Range(f func(T) bool) {
	if s.mtx != nil {
		s.mtx.RLock()
		defer s.mtx.RUnlock()
	}
	for e := s.Front(); e != nil; e = e.Next() {
		if !f(e.Value) {
			break
		}
	}
}
