package syncx_test

import (
	"testing"

	"github.com/xoctopus/x/syncx"
	. "github.com/xoctopus/x/testx"
)

type Bytes struct {
	data []byte
}

func (b *Bytes) Reset() {
	b.data = b.data[:0]
}

func TestNewPool(t *testing.T) {
	t.Run("Slice", func(t *testing.T) {
		p := syncx.NewPool(func() []byte { return make([]byte, 1024) })
		v := p.Get()
		defer p.Put(v)

		Expect(t, len(v), Equal(0))
		Expect(t, cap(v), Equal(1024))
	})

	t.Run("Reset", func(t *testing.T) {
		p := syncx.NewPool(func() *Bytes {
			return &Bytes{
				data: make([]byte, 1024),
			}
		})
		v := p.Get()
		defer p.Put(v)

		Expect(t, len(v.data), Equal(0))
		Expect(t, cap(v.data), Equal(1024))
	})
}
