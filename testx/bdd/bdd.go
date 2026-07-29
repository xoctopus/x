package bdd

import (
	"testing"
)

// TB is an alias for testing.TB to simplify the interface signature.
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

	// Given defines a precondition or initial state for the test.
	Given(preconditionSummary string, do func(t T))
	// When defines an action or event that triggers the behavior being tested.
	When(actionSummary string, do func(t T))
	// Then defines the expected outcome or assertions to be verified.
	Then(outcomeSummary string, checkers ...Checker)

	// Unwrap returns the underlying *testing.T instance.
	Unwrap() *testing.T
}

// From wraps a standard *testing.T into a BDD-style T interface.
func From(t *testing.T) T {
	return &bddT{T: t}
}

// Given is a helper function that creates a standard testing function
// which initializes a BDD context and runs the provided setup function.
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
	if !t.Skipped() {
		t.T.Run("GIVEN  "+summary, func(t *testing.T) {
			setup(From(t))
		})
	}
}

func (t *bddT) When(summary string, setup func(b T)) {
	t.Helper()

	if !t.Skipped() {
		t.T.Run("WHEN  "+summary, func(t *testing.T) {
			setup(From(t))
		})
	}
}

func (t *bddT) Then(summary string, checkers ...Checker) {
	t.Helper()

	if !t.Skipped() {
		t.T.Run("THEN  "+summary, func(t *testing.T) {
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
}
