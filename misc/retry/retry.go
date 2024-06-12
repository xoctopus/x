package retry

import (
	"log"
	"reflect"
	"time"

	"github.com/xoctopus/x/reflectx"
)

type Retry struct {
	Repeats  int
	Interval time.Duration
}

func (r *Retry) SetDefault() {
	reflectx.Set(reflect.ValueOf(r), reflect.ValueOf(Default))
}

func (r Retry) Do(exec func() error) (err error) {
	if r.Repeats <= 0 {
		return exec()
	}
	for i := 0; i < r.Repeats; i++ {
		if err = exec(); err != nil {
			log.Printf("retry in %s [err: %v]", r.Interval, err)
			time.Sleep(r.Interval)
			continue
		}
		break
	}
	return
}

var Default = &Retry{3, 3 * time.Second}

func Do(retry *Retry, exec func() error) error {
	return retry.Do(exec)
}
