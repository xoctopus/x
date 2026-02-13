package bdd

import (
	"testing"
)

type TB = testing.TB

// T defines a Behavior Driver Development testing.
// It provides a domain-specific language (DSL) that is close to natural language,
// making test structures read like:
//
//	bdd.From(t).Given("Precondition", func(b bdd.T) {
//		// setup precondition
//		b.When("SomeAction", func(b bdd.T) {
//			// do action
//			b.Then("Expects", bdd.Equal("1", v))
//		})
//	})
type T interface {
	TB

	Given(summary string, do func(t T))
	When(summary string, do func(t T))
	Then(summary string, checkers ...Checker)

	Unwrap() *testing.T
}

func From(t *testing.T) T {
	return &bddT{T: t}
}

func Given(setup func(T)) func(t *testing.T) {
	return func(t *testing.T) {
		setup(From(t))
	}
}

type bddT struct {
	*testing.T
}

func (t *bddT) Unwrap() *testing.T {
	return t.T
}

func (t *bddT) Given(summary string, setup func(b T)) {
	if t.Skipped() {
		return
	}

	t.T.Run("GIVEN  "+summary, func(t *testing.T) {
		setup(From(t))
	})
}

func (t *bddT) When(summary string, setup func(b T)) {
	if t.Skipped() {
		return
	}

	t.T.Run("WHEN  "+summary, func(t *testing.T) {
		setup(From(t))
	})
}

func (t *bddT) Then(summary string, checkers ...Checker) {
	if t.Skipped() {
		return
	}

	t.T.Helper()

	t.T.Run("THEN  "+summary, func(t *testing.T) {
		if t.Skipped() {
			return
		}

		t.Helper()

		tt := From(t)

		checked := 0
		for _, c := range checkers {
			if c != nil {
				c.Check(tt)
				checked++
			}
		}
		if checked == 0 {
			t.Logf("case %s has no checkers", t.Name())
		}
	})
}
