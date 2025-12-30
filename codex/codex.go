package codex

import (
	"errors"
	"fmt"
)

// Code is a generic interface that represents an error code with an underlying
// integer type. It uses Go's type approximation (~int, ~int8, etc.) to allow
// any integer-based type. The Message method returns the error description
// associated with the code.
type Code interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	// Message returns the error description string for the code.
	Message() string
}

// Error is a generic interface defining an error with a typed code.
// It extends the standard error interface with methods to access the typed code,
// unwrap the underlying cause, and compare errors by code. The generic parameter
// C must satisfy the Code interface.
type Error[C Code] interface {
	// Code returns typed code value
	Code() C
	// Error implements an error
	Error() string
	// Unwrap unwraps caused error
	Unwrap() error
	// Is returns if input is same type of CodeError and with same code
	Is(error) bool
}

func New[C Code](e C) error {
	return &coderr[C]{code: e}
}

func Errorf[C Code](e C, format string, args ...any) error {
	return &coderr[C]{code: e, msg: format, args: args}
}

func Wrap[C Code](e C, cause error) error {
	if cause == nil {
		return nil
	}
	return &coderr[C]{code: e, cause: cause, args: []any{cause}}
}

func Wrapf[C Code](e C, cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}
	return &coderr[C]{code: e, cause: cause, msg: format, args: append(args, cause)}
}

func IsCode[C Code](e error, code C) bool {
	return errors.Is(e, New(code))
}

func Is[C Code](e error) (C, bool) {
	var target *coderr[C]
	if errors.As(e, &target) {
		return target.code, true
	}
	return *new(C), false
}

func As[C Code](e error) (Error[C], bool) {
	var target *coderr[C]
	if errors.As(e, &target) {
		return target, true
	}
	return nil, false
}

type coderr[C Code] struct {
	code  C
	msg   string
	args  []any
	cause error
}

func (e *coderr[C]) Error() string {
	msg := e.code.Message()
	if len(e.msg) > 0 {
		msg += ". " + e.msg
	}
	if e.cause != nil {
		msg += ". [cause: %+v]"
	}
	return fmt.Sprintf(msg, e.args...)
}

func (e *coderr[C]) Code() C {
	return e.code
}

func (e *coderr[C]) Unwrap() error {
	return e.cause
}

func (e *coderr[C]) Is(err error) bool {
	var target *coderr[C]
	return errors.As(err, &target) && target.Code() == e.Code()
}
