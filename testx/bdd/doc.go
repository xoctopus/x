// Package bdd provides a Behavior-Driven Development (BDD) style testing framework.
//
// It offers a domain-specific language (DSL) that is close to natural language,
// making test structures highly readable and expressive.
// The framework is built on Go's standard `testing` package and integrates
// seamlessly with `testing.T`.
//
// # Key Concepts
//
// The BDD framework is structured around three main keywords:
//
//   - Given: Sets up the initial state or preconditions for the test.
//   - When: Defines the action, event, or behavior being tested.
//   - Then: Asserts the expected outcomes using various Checkers.
//
// # Usage Example
//
// You can start a BDD test by wrapping a standard `*testing.T` using `bdd.From(t)`.
//
//	func TestCalculator(t *testing.T) {
//		bdd.From(t).
//			Given("initial value is 1", func(t bdd.T) {
//				v := 1
//
//				t.When("we add 1", func(t bdd.T) {
//					v += 1
//					t.Then("it should equal 2", bdd.Equal(2, v))
//				})
//
//				t.When("we add 2", func(t bdd.T) {
//					v += 2
//					t.Then("it should equal 3", bdd.Equal(3, v))
//					t.Then("it should not equal 4", bdd.NotEqual(4, v))
//				})
//			})
//	}
//
// Alternatively, you can use the `bdd.Given` helper to create a test function directly:
//
//	func TestStringConversion(t *testing.T) {
//		v := 0
//		bdd.Given(func(t bdd.T) {
//			v = 1
//			t.Then(
//				"the string representation should be '1'",
//				bdd.Equal("1", strconv.Itoa(v)),
//			)
//		})(t)
//	}
//
// # Checkers
//
// The `Then` step uses `Checker` interfaces to perform assertions. This package provides a wide range of built-in checkers (e.g., `Equal`, `BeNil`, `BeTrue`, `HaveLen`, `IsError`, etc.) which are wrappers around the matchers from the `github.com/xoctopus/x/testx` package.
package bdd
