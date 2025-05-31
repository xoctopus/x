package textx

import (
	"fmt"
	"reflect"
)

func typ(v any) string {
	switch x := v.(type) {
	case reflect.Value:
		return x.Type().String()
	default:
		return fmt.Sprintf("%T", x)
	}
}

func NewErrMarshalUnsupportedType(got any) error {
	return &ErrMarshalUnsupportedType{got: typ(got)}
}

type ErrMarshalUnsupportedType struct {
	got string
}

func (e *ErrMarshalUnsupportedType) Error() string {
	return "marshal text unsupported type, got `" + e.got + "`"
}

func NewErrMarshalFailed(input any, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrMarshalFailed{input: typ(input), cause: cause}
}

type ErrMarshalFailed struct {
	input string
	cause error
}

func (e *ErrMarshalFailed) Error() string {
	return fmt.Sprintf("failed to call (%s).MarshalText: [%+v]", e.input, e.cause)
}

func (e *ErrMarshalFailed) Unwrap() error {
	return e.cause
}

func NewErrUnmarshalUnsupportedType(got any) error {
	return &ErrUnmarshalUnsupportedType{typ(got)}
}

type ErrUnmarshalUnsupportedType struct {
	got string
}

func (e *ErrUnmarshalUnsupportedType) Error() string {
	return "unmarshal text unsupported type, got `" + e.got + "`"
}

func NewErrUnmarshalInvalidInput(got any) error {
	return &ErrUnmarshalInvalidInput{got: typ(got)}
}

type ErrUnmarshalInvalidInput struct {
	got string
}

func (e *ErrUnmarshalInvalidInput) Error() string {
	return "unmarshal text invalid input, got `" + e.got + "`, it MUST be valid and can set"
}

func NewErrUnmarshalParseFailed(from []byte, to any, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrUnmarshalParseFailed{from, typ(to), cause}
}

type ErrUnmarshalParseFailed struct {
	from  []byte
	to    string
	cause error
}

func (e *ErrUnmarshalParseFailed) Error() string {
	return fmt.Sprintf("failed to parse `%s` from `%s`: [%+v]", e.to, e.from, e.cause)
}

func (e *ErrUnmarshalParseFailed) Unwrap() error {
	return e.cause
}

func NewErrUnmarshalFailed(from []byte, to any, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrUnmarshalFailed{from, typ(to), cause}
}

type ErrUnmarshalFailed struct {
	from  []byte
	to    string
	cause error
}

func (e *ErrUnmarshalFailed) Error() string {
	return fmt.Sprintf("failed to call (%s).UnmarshalText(%s): [%+v]", e.to, e.from, e.cause)
}

func (e *ErrUnmarshalFailed) Unwrap() error {
	return e.cause
}

func NewErrMarshalURLInvalidInput(got any) error {
	return &ErrMarshalURLInvalidInput{got: typ(got)}
}

type ErrMarshalURLInvalidInput struct {
	got string
}

func (e *ErrMarshalURLInvalidInput) Error() string {
	return "marshal url invalid input, got `" + e.got + "`, its underlying MUST be a struct"
}

func NewErrMarshalURLFailed(input any, field string, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrMarshalURLFailed{
		input: typ(input),
		field: field,
		cause: cause,
	}
}

type ErrMarshalURLFailed struct {
	input string
	field string
	cause error
}

func (e *ErrMarshalURLFailed) Error() string {
	return fmt.Sprintf("failed to call MarshalText at field `%s.%s`: [%+v]", e.input, e.field, e.cause)
}

func (e *ErrMarshalURLFailed) Unwrap() error {
	return e.cause
}

func NewErrUnmarshalURLInvalidInput(got any) error {
	return &ErrUnmarshalURLInvalidInput{got: typ(got)}
}

type ErrUnmarshalURLInvalidInput struct {
	got string
}

func (e *ErrUnmarshalURLInvalidInput) Error() string {
	return "unmarshal url invalid input, got `" + e.got + "`, its underlying MUST be a struct and can set"
}

func NewErrUnmarshalURLFailed(tpe any, field, input string, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrUnmarshalURLFailed{
		typename: typ(tpe),
		field:    field,
		input:    input,
		cause:    cause,
	}
}

type ErrUnmarshalURLFailed struct {
	typename string
	field    string
	input    string
	cause    error
}

func (e *ErrUnmarshalURLFailed) Error() string {
	return fmt.Sprintf("unmarshal`%s.%s` failed from '%s': [%+v]", e.typename, e.field, e.input, e.cause)
}

func (e *ErrUnmarshalURLFailed) Unwrap() error {
	return e.cause
}
