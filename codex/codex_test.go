package codex_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/xoctopus/x/codex"
	. "github.com/xoctopus/x/testx"
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

type ECode2 int8

func (e ECode2) Message() string { return "" }

type ECode3 int8

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

	e3 := New(ECode3(2))
	fmt.Printf("expect ecode without Message: %s\n", e3.Error())

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
	// expect ecode without Message: codex_test.ECode3[2]
}

func TestIs(t *testing.T) {
	Expect(t, Is[ECode](New(ECODE__REASON1)), BeTrue())
	Expect(t, Is[ECode2](New(ECODE__REASON2)), BeFalse())
	Expect(t, Is[ECode](errors.New("any")), BeFalse())

	Expect(t, IsCode(New(ECODE__REASON1), ECODE__REASON1), BeTrue())
	Expect(t, IsCode(New(ECODE__REASON1), ECODE__REASON2), BeFalse())

	as1, asserted := As[ECode](New(ECODE__REASON1))
	Expect(t, asserted, BeTrue())
	Expect(t, as1.Code(), Equal(ECODE__REASON1))

	c, as := AsCode[ECode](New(ECODE__REASON1))
	Expect(t, as, BeTrue())
	Expect(t, c, Equal(ECODE__REASON1))

	as2, asserted := As[ECode2](New(ECODE__REASON2))
	Expect(t, as2 == nil, BeTrue())
	Expect(t, asserted, BeFalse())

	_, ok2 := AsCode[ECode2](New(ECODE__REASON2))
	Expect(t, ok2, BeFalse())
}
