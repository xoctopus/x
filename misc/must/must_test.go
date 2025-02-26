package must_test

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

func ReturnError() error {
	return errors.New("some error")
}

func ReturnNoError() error {
	return nil
}

func ReturnIntError() (int, error) {
	return 100, errors.New("some error")
}

func ReturnIntNoError() (int, error) {
	return 100, nil
}

func ReturnTrue() bool {
	return true
}

func ReturnFalse() bool {
	return false
}

func ReturnIntTrue() (int, bool) {
	return 100, true
}

func ReturnIntFalse() (int, bool) {
	return 100, false
}

func ExampleNoError() {
	must.NoError(ReturnNoError())

	defer func() {
		fmt.Println(recover())
	}()
	must.NoError(ReturnError())

	// Output:
	// some error
}

func ExampleNoErrorV() {
	fmt.Println(must.NoErrorV(ReturnIntNoError()))

	defer func() {
		fmt.Println(recover())
	}()
	must.NoErrorV(ReturnIntError())

	// Output:
	// 100
	// some error
}

func ExampleNoErrorWrap() {
	must.NoErrorWrap(ReturnNoError(), "any")

	defer func() {
		fmt.Println(recover())
	}()
	must.NoErrorWrap(ReturnError(), "some message: %d", 10)

	// Output:
	// some message: 10: some error
}

func ExampleBeTrue() {
	must.BeTrue(ReturnTrue())

	defer func() {
		fmt.Println(recover())
	}()
	must.BeTrue(ReturnFalse())

	// Output:
	// must be true
}

func ExampleBeTrueV() {
	fmt.Println(must.BeTrueV(ReturnIntTrue()))

	defer func() {
		fmt.Println(recover())
	}()
	_ = must.BeTrueV(ReturnIntFalse())

	// Output:
	// 100
	// must be true
}

func ExampleBeTrueWrap() {
	must.BeTrueWrap(ReturnTrue(), "any")

	defer func() {
		fmt.Println(recover())
	}()
	must.BeTrueWrap(ReturnFalse(), "required exists")

	// Output:
	// must be true: required exists
}

func ExampleNotNilV() {
	rv := reflect.ValueOf(struct {
		V0 any
		V7 error
		V8 reflect.Type
		V1 chan error
		V2 func()
		V3 *int
		V4 unsafe.Pointer
		V5 []int
		V6 map[string]int
	}{})
	for i := range rv.NumField() {
		func(v any) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			must.NotNilV(v)
		}(rv.Field(i).Interface())
	}
	fmt.Println(must.NotNilV(1))
	fmt.Println(*must.NotNilV(new(int)))

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	must.NotNilWrap((*int)(nil), "business message %v", 100)

	// Output:
	// must not nil, but got invalid value
	// must not nil, but got invalid value
	// must not nil, but got invalid value
	// must not nil for type `chan error`
	// must not nil for type `func()`
	// must not nil for type `*int`
	// must not nil for type `unsafe.Pointer`
	// must not nil for type `[]int`
	// must not nil for type `map[string]int`
	// 1
	// 0
	// must not nil for type `*int` business message 100
}
