package codex

import (
	"errors"
	"fmt"
	"reflect"
)

// Code constrains the underlying integer type of an error code. It uses Go's
// type approximation (~int, ~int8, etc.) to allow integer-based types whose
// values always fit in int64, so ~uint is excluded.
type Code interface {
	~int | ~int8 | ~int16 | ~int32 | ~uint8 | ~uint16 | ~uint32
}

// ErrorMessage is implemented by code types that describe themselves. When the
// carried code implements it, the error message starts with Message instead of
// the default `Type[value]` form.
type ErrorMessage interface {
	// Message returns the error description string for the code.
	Message() string
}

// CanBeCode is implemented by code types that can list all of their values,
// for example to document every code of a domain.
type CanBeCode interface {
	// CodeValues returns all defined values of the code type.
	CodeValues() []any
}

// Valuer exposes the typed code carried by an error to callers which do not know
// the code type, so they can assert capabilities on the code by their own interfaces.
type Valuer interface {
	// CodeValue returns the carried code as its original type.
	CodeValue() any
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

// New returns an error carrying code e.
func New[C Code](e C) error {
	return &coderr[C]{code: e}
}

// Errorf returns an error carrying code e, with a formatted message appended to
// the code message.
func Errorf[C Code](e C, format string, args ...any) error {
	return &coderr[C]{code: e, msg: format, args: args}
}

// Wrap returns an error carrying code e and caused by cause. It returns nil if
// cause is nil.
func Wrap[C Code](e C, cause error) error {
	if cause == nil {
		return nil
	}
	return &coderr[C]{code: e, cause: cause, args: []any{cause}}
}

// Wrapf is like Wrap, with a formatted message appended to the code message.
func Wrapf[C Code](e C, cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}
	return &coderr[C]{code: e, cause: cause, msg: format, args: append(args, cause)}
}

// IsCode reports whether any error in e's chain carries code.
func IsCode[C Code](e error, code C) bool {
	return errors.Is(e, New(code))
}

// Is reports whether any error in e's chain carries a code of type C.
func Is[C Code](e error) bool {
	_, ok := errors.AsType[*coderr[C]](e)
	return ok
}

// As finds the first error in e's chain carrying a code of type C.
func As[C Code](e error) (Error[C], bool) {
	if target, ok := errors.AsType[*coderr[C]](e); ok {
		return target, true
	}
	return nil, false
}

// AsCode is like As, but returns the carried code directly.
func AsCode[C Code](e error) (C, bool) {
	if ce, ok := As[C](e); ok {
		return ce.Code(), true
	}
	return *new(C), false
}

type coderr[C Code] struct {
	code  C
	msg   string
	args  []any
	cause error
}

// Error formats as `code message[. user message][. [cause: cause]]`.
func (e *coderr[C]) Error() string {
	msg := fmt.Sprintf("%T[%d]", e.code, e.code)
	if x, ok := any(e.code).(ErrorMessage); ok {
		msg = x.Message()
	}
	if len(e.msg) > 0 {
		msg += ". " + e.msg
	}
	if e.cause != nil {
		msg += ". [cause: %+v]"
	}
	return fmt.Sprintf(msg, e.args...)
}

// Code returns the carried code.
func (e *coderr[C]) Code() C {
	return e.code
}

// CodeValue implements Valuer.
func (e *coderr[C]) CodeValue() any {
	return e.code
}

// As lets errors.As and errors.AsType view the error as its code: the code is
// assigned to target when the code type is assignable to the target's element type.
func (e *coderr[C]) As(target any) bool {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return false
	}
	if rv = rv.Elem(); !reflect.TypeOf(e.code).AssignableTo(rv.Type()) {
		return false
	}
	rv.Set(reflect.ValueOf(e.code))
	return true
}

// Unwrap returns the cause.
func (e *coderr[C]) Unwrap() error {
	return e.cause
}

// Is reports whether err carries the same code of the same type.
func (e *coderr[C]) Is(err error) bool {
	target, ok := errors.AsType[*coderr[C]](err)
	return ok && target.Code() == e.Code()
}
