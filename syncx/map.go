package syncx

import (
	"sync"
)

type Map[K comparable, V any] interface {
	Store(K, V)

	Clear()
	Delete(K)
	CompareAndDelete(key K, old V) (deleted bool)

	Swap(k K, new V) (previous V, loaded bool)
	CompareAndSwap(key K, old, new V) (swapped bool)

	Load(K) (V, bool)
	LoadAndDelete(K) (v V, loaded bool)
	LoadOrStore(K, V) (actual V, loaded bool)
	Range(func(K, V) bool)

	Len() int
	Keys() []K
	Values() []V
}

func NewXmap[K comparable, V any]() Map[K, V] {
	return &xmap[K, V]{
		mtx: &sync.RWMutex{},
		m:   make(map[K]V),
	}
}

func AsXmap[K comparable, V any, M ~map[K]V](m M) Map[K, V] {
	return &xmap[K, V]{
		mtx: &sync.RWMutex{},
		m:   m,
	}
}

type xmap[K comparable, V any] struct {
	mtx *sync.RWMutex
	m   map[K]V
}

func (m *xmap[K, V]) Store(k K, v V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.m[k] = v
}

func (m *xmap[K, V]) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.m = make(map[K]V)
}

func (m *xmap[K, V]) Delete(k K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.m, k)
}

func (m *xmap[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if x, ok := m.m[k]; ok && any(x) == any(old) {
		delete(m.m, k)
		return true
	}
	return false
}

func (m *xmap[K, V]) Swap(k K, new V) (previous V, loaded bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	previous, loaded = m.m[k]
	m.m[k] = new
	return previous, loaded
}

func (m *xmap[K, V]) CompareAndSwap(k K, old, new V) (swapped bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if x, ok := m.m[k]; ok && any(x) == any(old) {
		m.m[k] = new
		return true
	}
	return false
}

func (m *xmap[K, V]) Load(k K) (v V, loaded bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	v, loaded = m.m[k]
	return
}

func (m *xmap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	v, loaded = m.m[k]
	if loaded {
		delete(m.m, k)
	}
	return
}

func (m *xmap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	actual, loaded = m.m[k]
	if loaded {
		return
	}
	m.m[k] = v
	return v, loaded
}

func (m *xmap[K, V]) Range(f func(K, V) bool) {
	m.mtx.RLock()
	m.mtx.RUnlock()
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

func (m *xmap[K, V]) Len() int {
	i := 0
	for _, _ = range m.Range {
		i++
	}
	return i
}

func (m *xmap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for k := range m.Range {
		keys = append(keys, k)
	}
	return keys
}

func (m *xmap[K, V]) Values() []V {
	values := make([]V, 0)
	for _, v := range m.Range {
		values = append(values, v)
	}
	return values
}

func NewSmap[K comparable, V any]() Map[K, V] {
	return &smap[K, V]{m: &sync.Map{}}
}

func AsSmap[K comparable, V any, M ~map[K]V](m M) Map[K, V] {
	sm := &sync.Map{}
	for k, v := range m {
		sm.Store(k, v)
	}
	return &smap[K, V]{m: sm}
}

type smap[K comparable, V any] struct {
	m *sync.Map
}

func (m *smap[K, V]) Store(k K, v V) {
	m.m.Store(k, v)
}

func (m *smap[K, V]) Clear() {
	m.m.Clear()
}

func (m *smap[K, V]) Delete(k K) {
	m.m.Delete(k)
}

func (m *smap[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	return m.m.CompareAndDelete(k, old)
}

func (m *smap[K, V]) Swap(k K, new V) (previous V, loaded bool) {
	return m.result(m.m.Swap(k, new))
}

func (m *smap[K, V]) CompareAndSwap(k K, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(k, old, new)
}

func (m *smap[K, V]) Load(k K) (v V, loaded bool) {
	return m.result(m.m.Load(k))
}

func (m *smap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	return m.result(m.m.LoadAndDelete(k))
}

func (m *smap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	return m.result(m.m.LoadOrStore(k, v))
}

func (m *smap[K, V]) Range(f func(K, V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *smap[K, V]) Len() int {
	i := 0
	for _, _ = range m.Range {
		i++
	}
	return i
}

func (m *smap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for k, _ := range m.Range {
		keys = append(keys, k)
	}
	return keys
}

func (m *smap[K, V]) Values() []V {
	values := make([]V, 0)
	for _, v := range m.Range {
		values = append(values, v)
	}
	return values
}

func (m *smap[K, V]) result(x any, b bool) (V, bool) {
	if v, ok := x.(V); ok {
		return v, b
	}
	return *new(V), b
}
