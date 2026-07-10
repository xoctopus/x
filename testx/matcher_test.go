package testx_test

import (
	"errors"
	"testing"

	"github.com/xoctopus/x/codex"
	. "github.com/xoctopus/x/testx"
)

func TestMatchers(t *testing.T) {
	t.Run("Matchers", func(t *testing.T) {
		Expect(t, "1", Equal("1"))
		Expect(t, map[string]int{
			"1": 10,
			"2": 11,
			"9": 19,
		}, Equal(map[string]int{
			"1": 10,
			"2": 11,
			"9": 19,
		}))
		Expect(t, "1", NotEqual("2"))
		Expect(t, error(nil), BeNil[error]())
		Expect(t, true, BeTrue())
		Expect(t, false, BeFalse())
		Expect(t, "", IsZero[string]())
		Expect(t, "a", IsNotZero[string]())
		Expect(t, (*int)(nil), IsZero[*int]())
		Expect(t, errors.New("any"), NotBe[error](nil))
		Expect(t, errors.New("any"), Not(BeNil[error]()))
		Expect(t, errors.New("any"), NotBeNil[error]())
		Expect(t, "123", HaveLen[string](3))
		Expect(t, []int{1, 2, 3},
			HaveLen[[]int](3),
			HaveCap[[]int](3),
			Contains[int, []int](2),
		)
		Expect(t, "123",
			ContainsSubString("1"),
			HavePrefix("1"),
			HaveSuffix("3"),
		)

		Expect(t, 1, Not(HaveLen[int](0)))
		Expect(t, 1, Not(HaveCap[int](0)))
		Expect(t, map[string]int{"a": 0}, HaveKey[string, int, map[string]int]("a"))

		Expect(t, "test@example.com", MatchRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))

		Expect(t, 2, BeGt(1))
		Expect(t, 1, BeGte(1))
		Expect(t, 1, BeLt(2))
		Expect(t, 1, BeLte(1))

		ExpectPanic[string](t, func() { panic("any") })
		ExpectPanic[string](t, func() { panic("any") },
			Equal("any"),
			HavePrefix("a"),
			HaveSuffix("y"),
		)

		var (
			err     = errors.New("any")
			closure = func() { panic(err) }
		)

		ExpectPanic[error](t, func() { panic(err) })
		ExpectPanic[error](t, closure,
			NotBeNil[error](),
			Not(IsError(errors.New("any"))),
			ErrorEqual("any"),
			IsError(err),
			ErrorContains("a"),
		)

		Expect[error](t, nil, Not(ErrorEqual("any")))
		Expect[error](t, nil, Not(ErrorContains("any")))

		ExpectPanic(t, func() { panic("any") }, Equal("any"))

		// refuse comparing nils
		Expect(t, nil, Not(IsType[any]()))
		Expect(t, nil, Not(BeAssignableTo[any]()))
		Expect(t, nil, Not(BeConvertibleTo[any]()))

		type String string
		Expect[any](t, String("1"), Not(IsType[string]()))
		Expect[any](t, String("1"), Not(BeAssignableTo[string]()))
		Expect[any](t, String("1"), BeConvertibleTo[string]())

		Expect(t, errors.New("any"), Failed())

		Expect(t, []int{1, 2}, EquivalentSlice([]int{2, 1}))
		Expect(t, []int{1, 2}, ConsistOfSlice([]int{2, 1}))

		err = codex.New(MockCodeErr(0))
		Expect(t, err, Not(IsCodeError(MockCodeErr(1))))
		Expect(t, err, IsCodeError(MockCodeErr(0)))
		Expect(t, err, AsError(new(codex.New(MockCodeErr(0)))))
		Expect(t, err, AsErrorType[codex.Error[MockCodeErr]]())
		Expect(t, nil, Not(AsErrorType[codex.Error[MockCodeErr]]()))
	})
}
