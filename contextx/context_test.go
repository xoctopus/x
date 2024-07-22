package contextx_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/contextx"
)

type Value struct {
	Int int
}

var ValueContext = contextx.NewValue(&Value{1})

func TestContext(t *testing.T) {
	empty := context.Background()

	t.Run("FailedToExtract", func(t *testing.T) {
		t.Run("MustFrom", func(t *testing.T) {
			defer func() {
				err := recover().(error)
				NewWithT(t).Expect(err.Error()).To(ContainSubstring("not found in context"))
			}()

			ctx := contextx.New[*Value](nil)
			_ = ctx.MustFrom(empty)
		})

		t.Run("From", func(t *testing.T) {
			ctx := contextx.New[*Value](contextx.With[*Value](&Value{2}))
			v, ok := ctx.From(empty)
			NewWithT(t).Expect(ok).To(BeFalse())
			NewWithT(t).Expect(v).To(BeNil())
		})
	})

	t.Run("FromValuer", func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			t.Run("With", func(t *testing.T) {
				ctx := contextx.New(contextx.With(&Value{2}))
				v := ctx.MustFrom(empty)
				NewWithT(t).Expect(v).NotTo(BeNil())
				NewWithT(t).Expect(v.Int).To(Equal(2))
			})
			t.Run("WithValuer", func(t *testing.T) {
				ctx := contextx.New(contextx.WithValuer(func() *Value { return &Value{3} }))
				v := ctx.MustFrom(empty)
				NewWithT(t).Expect(v).NotTo(BeNil())
				NewWithT(t).Expect(v.Int).To(Equal(3))
			})
		})
		t.Run("NewValue", func(t *testing.T) {
			ctx := contextx.NewValue(&Value{4})
			v := ctx.MustFrom(empty)
			NewWithT(t).Expect(v).NotTo(BeNil())
			NewWithT(t).Expect(v.Int).To(Equal(4))
		})
	})

	t.Run("OverwriteByInjection", func(t *testing.T) {
		ctx := contextx.NewValue(&Value{5})
		root := ctx.With(empty, &Value{6})

		v, ok := ctx.From(root)
		NewWithT(t).Expect(ok).To(BeTrue())
		NewWithT(t).Expect(v.Int).To(Equal(6))

		v = ctx.MustFrom(root)
		NewWithT(t).Expect(v.Int).To(Equal(6))
	})
}

func BenchmarkCtx_MustFrom(b *testing.B) {
	empty := context.Background()
	injected := ValueContext.With(empty, &Value{1})

	b.Run("FromValuer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValueContext.MustFrom(empty)
		}
	})

	b.Run("ExtractFromContext", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValueContext.MustFrom(injected)
		}
	})

	overinjected := context.Background()
	for i := 0; i < 1000; i++ {
		overinjected = context.WithValue(overinjected, i, i)
	}
	overinjected = ValueContext.With(overinjected, &Value{1})

	b.Run("OverInjected", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValueContext.MustFrom(overinjected)
		}
	})

	overinjected2 := context.Background()
	for i := 0; i < 1000; i++ {
		overinjected = contextx.WithValue(overinjected, i, i)
	}
	overinjected2 = ValueContext.With(overinjected, &Value{1})

	b.Run("OverInjected2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValueContext.MustFrom(overinjected2)
		}
	})
}
