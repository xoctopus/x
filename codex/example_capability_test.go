package codex_test

import (
	"fmt"

	"github.com/xoctopus/x/codex"
)

// OrderError and StatusDescriber are defined in the package example.

func ExampleValuer() {
	err := codex.New(ORDER_ERROR__CLOSED)

	if v, ok := err.(codex.Valuer); ok {
		code := v.CodeValue()
		fmt.Printf("%T(%d)\n", code, code)
		if d, ok := code.(StatusDescriber); ok {
			fmt.Println(d.StatusCode())
		}
	}

	// Output:
	// codex_test.OrderError(2)
	// 409
}

func ExampleCanBeCode() {
	var c codex.CanBeCode = ORDER_ERROR_UNDEFINED

	for _, v := range c.CodeValues() {
		fmt.Printf("%d %s\n", v, v.(codex.ErrorMessage).Message())
	}

	// Output:
	// 1 order not found
	// 2 order closed
}
