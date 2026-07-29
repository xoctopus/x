// Package testx provides a comprehensive set of generic matchers and assertion
// utilities for Go testing.
//
// It is designed to be highly extensible, type-safe, and easy to use alongside
// the standard `testing` package.
//
// # Core Assertions
//
// The package provides two primary assertion functions:
//
//   - `Expect`: Asserts that an actual value satisfies one or more matchers.
//   - `ExpectPanic`: Asserts that a function panics, and optionally verifies
//     the panic value against matchers.
//
// Example:
//
//	func TestMath(t *testing.T) {
//		result := 1 + 1
//		testx.Expect(t, result, testx.Equal(2), testx.BeGt(0))
//
//		testx.ExpectPanic(t, func() {
//			panic("something went wrong")
//		}, testx.Equal("something went wrong"))
//	}
//
// # Matchers
//
// Matchers are the building blocks of assertions. A `Matcher[T]` evaluates an actual
// value of type `T` and returns whether it meets the expected criteria.
// The package includes a wide variety of built-in matchers:
//
//   - Equality & Comparison: `Equal`, `Be`, `BeGt`, `BeLt`, `IsZero`, etc.
//   - Collections: `HaveLen`, `HaveCap`, `HaveKey`, `Contains`, `EquivalentSlice`, etc.
//   - Strings: `HavePrefix`, `HaveSuffix`, `ContainsSubString`, `MatchRegexp`.
//   - Types & Interfaces: `BeAssignableTo`, `BeConvertibleTo`, `IsType`.
//   - Errors: `IsError`, `AsError`, `ErrorEqual`, `Succeed`, `Failed`, etc.
//
// # Custom Matchers
//
// You can easily create custom matchers using `NewMatcher` or `NewComparedMatcher`:
//
//	func BeEven() testx.Matcher[int] {
//		return testx.NewMatcher("BeEven", func(actual int) bool {
//			return actual%2 == 0
//		})
//	}
//
// # BDD Integration
//
// For a more expressive, Behavior-Driven Development style, consider using the
// `github.com/xoctopus/x/testx/bdd` sub-package, which builds upon the matchers
// defined here.
package testx
