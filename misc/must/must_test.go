package must_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

func ExampleNoError() {
	must.NoError(nil)
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.NoError(errors.New("NoError: some error"))
	}()

	fmt.Println(must.NoErrorV(100, nil))
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.NoErrorV(100, errors.New("NoErrorV: some error"))
	}()

	must.NoErrorF(nil, "any")

	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.NoErrorF(errors.New("some error"), "NoErrorF: some message: %d", 10)
	}()

	// Output:
	// NoError: some error
	// 100
	// NoErrorV: some error
	// NoErrorF: some message: 10: some error
}

func ExampleBeTrue() {
	must.BeTrue(true)
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.BeTrue(false)
	}()

	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.BeTrueWrap(false, errors.New("BeTrueWrap: some error"))
	}()

	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.BeTrueWrapF(false, errors.New("some error"), "BeTrueWrapF: message and args: %d", 100)
	}()

	f1 := func() (int, bool) { return 100, true }
	f2 := func() (int, bool) { return 100, false }

	fmt.Println(must.BeTrueV(f1()))
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		_ = must.BeTrueV(f2())
	}()

	must.BeTrueF(true, "any")
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		must.BeTrueF(false, "BeTrueF: required exists")
	}()

	// Output:
	// must be true
	// must be true: [err: BeTrueWrap: some error]
	// must be true: BeTrueWrapF: message and args: 100 [err: some error]
	// 100
	// must be true
	// must be true: BeTrueF: required exists
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
				fmt.Println(recover())
			}()
			must.NotNilV(v)
		}(rv.Field(i).Interface())

		func(v any) {
			defer func() {
				fmt.Println(recover())
			}()
			must.NotNilF(v, "business message: %v", 100)
		}(rv.Field(i).Interface())

		func(v any) {
			defer func() {
				fmt.Println(recover())
			}()
			must.NotNilWrap(v, errors.Errorf("business message: %v", 101))
		}(rv.Field(i).Interface())

		func(v any) {
			defer func() {
				fmt.Println(recover())
			}()
			must.NotNilWrapF(v, errors.Errorf("business message: %v", 101), "more custom message and args: %d", 102)
		}(rv.Field(i).Interface())
	}
	fmt.Println(must.NotNilV(1))
	fmt.Println(*must.NotNilV(new(int)))

	// Output:
	// must not nil, but got invalid value.
	// must not nil, but got invalid value. business message: 100
	// must not nil, but got invalid value. [err: business message: 101]
	// must not nil, but got invalid value. more custom message and args: 102 [err: business message: 101]
	// must not nil, but got invalid value.
	// must not nil, but got invalid value. business message: 100
	// must not nil, but got invalid value. [err: business message: 101]
	// must not nil, but got invalid value. more custom message and args: 102 [err: business message: 101]
	// must not nil, but got invalid value.
	// must not nil, but got invalid value. business message: 100
	// must not nil, but got invalid value. [err: business message: 101]
	// must not nil, but got invalid value. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `chan error`.
	// must not nil for type `chan error`. business message: 100
	// must not nil for type `chan error`. [err: business message: 101]
	// must not nil for type `chan error`. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `func()`.
	// must not nil for type `func()`. business message: 100
	// must not nil for type `func()`. [err: business message: 101]
	// must not nil for type `func()`. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `*int`.
	// must not nil for type `*int`. business message: 100
	// must not nil for type `*int`. [err: business message: 101]
	// must not nil for type `*int`. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `unsafe.Pointer`.
	// must not nil for type `unsafe.Pointer`. business message: 100
	// must not nil for type `unsafe.Pointer`. [err: business message: 101]
	// must not nil for type `unsafe.Pointer`. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `[]int`.
	// must not nil for type `[]int`. business message: 100
	// must not nil for type `[]int`. [err: business message: 101]
	// must not nil for type `[]int`. more custom message and args: 102 [err: business message: 101]
	// must not nil for type `map[string]int`.
	// must not nil for type `map[string]int`. business message: 100
	// must not nil for type `map[string]int`. [err: business message: 101]
	// must not nil for type `map[string]int`. more custom message and args: 102 [err: business message: 101]
	// 1
	// 0
}

func ExampleSuccess() {
	x := must.Success(func() (int, error) { return 100, nil })
	fmt.Println(x)

	x = must.OK(func() (int, bool) { return 100, true })
	fmt.Println(x)

	func() {
		defer func() {
			fmt.Println(recover())
		}()
		_ = must.Success(func() (int, error) { return 0, errors.New("some error") })
	}()
	func() {
		defer func() {
			fmt.Println(recover())
		}()
		_ = must.OK(func() (int, bool) { return 0, false })
	}()

	// Output:
	// 100
	// 100
	// some error
	// must be true
}

func TestIdenticalTypes(t *testing.T) {
	NewWithT(t).Expect(must.IdenticalTypes(1, 1)).To(BeTrue())
}
