package mapx

import (
	"reflect"
	"sync"

	"github.com/xoctopus/x/misc/must"
)

type Map[K comparable, V any] interface {
	Exists(K) bool
	Load(K) (V, bool)
	LoadEq(func(any) bool) (V, bool)
	LoadEqs(func(any) bool) []V
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
	Clone() Map[K, V]
}

func Wrap[K comparable, V any](m map[K]V) Map[K, V] {
	return &xmap[K, V]{m: m}
}

func SafeWrap[K comparable, V any](m map[K]V) Map[K, V] {
	return &xmap[K, V]{m: m, mtx: &sync.RWMutex{}}
}

func NewXmap[K comparable, V any]() Map[K, V] {
	return &xmap[K, V]{m: make(map[K]V)}
}

func NewSafeXmap[K comparable, V any]() Map[K, V] {
	return &xmap[K, V]{m: make(map[K]V), mtx: &sync.RWMutex{}}
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

func (m *xmap[K, V]) Load(k K) (v V, loaded bool) {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	v, loaded = m.m[k]
	return
}

func (m *xmap[K, V]) LoadEq(eq func(any) bool) (value V, loaded bool) {
	m.Range(func(key K, val V) bool {
		if eq(key) {
			value, loaded = val, true
			return false
		}
		return true
	})
	return
}

func (m *xmap[K, V]) LoadEqs(eq func(any) bool) (values []V) {
	m.Range(func(key K, value V) bool {
		if eq(key) {
			values = append(values, value)
		}
		return true
	})
	return values
}

func (m *xmap[K, V]) BatchLoad(ks ...K) []Result[V] {
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

func (m *xmap[K, V]) Store(k K, v V) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	m.m[k] = v
}

func (m *xmap[K, V]) BatchStore(ks []K, vs []V) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	must.BeTrue(len(ks) == len(vs))
	for i := range len(ks) {
		m.m[ks[i]] = vs[i]
	}
}

func (m *xmap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
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

func (m *xmap[K, V]) LoadAndDelete(k K) (value V, loaded bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	if value, loaded = m.m[k]; loaded {
		delete(m.m, k)
	}
	return
}

func (m *xmap[K, V]) Delete(k K) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	delete(m.m, k)
}

func (m *xmap[K, V]) BatchDelete(ks ...K) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	for i := range ks {
		delete(m.m, ks[i])
	}
}

func (m *xmap[K, V]) Swap(k K, v V) (previous V, loaded bool) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	previous, loaded = m.m[k]
	m.m[k] = v
	return
}

func (m *xmap[K, V]) CompareAndSwap(k K, old, new V) (swapped bool) {
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

func (m *xmap[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
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

func (m *xmap[K, V]) Range(f func(K, V) bool) {
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

func (m *xmap[K, V]) Clear() {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}
	clear(m.m)
}

func (m *xmap[K, V]) Clone() Map[K, V] {
	if m.mtx != nil {
		m.mtx.RLock()
		defer m.mtx.RUnlock()
	}

	mm := &xmap[K, V]{m: make(map[K]V)}
	if m.mtx != nil {
		mm.mtx = &sync.RWMutex{}
	}

	m.Range(func(k K, v V) bool {
		mm.m[k] = v
		return true
	})

	return mm
}

func NewSmap[K comparable, V any]() Map[K, V] {
	return &smap[K, V]{Map: &sync.Map{}}
}

type smap[K comparable, V any] struct {
	*sync.Map
}

func (s *smap[K, V]) Exists(k K) bool {
	_, exists := s.Load(k)
	return exists
}

func (s *smap[K, V]) Load(k K) (value V, loaded bool) {
	v, exists := s.Map.Load(k)
	if !exists {
		return *new(V), false
	}
	return v.(V), true
}

func (s *smap[K, V]) LoadAndDelete(k K) (value V, loaded bool) {
	v, ok := s.Map.LoadAndDelete(k)
	if !ok {
		return *new(V), ok
	}
	return v.(V), ok
}

func (s *smap[K, V]) LoadEq(eq func(any) bool) (value V, loaded bool) {
	s.Map.Range(func(key, val any) bool {
		if eq(key.(K)) {
			value, loaded = val.(V), true
			return false
		}
		return true
	})
	return
}

func (s *smap[K, V]) LoadEqs(eq func(any) bool) (values []V) {
	s.Map.Range(func(key, value any) bool {
		if eq(key) {
			values = append(values, value.(V))
		}
		return true
	})
	return values
}

func (s *smap[K, V]) BatchLoad(ks ...K) []Result[V] {
	results := make([]Result[V], 0, len(ks))
	for i := range ks {
		vi, loaded := s.Load(ks[i])
		results = append(results, Result[V]{vi, loaded})
	}
	return results
}

func (s *smap[K, V]) Store(k K, v V) {
	s.Map.Store(k, v)
}

func (s *smap[K, V]) BatchStore(ks []K, vs []V) {
	must.BeTrue(len(ks) == len(vs))
	for i := range ks {
		s.Store(ks[i], vs[i])
	}
}

func (s *smap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	val, ok := s.Map.LoadOrStore(k, v)
	return val.(V), ok
}

func (s *smap[K, V]) Delete(k K) {
	s.Map.Delete(k)
}

func (s *smap[K, V]) BatchDelete(ks ...K) {
	for i := range ks {
		s.Delete(ks[i])
	}
}

func (s *smap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return s.Map.CompareAndDelete(key, old)
}

func (s *smap[K, V]) Swap(k K, v V) (previous V, loaded bool) {
	val, ok := s.Map.Swap(k, v)
	return val.(V), ok
}

func (s *smap[K, V]) CompareAndSwap(key K, old V, new V) (swapped bool) {
	return s.Map.CompareAndSwap(key, old, new)
}

func (s *smap[K, V]) Range(f func(K, V) bool) {
	ff := func(k, v any) bool { return f(k.(K), v.(V)) }
	s.Map.Range(ff)
}

func (s *smap[K, V]) Clone() Map[K, V] {
	mm := &sync.Map{}
	s.Range(func(k K, v V) bool {
		mm.Store(k, v)
		return true
	})
	return &smap[K, V]{mm}
}

type Result[V any] struct {
	V      V
	Loaded bool
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

func Equal[K comparable, V any](x, y Map[K, V]) bool {
	if x == nil && y == nil || x == y {
		return true
	}

	if u, ok := x.(*xmap[K, V]); ok {
		if u.mtx != nil {
			u.mtx.RLock()
			defer u.mtx.RUnlock()
		}
	}

	equal := true
	length1 := 0
	length2 := 0
	y.Range(func(k K, v V) bool {
		vv, exist := x.Load(k)
		if !exist || !reflect.DeepEqual(vv, v) {
			equal = false
			return false
		}
		length1++
		return true
	})
	if !equal {
		return false
	}

	x.Range(func(k K, v V) bool {
		length2++
		return true
	})

	return length1 == length2
}

func NewSet[K comparable]() Set[K] {
	return Set[K]{m: NewXmap[K, struct{}]()}
}

func NewSafeSet[K comparable]() Set[K] {
	return Set[K]{m: NewSafeXmap[K, struct{}]()}
}

type Set[K comparable] struct {
	m Map[K, struct{}]
}

func (s *Set[K]) Exists(k K) bool {
	return s.m.Exists(k)
}

func (s *Set[K]) Store(keys ...K) {
	s.m.BatchStore(keys, make([]struct{}, len(keys)))
}

func (s *Set[K]) Delete(keys ...K) {
	s.m.BatchDelete(keys...)
}

func (s *Set[K]) Clear() {
	s.m.Clear()
}

func (s *Set[K]) Keys() []K {
	return Keys(s.m)
}

func (s *Set[K]) Len() int {
	return Len(s.m)
}

func (s *Set[K]) Range(f func(k K) bool) {
	s.m.Range(func(k K, _ struct{}) (shouldContinue bool) {
		return f(k)
	})
}

func (s *Set[K]) Clone() Set[K] {
	return Set[K]{m: s.m.Clone()}
}

func (s *Set[K]) Equal(x Set[K]) bool {
	return Equal(s.m, x.m)
}
