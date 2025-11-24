package contextx_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/xoctopus/x/contextx"
	. "github.com/xoctopus/x/testx"
)

func TestContext(t *testing.T) {
	t.Run("EmptyContext", func(t *testing.T) {
		ctx := context.Background()
		t.Run("NoDefaults", func(t *testing.T) {
			t.Run("From", func(t *testing.T) {
				val, ok := contextx.NewT[string]().From(ctx)
				Expect(t, val, Equal(""))
				Expect(t, ok, BeFalse())
			})
			t.Run("MustFrom", func(t *testing.T) {
				ExpectPanic[error](t, func() {
					_ = contextx.NewT[string]().MustFrom(ctx)
				})
			})
		})
		t.Run("HasDefaults", func(t *testing.T) {
			t.Run("From", func(t *testing.T) {
				val, ok := contextx.NewV(t.Name()).From(ctx)
				Expect(t, ok, BeTrue())
				Expect(t, val, Equal(t.Name()))
			})
			t.Run("MustFrom", func(t *testing.T) {
				val := contextx.NewT(contextx.WithDefault(t.Name())).MustFrom(ctx)
				Expect(t, val, Equal(t.Name()))
			})
		})
	})
	t.Run("ValueContext", func(t *testing.T) {
		c := contextx.NewT[string]()
		ctx := c.With(context.Background(), t.Name())

		val, ok := c.From(ctx)
		Expect(t, ok, BeTrue())
		Expect(t, val, Equal(t.Name()))

		Expect(t, c.MustFrom(ctx), Equal(t.Name()))

		t.Run("Once", func(t *testing.T) {
			c = contextx.NewT[string](contextx.WithOnce[string]())

			ctx := c.With(ctx, t.Name())
			Expect(t, c.MustFrom(ctx), Equal(t.Name()))

			ctx = c.With(context.Background(), "replace")
			Expect(t, c.MustFrom(ctx), Equal(t.Name()))
		})
	})
}

func ExampleCompose() {
	c1 := contextx.NewT[string]()
	c2 := contextx.NewT[int]()
	c3 := contextx.NewT[fmt.Stringer]()

	ctx := contextx.Compose(
		c1.Carry("1"),
		c2.Carry(2),
		c3.Carry(net.IPv4(1, 1, 1, 1)),
	)(context.Background())
	fmt.Println(ctx)
	fmt.Println(c1.MustFrom(ctx))
	fmt.Println(c2.MustFrom(ctx))
	fmt.Println(c3.MustFrom(ctx))

	// Output:
	// context.Background.WithValue(key:*contextx.ctx[string], val:1).WithValue(key:*contextx.ctx[int], val:<not Stringer>).WithValue(key:*contextx.ctx[fmt.Stringer], val:1.1.1.1)
	// 1
	// 2
	// 1.1.1.1
}
