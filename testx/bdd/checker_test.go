package bdd_test

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/bdd"
)

type CodeError int8

func (c CodeError) Message() string {
	return strconv.Itoa(int(c))
}

func TestCheckers(t *testing.T) {
	err1 := errors.New("any")
	_, err2 := os.Open("non_existent_file.txt")

	tt := bdd.From(t)

	for _, c := range []bdd.Checker{
		bdd.BeNil[any](nil),
		bdd.NotBeNil[any](1),
		bdd.BeTrue(true),
		bdd.BeFalse(false),
		bdd.IsZero(0),
		bdd.IsNotZero(1),
		bdd.IsZero(struct{}{}),
		bdd.IsZero(map[string]any{}),
		bdd.Be(1, 1),
		bdd.NotBe(1, 2),
		bdd.Equal(1, 1),
		bdd.NotEqual(1, 2),
		bdd.BeGt(2, 1),
		bdd.BeGte(1, 1),
		bdd.BeGte(2, 1),
		bdd.BeLt(1, 2),
		bdd.BeLte(1, 1),
		bdd.BeLte(1, 2),
		bdd.HaveCap([]int{1, 2, 3}, 3),
		bdd.HaveLen([]int{1, 2, 3}, 3),
		bdd.HaveKey(map[string]int{"a": 1}, "a"),
		bdd.HavePrefix("abc", "a"),
		bdd.HaveSuffix("abc", "c"),
		bdd.ContainsSubString("abc", "b"),
		bdd.MatchRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "test@example.com"),
		bdd.Contains([]int{1, 2, 3}, 1),
		bdd.EquivalentSlice([]int{1, 2, 3}, []int{3, 2, 1}),
		bdd.ConsistOfSlice([]int{1, 2, 3}, []int{3, 2, 1}),
		bdd.BeAssignableTo[any](1),
		bdd.BeAssignableTo[int](any(1)),
		bdd.BeConvertibleTo[any](1),
		bdd.BeConvertibleTo[int](any(1)),
		bdd.IsType[int](any(1)),
		bdd.IsError(err1, err1),
		bdd.AsError(&err1, err1),
		bdd.AsErrorType[*os.PathError](err2),
		bdd.IsCodeError(codex.New(CodeError(1)), CodeError(1)),
		bdd.ErrorEqual(err1, "any"),
		bdd.AsNegativeChecker(testx.ErrorEqual("other"), err1),
		bdd.ErrorContains(err1, "an"),
		bdd.Succeed(nil),
		bdd.Failed(err1),
	} {
		c.Check(t)
		c.Check(tt)
	}
}
