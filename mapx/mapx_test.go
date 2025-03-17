package mapx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/ptrx"
)

func results(v ...any) []any {
	return v
}

func TestMap(t *testing.T) {
	k1 := ptrx.Ptr(1)
	k2 := ptrx.Ptr(2)
	k3 := ptrx.Ptr(3)
	equal := func(x any) func(k any) bool {
		return func(k any) bool {
			if k == x {
				return true
			}
			return reflect.DeepEqual(x, k)
		}
	}

	NewWithT(t).Expect(k1).To(Equal(ptrx.Ptr(1)))
	NewWithT(t).Expect(k1).To(BeEquivalentTo(ptrx.Ptr(1)))   // called reflect.DeepEqual
	NewWithT(t).Expect(k1).NotTo(BeIdenticalTo(ptrx.Ptr(1))) // simple v1 == v2

	for _, m := range []mapx.Map[any, any]{
		mapx.NewXmap[any, any](),
		mapx.NewSafeXmap[any, any](),
		mapx.Wrap(map[any]any{}),
		mapx.SafeWrap(map[any]any{}),
		mapx.NewSmap[any, any](),
	} {
		NewWithT(t).Expect(m.Exists(k1)).To(BeFalse())
		NewWithT(t).Expect(results(m.Load(k1))).To(Equal([]any{nil, false}))

		m.Store(k1, 100)
		NewWithT(t).Expect(m.BatchLoad(k1, k2)).To(Equal([]mapx.Result[any]{{100, true}, {}}))

		m.BatchStore([]any{k1, k2}, []any{101, 102})
		NewWithT(t).Expect(results(m.LoadOrStore(k1, 100))).To(Equal([]any{101, true}))
		NewWithT(t).Expect(results(m.LoadOrStore(k3, 103))).To(Equal([]any{103, false}))

		NewWithT(t).Expect(results(m.LoadAndDelete(k3))).To(Equal([]any{103, true}))
		NewWithT(t).Expect(results(m.LoadAndDelete(ptrx.Ptr(1)))).To(Equal([]any{nil, false}))

		m.Delete(ptrx.Ptr(1))
		m.BatchDelete(ptrx.Ptr(1), k2)
		NewWithT(t).Expect(results(m.Swap(k1, 100))).To(Equal([]any{101, true}))
		NewWithT(t).Expect(results(m.Swap(k1, 102))).To(Equal([]any{100, true}))

		NewWithT(t).Expect(m.CompareAndSwap(k1, 101, 102)).To(BeFalse())
		NewWithT(t).Expect(m.CompareAndSwap(k1, 102, 101)).To(BeTrue())
		NewWithT(t).Expect(m.CompareAndDelete(k1, 100)).To(BeFalse())
		NewWithT(t).Expect(m.CompareAndDelete(k1, 101)).To(BeTrue())

		m.BatchStore([]any{k1, k2}, []any{101, 102})

		NewWithT(t).Expect(results(m.LoadEq(equal(k1)))).To(Equal([]any{101, true}))
		NewWithT(t).Expect(results(m.LoadEq(equal(ptrx.Ptr(1))))).To(Equal([]any{101, true}))
		NewWithT(t).Expect(results(m.LoadEq(equal(ptrx.Ptr(5))))).To(Equal([]any{nil, false}))

		k := ptrx.Ptr(1)
		m.Store(k, 105)
		NewWithT(t).Expect(m.LoadEqs(equal(k1))).To(ConsistOf(101, 105))
		NewWithT(t).Expect(m.LoadEqs(equal(k))).To(ConsistOf(101, 105))
		NewWithT(t).Expect(m.LoadEqs(equal(ptrx.Ptr(5)))).To(HaveLen(0))

		NewWithT(t).Expect(mapx.Keys(m)).To(ConsistOf(k1, k2, k))
		NewWithT(t).Expect(mapx.Values(m)).To(ConsistOf(101, 102, 105))

		NewWithT(t).Expect(mapx.Len(m)).To(Equal(3))

		mm := m.Clone()
		NewWithT(t).Expect(mapx.Len(m)).To(Equal(mapx.Len(mm)))
		NewWithT(t).Expect(mapx.Equal(m, mm)).To(BeTrue())

		m.Clear()
		NewWithT(t).Expect(mapx.Len(m)).To(Equal(0))
	}

	t.Run("Equal", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			NewWithT(t).Expect(mapx.Equal[any, any](nil, nil)).To(BeTrue())
		})

		m1 := mapx.NewSafeXmap[int, int]()
		m1.BatchStore([]int{1, 2, 3}, []int{1, 2, 3})
		m2 := mapx.NewSmap[int, int]()
		m2.BatchStore([]int{1, 2, 3, 4}, []int{1, 2, 3, 4})

		NewWithT(t).Expect(mapx.Equal(m1, m2)).To(BeFalse())
		NewWithT(t).Expect(mapx.Equal(m2, m1)).To(BeFalse())

		m1.Store(4, 4)
		NewWithT(t).Expect(mapx.Equal(m1, m2)).To(BeTrue())
		NewWithT(t).Expect(mapx.Equal(m2, m1)).To(BeTrue())
	})
}

func TestSet(t *testing.T) {
	for _, set := range []mapx.Set[int]{mapx.NewSet[int](), mapx.NewSafeSet[int]()} {
		set.Store(1, 2, 3)
		NewWithT(t).Expect(set.Exists(1)).To(BeTrue())
		set.Delete(1)
		NewWithT(t).Expect(set.Exists(1)).To(BeFalse())
		NewWithT(t).Expect(set.Keys()).To(ConsistOf(2, 3))
		NewWithT(t).Expect(set.Len()).To(Equal(2))

		f := func(expect int, has *bool) func(k int) bool {
			return func(k int) bool {
				if k == expect {
					*has = true
					return false
				}
				return true
			}
		}

		has := false
		expect := 2
		set.Range(f(expect, &has))
		NewWithT(t).Expect(has).To(BeTrue())

		has = false
		expect = 1
		set.Range(f(expect, &has))
		NewWithT(t).Expect(has).To(BeFalse())

		set2 := set.Clone()
		NewWithT(t).Expect(set.Equal(set2)).To(BeTrue())

		set.Clear()
		NewWithT(t).Expect(set.Len()).To(Equal(0))
	}
}
