package contextx_test

import (
	"context"
	"fmt"

	"github.com/xoctopus/x/contextx"
)

type tCtxT struct{}

func Example() {
	ctx := contextx.Carry[tCtxT]("100")(context.Background())
	fmt.Println(contextx.Must[tCtxT, string](ctx))
	fmt.Println(contextx.FromOr[tCtxT](ctx, "101"))

	ctx = context.Background()
	fmt.Println(contextx.FromOr[tCtxT](ctx, "102"))

	// Output:
	// 100
	// 100
	// 102
}
