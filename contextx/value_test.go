package contextx_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/xoctopus/x/contextx"
)

type key struct{}

func getValue(ctx context.Context) bool {
	v := ctx.Value(key{})
	_ = v
	return true
}

func BenchmarkWithValue(b *testing.B) {
	parent := context.Background()

	b.Run("std.Context", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx := context.WithValue(parent, key{}, nil)
			getValue(ctx)
		}
	})
	b.Run("x.Contextx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx := contextx.WithValue(parent, key{}, nil)
			getValue(ctx)
		}
	})
}

func TestWithValue(t *testing.T) {
	t.Run("CatchParentIsNil", func(t *testing.T) {
		defer func() {
			t.Log(recover())
		}()
		contextx.WithValue(nil, nil, nil)
	})
	t.Run("CatchKeyIsNil", func(t *testing.T) {
		defer func() {
			t.Log(recover())
		}()
		contextx.WithValue(context.Background(), nil, nil)
	})
}

type fakeContext struct{}

func (fakeContext) Deadline() (time.Time, bool) { return time.Time{}, false }

func (fakeContext) Done() <-chan struct{} { return nil }

func (fakeContext) Err() error { return nil }

func (fakeContext) Value(any) any { return nil }

func ExampleWithValue() {
	ctx := contextx.WithValue(context.Background(), key{}, 100)

	fmt.Println(ctx.Value(key{}))
	fmt.Println(ctx.Value(1))
	fmt.Println(ctx)

	// Output:
	// 100
	// <nil>
	// context.Background.WithValue(type contextx_test.key, val<not Stringer>)
}

func ExampleWithContextCompose() {
	with := []contextx.WithContext{
		func(ctx context.Context) context.Context {
			return contextx.WithValue(ctx, key{}, 100)
		},
		func(ctx context.Context) context.Context {
			return contextx.WithValue(ctx, "key", "200")
		},
		func(ctx context.Context) context.Context {
			return contextx.WithValue(ctx, "key", reflect.ValueOf(1))
		},
	}

	compose := contextx.WithContextCompose(with...)

	ctx := compose(fakeContext{})
	fmt.Println(ctx.Value(key{}))
	fmt.Println(ctx)

	ctx = compose(context.Background())
	fmt.Println(ctx.Value(key{}))
	fmt.Println(ctx)

	// Output:
	// 100
	// contextx_test.fakeContext.WithValue(type contextx_test.key, val<not Stringer>).WithValue(type string, val200).WithValue(type string, val<int Value>)
	// 100
	// context.Background.WithValue(type contextx_test.key, val<not Stringer>).WithValue(type string, val200).WithValue(type string, val<int Value>)
}
