package contextx_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/xoctopus/x/contextx"
	. "github.com/xoctopus/x/testx"
)

type key struct{}

func BenchmarkWithValue(b *testing.B) {
	parent := context.Background()

	b.Run("WithValue", func(b *testing.B) {
		b.Run("std.Context", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = context.WithValue(parent, key{}, "value")
			}
		})
		b.Run("x.Contextx", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = contextx.WithValue(parent, key{}, "value")
			}
		})
	})

	b.Run("Value", func(b *testing.B) {
		b.Run("std.Context", func(b *testing.B) {
			for b.Loop() {
				_ = parent.Value(key{})
			}
		})
		b.Run("x.Contextx", func(b *testing.B) {
			c := contextx.NewT[key]()
			for b.Loop() {
				_, _ = c.From(parent)
			}
		})
	})

}

func TestWithValue(t *testing.T) {
	t.Run("CatchParentIsNil", func(t *testing.T) {
		ExpectPanic(t, func() { contextx.WithValue(nil, nil, nil) }, NotBeNil[any]())
	})
	t.Run("CatchKeyIsNil", func(t *testing.T) {
		ExpectPanic(t, func() { contextx.WithValue(context.Background(), nil, nil) }, NotBeNil[any]())
	})

	ctx := contextx.WithValue(context.Background(), key{}, t.Name())
	Expect(t, ctx.Value(key{}).(string), Equal(t.Name()))
}

type MockContext struct{}

func (MockContext) Deadline() (time.Time, bool) { return time.Time{}, true }

func (MockContext) Err() error { return nil }

func (MockContext) Done() <-chan struct{} { return nil }

func (MockContext) Value(any) any { return nil }

func ExampleWithValue() {
	var ctx context.Context

	type key1 struct{}
	type key2 struct{}
	type key3 struct{}
	type key4 struct{}

	ctx = contextx.WithValue(MockContext{}, key1{}, "1")
	fmt.Println(ctx)
	ctx = contextx.WithValue(ctx, key2{}, "2")
	fmt.Println(ctx)
	ctx = contextx.WithValue(ctx, key3{}, net.IPv4(1, 1, 1, 1))
	fmt.Println(ctx)
	ctx = contextx.WithValue(ctx, key4{}, 4)
	fmt.Println(ctx)

	fmt.Println("context value key1:", ctx.Value(key1{}))
	fmt.Println("context value key2:", ctx.Value(key2{}))
	fmt.Println("context value key3:", ctx.Value(key3{}))
	fmt.Println("context value key4:", ctx.Value(key4{}))

	// Output:
	// contextx_test.MockContext.WithValue(key:contextx_test.key1, val:1)
	// contextx_test.MockContext.WithValue(key:contextx_test.key1, val:1).WithValue(key:contextx_test.key2, val:2)
	// contextx_test.MockContext.WithValue(key:contextx_test.key1, val:1).WithValue(key:contextx_test.key2, val:2).WithValue(key:contextx_test.key3, val:1.1.1.1)
	// contextx_test.MockContext.WithValue(key:contextx_test.key1, val:1).WithValue(key:contextx_test.key2, val:2).WithValue(key:contextx_test.key3, val:1.1.1.1).WithValue(key:contextx_test.key4, val:<not Stringer>)
	// context value key1: 1
	// context value key2: 2
	// context value key3: 1.1.1.1
	// context value key4: 4
}
