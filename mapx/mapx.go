package mapx

import (
	"sync"

	"github.com/xoctopus/x/misc/must"
)

type U[K comparable, V any] interface {
	xmap[K, V] | *sync.Map
}

func NewXmap[K comparable, V any](m map[K]V, safe bool) Map[K, V] {
	xm := &xmap[K, V]{m: m}
	if safe {
		xm.mtx = &sync.RWMutex{}
	}
	return xm
}

type xmap[K comparable, V any] struct {
	m   map[K]V
	mtx *sync.RWMutex
}

func (m *xmap[K, V]) Exists(k K) bool {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	_, ok := m.m[k]
	return ok
}

func (m xmap[K, V]) Load(k K) (v V, loaded bool) {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	v, loaded = m.m[k]
	return
}

func (m xmap[K, V]) LoadEq(k K, eq func(any, any) bool) (value V, loaded bool) {
	m.Range(func(key K, val V) bool {
		if eq(key, k) {
			value, loaded = val, true
			return false
		}
		return true
	})
	return
}

func (m xmap[K, V]) LoadEqs(k K, eq func(any, any) bool) (values []V) {
	m.Range(func(key K, value V) bool {
		if eq(key, k) {
			values = append(values, value)
		}
		return true
	})
	return values
}

func (m xmap[K, V]) BatchLoad(ks ...K) []Result[V] {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	results := make([]Result[V], 0, len(ks))
	for i := range ks {
		vi, ok := m.m[ks[i]]
		results = append(results, Result[V]{vi, ok})
	}
	return results
}

func (m xmap[K, V]) Store(k K, v V) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	m.m[k] = v
}

func (m xmap[K, V]) BatchStore(ks []K, vs []V) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	must.BeTrue(len(ks) == len(vs))
	for i := range len(ks) {
		m.m[ks[i]] = vs[i]
	}
}

func (m xmap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	actual, loaded = m.m[k]
	if loaded {
		return
	}
	m.m[k] = v
	actual = v
	return
}

func (m xmap[K, V]) LoadAndDelete(k K) (value V, loaded bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	if value, loaded = m.m[k]; loaded {
		delete(m.m, k)
	}
	return
}

func (m xmap[K, V]) Delete(k K) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	delete(m.m, k)
}

func (m xmap[K, V]) BatchDelete(ks ...K) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	for i := range ks {
		delete(m.m, ks[i])
	}
}

func (m xmap[K, V]) Swap(k K, v V) (previous V, loaded bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	previous, loaded = m.m[k]
	m.m[k] = v
	return
}

func (m xmap[K, V]) CompareAndSwap(k K, old, new V) (swapped bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	if v, ok := m.m[k]; ok && any(v) == any(old) {
		m.m[k] = new
		return true
	}
	return
}

func (m xmap[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	if v, ok := m.m[k]; ok && any(v) == any(old) {
		delete(m.m, k)
		return true
	}
	return
}

func (m xmap[K, V]) Range(f func(K, V) bool) {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

func (m xmap[K, V]) Clear() {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	clear(m.m)
}

func NewSmap(m *sync.Map) Map[any, any] {
	return &smap{m}
}

type smap struct {
	*sync.Map
}

func (s *smap) Exists(k any) bool {
	_, exists := s.Load(k)
	return exists
}

func (s *smap) LoadEq(k any, eq func(any, any) bool) (value any, loaded bool) {
	s.Range(func(key, val any) bool {
		if eq(key, k) {
			value, loaded = val, true
			return false
		}
		return true
	})
	return
}

func (s *smap) LoadEqs(k any, eq func(any, any) bool) (values []any) {
	s.Range(func(key, value any) bool {
		if eq(key, k) {
			values = append(values, value)
		}
		return true
	})
	return values
}

func (s *smap) BatchLoad(ks ...any) []Result[any] {
	results := make([]Result[any], 0, len(ks))
	for i := range ks {
		vi, loaded := s.Load(ks[i])
		results = append(results, Result[any]{vi, loaded})
	}
	return results
}

func (s *smap) BatchStore(ks []any, vs []any) {
	must.BeTrue(len(ks) == len(vs))
	for i := range ks {
		s.Store(ks[i], vs[i])
	}
}

func (s *smap) BatchDelete(ks ...any) {
	for i := range ks {
		s.Delete(ks[i])
	}
}

type Result[V any] struct {
	V      V
	Loaded bool
}

type Map[K comparable, V any] interface {
	Exists(K) bool
	Load(K) (V, bool)
	LoadEq(K, func(any, any) bool) (V, bool)
	LoadEqs(K, func(any, any) bool) []V
	BatchLoad(...K) []Result[V]
	Store(K, V)
	BatchStore([]K, []V)
	LoadOrStore(K, V) (actual V, loaded bool)
	LoadAndDelete(K) (value V, loaded bool)
	Delete(K)
	BatchDelete(...K)
	Swap(K, V) (previous V, loaded bool)
	CompareAndSwap(key K, old, new V) (swapped bool)
	CompareAndDelete(key K, old V) (deleted bool)
	Range(func(key K, value V) (shouldContinue bool))
	Clear()
}

func Keys[K comparable, V any](m Map[K, V]) []K {
	keys := make([]K, 0)
	m.Range(func(key K, _ V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

func Values[K comparable, V any](m Map[K, V]) []V {
	values := make([]V, 0)
	m.Range(func(_ K, value V) bool {
		values = append(values, value)
		return true
	})
	return values
}

func Len[K comparable, V any](m Map[K, V]) int {
	size := 0
	m.Range(func(key K, value V) bool {
		size++
		return true
	})
	return size
}
