package ptrx_test

import (
	"fmt"
	"time"

	"github.com/xoctopus/x/ptrx"
)

func Example() {
	fmt.Println(*new(10))
	fmt.Println(*new(uint64(10)))
	fmt.Println(*new("abc"))
	fmt.Println(ptrx.Ptr(time.Second))
	fmt.Println(ptrx.Ptr(time.Hour))

	// Output:
	// 10
	// 10
	// abc
	// 1s
	// 1h0m0s
}
