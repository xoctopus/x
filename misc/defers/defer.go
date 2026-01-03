package defers

import (
	"errors"

	"github.com/xoctopus/x/misc/must"
)

func Collect(f func() error, dst *error) {
	if err := f(); err != nil {
		must.NotNilF(dst, "destinational error must be assigned")
		if *dst == nil {
			*dst = err
			return
		}
		*dst = errors.Join(*dst, err)
	}
}
