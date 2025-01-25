package contextx_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/contextx"
)

func TestContext(t *testing.T) {
	t.Run("EmptyContext", func(t *testing.T) {
		ctx := context.Background()
		t.Run("NoDefaults", func(t *testing.T) {
			t.Run("From", func(t *testing.T) {
				val, ok := contextx.NewT[string]().From(ctx)
				NewWithT(t).Expect(val).To(Equal(""))
				NewWithT(t).Expect(ok).To(BeFalse())
			})
			t.Run("MustFrom", func(t *testing.T) {
				defer func() {
					NewWithT(t).Expect(recover()).NotTo(BeNil())
				}()
				_ = contextx.NewT[string]().MustFrom(ctx)
			})
		})
		t.Run("HasDefaults", func(t *testing.T) {
			t.Run("From", func(t *testing.T) {
				val, ok := contextx.NewV(t.Name()).From(ctx)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(val).To(Equal(t.Name()))
			})
			t.Run("MustFrom", func(t *testing.T) {
				val := contextx.NewT(contextx.WithDefault(t.Name())).MustFrom(ctx)
				NewWithT(t).Expect(val).To(Equal(t.Name()))
			})
		})
	})
	t.Run("ValueContext", func(t *testing.T) {
		c := contextx.NewT[string]()
		ctx := c.With(context.Background(), t.Name())

		val, ok := c.From(ctx)
		NewWithT(t).Expect(ok).To(BeTrue())
		NewWithT(t).Expect(val).To(Equal(t.Name()))

		val = c.MustFrom(ctx)
		NewWithT(t).Expect(val).To(Equal(t.Name()))
	})
}

func ExampleCompose() {
	c1 := contextx.NewT[string]()
	c2 := contextx.NewT[int]()
	c3 := contextx.NewT[fmt.Stringer]()

	ctx := contextx.Compose(
		c1.WithCompose("1"),
		c2.WithCompose(2),
		c3.WithCompose(net.IPv4(1, 1, 1, 1)),
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
