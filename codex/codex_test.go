package codex_test

import (
	"errors"
	"fmt"

	. "github.com/xoctopus/x/codex"
)

type ECode int8

const (
	ECODE_UNDEFINED ECode = iota + 1
	ECODE__REASON1        // reason1
	ECODE__REASON2        // reason2
)

func (e ECode) Message() string {
	prefix := fmt.Sprintf("[region:%d] ", e)

	switch e {
	case ECODE_UNDEFINED:
		return prefix + "undefined"
	case ECODE__REASON1:
		return prefix + "reason1"
	case ECODE__REASON2:
		return prefix + "reason2"
	default:
		return prefix + "unknown"
	}
}

func ExampleError() {
	fmt.Println(New(ECODE_UNDEFINED).Error())
	e0 := Errorf(ECODE__REASON1, "user message: %d", 1)
	fmt.Println(e0.Error())
	e1 := Wrap(ECODE__REASON1, errors.New("cause1"))
	fmt.Println(e1.Error())
	fmt.Printf("cause by unwrapping: %v\n", errors.Unwrap(e1))
	e2 := Wrapf(ECODE__REASON2, errors.New("cause2"), "user message: %s", "any")
	fmt.Println(e2.Error())
	fmt.Printf("cause by unwrapping: %v\n", errors.Unwrap(e2))

	fmt.Printf("expecting true errors.Is(e0, e1): %v\n", errors.Is(e0, e1))
	fmt.Printf("expecting false errors.Is(e1, e2): %v\n", errors.Is(e1, e2))
	fmt.Printf("expecting nil Wrap(ECODE__REASON1, nil): %v\n", Wrap(ECODE__REASON1, nil))
	fmt.Printf("expecting nil Wrapf(ECODE__REASON2, nil): %v\n", Wrapf(ECODE__REASON2, nil, "whatever"))

	// Output:
	// [region:1] undefined
	// [region:2] reason1. user message: 1
	// [region:2] reason1. [cause: cause1]
	// cause by unwrapping: cause1
	// [region:3] reason2. user message: any. [cause: cause2]
	// cause by unwrapping: cause2
	// expecting true errors.Is(e0, e1): true
	// expecting false errors.Is(e1, e2): false
	// expecting nil Wrap(ECODE__REASON1, nil): <nil>
	// expecting nil Wrapf(ECODE__REASON2, nil): <nil>
}
