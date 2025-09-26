package reflectx

import (
	"github.com/pkg/errors"
)

type ErrorCode int8

const (
	E_INVALID_FLAG_KEY ErrorCode = iota + 1
	E_INVALID_FLAG_VALUE
	E_INVALID_FLAG_NAME
	E_INVALID_OPTION_KEY
	E_INVALID_OPTION_VALUE
	E_INVALID_OPTION_UNQUOTED
)

var errMessages = map[ErrorCode]string{
	E_INVALID_FLAG_KEY:        "invalid flag key",
	E_INVALID_FLAG_VALUE:      "invalid flag value",
	E_INVALID_FLAG_NAME:       "invalid flag name",
	E_INVALID_OPTION_KEY:      "invalid option key",
	E_INVALID_OPTION_VALUE:    "invalid option value",
	E_INVALID_OPTION_UNQUOTED: "invalid option unquoted",
}

func NewError(code ErrorCode, inf string) error {
	return errors.WithStack(&Error{code: code, info: inf})
}

type Error struct {
	code ErrorCode
	info string
}

func (e *Error) Error() string {
	detail := ""
	if len(e.info) > 0 {
		detail = ": `" + e.info + "`"
	}
	return errMessages[e.code] + detail
}

func (e *Error) Is(err error) bool {
	var target *Error
	ok := errors.As(err, &target)
	return ok && e.code == e.code
}
