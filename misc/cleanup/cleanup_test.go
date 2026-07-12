package cleanup_test

import (
	"errors"
	"fmt"

	"github.com/xoctopus/x/misc/cleanup"
)

func ExampleCollector() {

	func() {
		c := cleanup.NewCollector()

		err := fmt.Errorf("main error")

		defer func() {
			fmt.Println(c.JoinTo(&err))
		}()

		c.Collect(func() error { return errors.New("a") })

		c.Collect(func() error { return errors.New("b") })

		c.Collect(func() error { return errors.New("c") })

		c.Collect((&closer{}).Close)

		c.Collect((&closer{"d"}).Close)

		_ = c.JoinTo(&err)
	}()

	func() {
		c := cleanup.NewCollector()

		var err error

		defer func() {
			fmt.Println(c.JoinTo(&err))
		}()
	}()

	// Output:
	// main error
	// d
	// c
	// b
	// a
	// <nil>
}

type closer struct{ name string }

func (c *closer) Close() error {
	if len(c.name) > 0 {
		return errors.New(c.name)
	}
	return nil
}
