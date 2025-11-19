package syncx_test

import (
	"fmt"

	"github.com/xoctopus/x/syncx"
)

func ExampleOnceOverride() {
	v := syncx.NewOnceOverride("default")
	fmt.Println(v.Value())
	v.Set("override")
	fmt.Println(v.Value())
	v.Set("override2")
	fmt.Println(v.Value())

	// Output:
	// default
	// override
	// override
}
