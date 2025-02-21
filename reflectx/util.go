package reflectx

import (
	"github.com/pkg/errors"
)

func mustBeTrue(b bool, msg string, args ...any) {
	if !b {
		panic(errors.Errorf(msg, args...))
	}
}
