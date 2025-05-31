package textx_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/textx"
)

type ErrorWrapper interface {
	Unwrap() error
}

func TestNewErrs(t *testing.T) {
	cause := errors.New("cause")

	v := new(int)
	r := reflect.ValueOf(v)
	text := []byte("invalid")

	for _, e := range []error{
		NewErrMarshalUnsupportedType(nil),
		NewErrMarshalUnsupportedType(v),
		NewErrMarshalUnsupportedType(r),
		NewErrMarshalFailed(nil, cause),
		NewErrMarshalFailed(v, cause),
		NewErrMarshalFailed(r, cause),
		NewErrUnmarshalUnsupportedType(nil),
		NewErrUnmarshalUnsupportedType(v),
		NewErrUnmarshalUnsupportedType(r),
		NewErrUnmarshalInvalidInput(nil),
		NewErrUnmarshalInvalidInput(v),
		NewErrUnmarshalInvalidInput(r),
		NewErrUnmarshalParseFailed(text, nil, cause),
		NewErrUnmarshalParseFailed(text, v, cause),
		NewErrUnmarshalParseFailed(text, r, cause),
		NewErrUnmarshalFailed(text, nil, cause),
		NewErrUnmarshalFailed(text, v, cause),
		NewErrUnmarshalFailed(text, r, cause),
		NewErrMarshalURLInvalidInput(nil),
		NewErrMarshalURLInvalidInput(v),
		NewErrMarshalURLInvalidInput(r),
		NewErrMarshalURLFailed(v, "field", cause),
		NewErrMarshalURLFailed(r, "field", cause),
		NewErrUnmarshalURLInvalidInput(nil),
		NewErrUnmarshalURLInvalidInput(v),
		NewErrUnmarshalURLInvalidInput(r),
		NewErrUnmarshalURLFailed(v, "field", "value", cause),
		NewErrUnmarshalURLFailed(r, "field", "value", cause),
	} {
		NewWithT(t).Expect(e).NotTo(BeNil())
		t.Log(e)
		if u, ok := e.(ErrorWrapper); ok {
			NewWithT(t).Expect(u.Unwrap()).To(Equal(cause))
		}
	}

	for _, e := range []error{
		NewErrMarshalFailed(nil, nil),
		NewErrMarshalURLFailed(nil, "field", nil),
		NewErrUnmarshalParseFailed(text, nil, nil),
		NewErrUnmarshalFailed(text, nil, nil),
		NewErrUnmarshalURLFailed(nil, "field", "value", nil),
	} {
		NewWithT(t).Expect(e).To(BeNil())
	}
}
