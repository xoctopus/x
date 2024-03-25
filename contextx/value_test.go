package contextx_test

import (
	"context"
	"testing"

	"github.com/sincospro/x/contextx"
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
