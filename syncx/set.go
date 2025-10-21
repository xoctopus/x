package syncx

func NewSet[K comparable]() Set[K] {
	return Set[K]{m: NewXmap[K, struct{}]()}
}

type Set[K comparable] struct {
	m Map[K, struct{}]
}

func (s *Set[K]) Exists(k K) bool {
	for x := range s.Range {
		if x == k {
			return true
		}
	}
	return false
}

func (s *Set[K]) Store(k K) {
	s.m.Store(k, struct{}{})
}

func (s *Set[K]) Delete(k K) {
	s.m.Delete(k)
}

func (s *Set[K]) Clear() {
	s.m.Clear()
}

func (s *Set[K]) Len() int {
	return s.m.Len()
}

func (s *Set[K]) Keys() []K {
	return s.m.Keys()
}

func (s *Set[K]) Range(f func(k K) bool) {
	s.m.Range(func(k K, _ struct{}) (shouldContinue bool) {
		return f(k)
	})
}
