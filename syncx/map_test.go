package syncx_test

import (
	"testing"

	"github.com/xoctopus/x/ptrx"
	"github.com/xoctopus/x/syncx"
	. "github.com/xoctopus/x/testx"
)

func results(v ...any) []any {
	return v
}

func TestMap(t *testing.T) {
	k1 := ptrx.Ptr(1)
	k2 := ptrx.Ptr(2)
	// k3 := ptrx.Ptr(3)
	// equal := func(x any) func(k any) bool {
	// 	return func(k any) bool {
	// 		if k == x {
	// 			return true
	// 		}
	// 		return reflect.DeepEqual(x, k)
	// 	}
	// }

	Expect(t, k1, Equal(ptrx.Ptr(1)))
	Expect(t, k1 == ptrx.Ptr(1), BeFalse())

	for _, m := range []syncx.Map[any, any]{
		syncx.NewXmap[any, any](),
		syncx.NewSmap[any, any](),
		syncx.AsSmap(map[any]any{}),
		syncx.AsXmap(map[any]any{}),
	} {
		Expect(t, results(m.Load(k1)), Equal([]any{nil, false}))

		m.Store(k1, 100)
		Expect(t, results(m.Load(k1)), Equal([]any{100, true}))
		Expect(t, results(m.Load(k2)), Equal([]any{nil, false}))

		m.Delete(ptrx.Ptr(1))
		Expect(t, results(m.Load(k1)), Equal([]any{100, true}))
		Expect(t, results(m.Load(k2)), Equal([]any{nil, false}))

		Expect(t, m.CompareAndSwap(k1, 101, 102), BeFalse())
		Expect(t, m.CompareAndSwap(k1, 100, 102), BeTrue())

		Expect(t, m.CompareAndDelete(k1, 100), BeFalse())
		Expect(t, m.CompareAndDelete(k1, 102), BeTrue())

		m.Store(k1, 100)
		Expect(t, results(m.Swap(k1, 102)), Equal([]any{100, true}))

		m.Delete(k1)
		Expect(t, results(m.Load(k1)), Equal([]any{nil, false}))

		m.Store(k1, 100)
		Expect(t, results(m.LoadOrStore(k1, 101)), Equal([]any{100, true}))
		Expect(t, results(m.LoadOrStore(k2, 101)), Equal([]any{101, false}))

		count := 0
		counter := func(k any, v any) bool {
			count++
			return true
		}
		m.Range(counter)
		Expect(t, count, Equal(2))

		Expect(t, results(m.LoadAndDelete(k1)), Equal([]any{100, true}))

		count = 0
		m.Range(counter)
		Expect(t, count, Equal(1))

		m.Clear()
		count = 0
		m.Range(counter)
		Expect(t, count, Equal(0))
	}

	m := syncx.AsSmap(map[any]any{k1: 100})
	Expect(t, results(m.Load(k1)), Equal([]any{100, true}))

	m = syncx.AsXmap(map[any]any{k1: 100})
	m.Range(func(k any, v any) bool {
		if v == 100 {
			x, ok := k.(*int)
			Expect(t, ok, BeTrue())
			Expect(t, x, Equal(k1))
			return false
		}
		return true
	})
}

// func TestSet(t *testing.T) {
// 	for _, set := range []mapx.Set[int]{mapx.NewSet[int](), mapx.NewSafeSet[int]()} {
// 		set.Store(1, 2, 3)
// 		NewWithT(t).Expect(set.Exists(1)).To(BeTrue())
// 		set.Delete(1)
// 		NewWithT(t).Expect(set.Exists(1)).To(BeFalse())
// 		NewWithT(t).Expect(set.Keys()).To(ConsistOf(2, 3))
// 		NewWithT(t).Expect(set.Len()).To(Equal(2))
//
// 		f := func(expect int, has *bool) func(k int) bool {
// 			return func(k int) bool {
// 				if k == expect {
// 					*has = true
// 					return false
// 				}
// 				return true
// 			}
// 		}
//
// 		has := false
// 		expect := 2
// 		set.Range(f(expect, &has))
// 		NewWithT(t).Expect(has).To(BeTrue())
//
// 		has = false
// 		expect = 1
// 		set.Range(f(expect, &has))
// 		NewWithT(t).Expect(has).To(BeFalse())
//
// 		set2 := set.Clone()
// 		NewWithT(t).Expect(set.Equal(set2)).To(BeTrue())
//
// 		set.Clear()
// 		NewWithT(t).Expect(set.Len()).To(Equal(0))
// 	}
// }
