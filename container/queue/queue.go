package queue

import (
	"sync"

	"github.com/xoctopus/x/container/list"
)

type Queue[T any] interface {
	// Len returns length of Queue
	Len() int
	// Push pushes element to end of Queue
	Push(v T)
	// Pop pops element from head of Queue
	Pop() (T, bool)
	// Head returns head element of Queue
	Head() (T, bool)
	// Tail returns tail element of Queue
	Tail() (T, bool)
	// Range calls f sequentially for each element present in the Queue.
	// If f returns false, range stops the iteration.
	Range(f func(T) bool)
	// Clear releases all the entries, resulting in an empty Queue.
	Clear()
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		List: list.New[T](),
	}
}

func NewSafeQueue[T any]() Queue[T] {
	return &queue[T]{
		mtx:  &sync.RWMutex{},
		List: list.New[T](),
	}
}

type queue[T any] struct {
	mtx *sync.RWMutex
	list.List[T]
}

func (q *queue[T]) Len() int {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	return q.List.Len()
}

func (q *queue[T]) Push(v T) {
	if q.mtx != nil {
		q.mtx.Lock()
		defer q.mtx.Unlock()
	}
	q.PushBack(v)
}

func (q *queue[T]) Pop() (T, bool) {
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

func (q *queue[T]) Head() (T, bool) {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	if e := q.List.Front(); e != nil {
		return e.Value, true
	}
	return *new(T), false
}

func (q *queue[T]) Tail() (T, bool) {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	if e := q.List.Back(); e != nil {
		return e.Value, true
	}
	return *new(T), false
}

func (q *queue[T]) Clear() {
	if q.mtx != nil {
		q.mtx.Lock()
		defer q.mtx.Unlock()
	}
	q.List.Clear()
}

func (q *queue[T]) Range(f func(T) bool) {
	if q.mtx != nil {
		q.mtx.RLock()
		defer q.mtx.RUnlock()
	}
	for e := q.Front(); e != nil; e = e.Next() {
		if !f(e.Value) {
			break
		}
	}
}
