package codex_test

import (
	"errors"
	"fmt"

	"github.com/xoctopus/x/codex"
)

type OrderError int8

const (
	ORDER_ERROR_UNDEFINED OrderError = iota
	ORDER_ERROR__NOT_FOUND
	ORDER_ERROR__CLOSED
)

func (e OrderError) Message() string {
	switch e {
	case ORDER_ERROR__NOT_FOUND:
		return "order not found"
	case ORDER_ERROR__CLOSED:
		return "order closed"
	default:
		return "undefined"
	}
}

func (OrderError) CodeValues() []any {
	return []any{ORDER_ERROR__NOT_FOUND, ORDER_ERROR__CLOSED}
}

// StatusCode is a capability defined for consumers, codex is unaware of it.
func (e OrderError) StatusCode() int {
	switch e {
	case ORDER_ERROR__NOT_FOUND:
		return 404
	case ORDER_ERROR__CLOSED:
		return 409
	default:
		return 0
	}
}

// StatusDescriber is defined by the consumer which needs the capability.
type StatusDescriber interface {
	StatusCode() int
}

func Example() {
	cause := errors.New("sql: no rows")
	err := fmt.Errorf("query order: %w", codex.Wrap(ORDER_ERROR__NOT_FOUND, cause))
	fmt.Println(err)

	// the code type is known
	fmt.Println(codex.IsCode(err, ORDER_ERROR__NOT_FOUND))
	code, _ := codex.AsCode[OrderError](err)
	fmt.Println(code.Message())

	// the code type is unknown: view the error as its code
	var d StatusDescriber
	if errors.As(err, &d) {
		fmt.Println(d.StatusCode())
	}

	// Output:
	// query order: order not found. [cause: sql: no rows]
	// true
	// order not found
	// 404
}
