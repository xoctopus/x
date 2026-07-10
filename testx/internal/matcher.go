package internal

type Matcher[Actual any] interface {
	// Action returns matching action name
	Action() string
	// Match defines matching action
	Match(Actual) bool
	// Negative returns if negative action. eg: Equal, NotEqual
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

type matcher[Actual any] struct {
	action string
	match  func(Actual) bool
}

func (m *matcher[Actual]) Action() string {
	return m.action
}

func (m *matcher[Actual]) Match(actual Actual) bool {
	return m.match(actual)
}

func (m *matcher[Actual]) Negative() bool {
	return false
}

func Not[Actual any](matcher Matcher[Actual]) Matcher[Actual] {
	return &negative[Actual]{Matcher: matcher}
}

type negative[Actual any] struct {
	Matcher[Actual]
}

func (m *negative[Actual]) Negative() bool {
	return true
}

type MatcherNewer[Actual any, Expect any] func(expect Expect) Matcher[Actual]

func NewComparedMatcher[Actual any, Expect any](action string, match func(Actual, Expect) bool) MatcherNewer[Actual, Expect] {
	return func(expect Expect) Matcher[Actual] {
		return &compared[Actual, Expect]{
			action: action,
			match:  match,
			expect: expect,
		}
	}
}

type compared[Actual any, Expect any] struct {
	action string
	match  func(Actual, Expect) bool
	expect Expect
}

func (m *compared[Actual, Expect]) Action() string {
	return m.action
}

func (m *compared[Actual, Expect]) Match(actual Actual) bool {
	return m.match(actual, m.expect)
}

func (m *compared[Actual, Expect]) Negative() bool {
	return false
}

func (m *compared[Actual, Expect]) NormalizeActual(actual Actual) any {
	return actual
}

func (m *compared[Actual, Expect]) NormalizeExpect() any {
	return m.expect
}
