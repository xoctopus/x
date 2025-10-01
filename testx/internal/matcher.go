package internal

type Matcher[A any] interface {
	Action() string
	Match(A) bool
	Negative() bool
}

type NormalizedMatcher[T any] interface {
	NormalizeActual(T) any
}

type NormalizedExpectedMatcher interface {
	NormalizeExpect() any
}

func NewMatcher[A any](action string, match func(A) bool) Matcher[A] {
	return &matcher[A]{
		action: action,
		match:  match,
	}
}

type matcher[A any] struct {
	action string
	match  func(A) bool
}

func (m *matcher[A]) Action() string { return m.action }

func (m *matcher[A]) Match(actual A) bool { return m.match(actual) }

func (m *matcher[A]) Negative() bool { return false }

func Not[T any](matcher Matcher[T]) Matcher[T] {
	return &negative[T]{Matcher: matcher}
}

type negative[A any] struct {
	Matcher[A]
}

func (m *negative[A]) Negative() bool { return true }

type MatcherNewer[A any, E any] func(e E) Matcher[A]

func NewComparedMatcher[A any, E any](action string, match func(A, E) bool) MatcherNewer[A, E] {
	return func(expected E) Matcher[A] {
		return &compared[A, E]{
			action:   action,
			match:    match,
			expected: expected,
		}
	}
}

type compared[A any, E any] struct {
	action   string
	match    func(A, E) bool
	expected E
}

func (m *compared[A, E]) Action() string { return m.action }

func (m *compared[A, E]) Match(actual A) bool {
	return m.match(actual, m.expected)
}

func (m *compared[A, E]) Negative() bool { return false }

func (m *compared[A, E]) NormalizeExpect() any { return m.expected }
