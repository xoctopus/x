package testx_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/xoctopus/x/testx"
)

func TestExpect(t *testing.T) {
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
		Expect(t, "", BeEmpty[string]())
		Expect(t, (*int)(nil), BeEmpty[*int]())
		Expect(t, errors.New("any"), NotBe[error](nil))
		Expect(t, errors.New("any"), Not(BeNil[error]()))
		Expect(t, errors.New("any"), NotBeNil[error]())
		Expect(t, "123", HaveLen[string](3))
		Expect(t, []int{1, 2, 3},
			HaveLen[[]int](3),
			HaveCap[[]int](3),
			Contains(2),
		)
		Expect(t, "123",
			ContainsSubString("1"),
			HavePrefix("1"),
			HaveSuffix("3"),
		)

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
		Expect(t, nil, Not(BeType[any]()))
		Expect(t, nil, Not(BeAssignableTo[any]()))
		Expect(t, nil, Not(BeConvertibleTo[any]()))

		type String string
		Expect[any](t, String("1"), Not(BeType[string]()))
		Expect[any](t, String("1"), Not(BeAssignableTo[string]()))
		Expect[any](t, String("1"), BeConvertibleTo[string]())

		Expect(t, errors.New("any"), Failed())

		Expect(t, []int{1, 2}, EquivalentSlice([]int{2, 1}))
	})

	// this help to verify failure point
	// TODO trace2 case's failure point is not correct
	// t.Run("Trace", func(t *testing.T) {
	// 	// Should be traced to current line as case failure point
	// 	ExpectPanic[error](t, func() { panic("any") })
	// })
	// t.Run("Trace2", func(t *testing.T) {
	// 	ExpectPanic[string](t, func() { crash(1) })
	// })
}

func crash(i int) {
	switch i {
	case 0:
		return
	case 1:
		panic(fmt.Errorf("any"))
	}
}
