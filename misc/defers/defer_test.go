package defers_test

import (
	"errors"
	"fmt"

	"github.com/xoctopus/x/misc/defers"
)

func ExampleCollect() {
	var err error

	func() {
		fmt.Println("Single:")
		defer func() { fmt.Println(err) }()

		defer defers.Collect(func() error { return errors.New("a") }, &err)
	}()

	err = nil
	func() {
		fmt.Println("Multi:")
		defer func() { fmt.Println(err) }()

		defer defers.Collect(func() error { return errors.New("d") }, &err)
		defer defers.Collect(func() error { return errors.New("c") }, &err)
		defer defers.Collect((&closer{"b"}).Close, &err)
		defer defers.Collect((&closer{"a"}).Close, &err)
	}()

	// Output:
	// Single:
	// a
	// Multi:
	// a
	// b
	// c
	// d
}

type closer struct{ name string }

func (c *closer) Close() error { return errors.New(c.name) }
