package stack

import (
	"sync"

	"github.com/xoctopus/x/container/list"
)

type Stack[T any] interface {
	Len() int
	Push(v T)
	Pop() (T, bool)
	Top() (T, bool)
	Clear()
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

func (q *stack[T]) Len() int {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	return q.List.Len()
}

func (q *stack[T]) Push(v T) {
	if q.mtx != nil {
		q.mtx.Lock()
		defer q.mtx.Unlock()
	}
	q.PushFront(v)
}

func (q *stack[T]) Top() (T, bool) {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	if e := q.List.Front(); e != nil {
		return e.Value, true
	}
	return *new(T), false
}

func (q *stack[T]) Pop() (T, bool) {
	if q.mtx != nil {
		q.mtx.Lock()
		defer q.mtx.Unlock()
	}
	if e := q.List.Front(); e != nil {
		q.List.Remove(e)
		return e.Value, true
	}
	return *new(T), false
}

func (q *stack[T]) Clear() {
	if q.mtx != nil {
		q.mtx.Lock()
		defer q.mtx.Unlock()
	}
	q.List.Clear()
}
