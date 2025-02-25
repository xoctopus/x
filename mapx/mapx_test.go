package mapx_test

import (
	"reflect"
	"sync"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/ptrx"
)

func results(v ...any) []any {
	return v
}

func TestXmap_Load(t *testing.T) {
	k1 := ptrx.Ptr(1)
	k2 := ptrx.Ptr(2)
	k3 := ptrx.Ptr(3)

	NewWithT(t).Expect(k1).To(Equal(ptrx.Ptr(1)))
	NewWithT(t).Expect(k1).To(BeEquivalentTo(ptrx.Ptr(1)))   // called reflect.DeepEqual
	NewWithT(t).Expect(k1).NotTo(BeIdenticalTo(ptrx.Ptr(1))) // simple v1 == v2

	m1 := mapx.NewXmap(map[*int]int{}, true)
	m2 := mapx.NewSmap(&sync.Map{})

	NewWithT(t).Expect(m1.Exists(k1)).To(BeFalse())
	NewWithT(t).Expect(m2.Exists(k1)).To(Equal(m1.Exists(k1)))

	NewWithT(t).Expect(results(m1.Load(k1))).To(Equal([]any{0, false}))
	NewWithT(t).Expect(results(m2.Load(k1))).To(Equal([]any{nil, false}))

	m1.Store(k1, 100)
	m2.Store(k1, 100)

	NewWithT(t).Expect(m1.BatchLoad(k1, k2)).To(Equal([]mapx.Result[int]{{100, true}, {}}))
	NewWithT(t).Expect(m2.BatchLoad(k1, k2)).To(Equal([]mapx.Result[any]{{100, true}, {}}))

	m1.BatchStore([]*int{k1, k2}, []int{101, 102})
	m2.BatchStore([]any{k1, k2}, []any{101, 102})

	NewWithT(t).Expect(results(m1.LoadOrStore(k1, 100))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m2.LoadOrStore(k1, 100))).To(Equal([]any{101, true}))

	NewWithT(t).Expect(results(m1.LoadOrStore(k3, 103))).To(Equal([]any{103, false}))
	NewWithT(t).Expect(results(m2.LoadOrStore(k3, 103))).To(Equal([]any{103, false}))

	NewWithT(t).Expect(results(m1.LoadAndDelete(k3))).To(Equal([]any{103, true}))
	NewWithT(t).Expect(results(m2.LoadAndDelete(k3))).To(Equal([]any{103, true}))

	NewWithT(t).Expect(results(m1.LoadAndDelete(ptrx.Ptr(1)))).To(Equal([]any{0, false}))
	NewWithT(t).Expect(results(m2.LoadAndDelete(ptrx.Ptr(1)))).To(Equal([]any{nil, false}))

	m1.Delete(k2)
	m2.Delete(k2)
	m1.BatchDelete(ptrx.Ptr(1))
	m2.BatchDelete(ptrx.Ptr(1))

	NewWithT(t).Expect(results(m1.Swap(k1, 100))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m2.Swap(k1, 100))).To(Equal([]any{101, true}))

	NewWithT(t).Expect(results(m1.Swap(k1, 102))).To(Equal([]any{100, true}))
	NewWithT(t).Expect(results(m2.Swap(k1, 102))).To(Equal([]any{100, true}))

	NewWithT(t).Expect(m1.CompareAndSwap(k1, 101, 102)).To(BeFalse())
	NewWithT(t).Expect(m2.CompareAndSwap(k1, 101, 102)).To(BeFalse())

	NewWithT(t).Expect(m1.CompareAndSwap(k1, 102, 101)).To(BeTrue())
	NewWithT(t).Expect(m2.CompareAndSwap(k1, 102, 101)).To(BeTrue())

	NewWithT(t).Expect(m1.CompareAndDelete(k1, 100)).To(BeFalse())
	NewWithT(t).Expect(m2.CompareAndDelete(k1, 100)).To(BeFalse())

	NewWithT(t).Expect(m1.CompareAndDelete(k1, 101)).To(BeTrue())
	NewWithT(t).Expect(m2.CompareAndDelete(k1, 101)).To(BeTrue())

	m1.BatchStore([]*int{k1, k2}, []int{101, 102})
	m2.BatchStore([]any{k1, k2}, []any{101, 102})

	equal := func(k1, k2 any) bool {
		if k1 == k2 {
			return true
		}
		return reflect.DeepEqual(k1, k2)
	}

	k1_ := ptrx.Ptr(1)

	NewWithT(t).Expect(results(m1.LoadEq(k1, equal))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m2.LoadEq(k1, equal))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m1.LoadEq(k1_, equal))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m2.LoadEq(k1_, equal))).To(Equal([]any{101, true}))
	NewWithT(t).Expect(results(m1.LoadEq(ptrx.Ptr(5), equal))).To(Equal([]any{0, false}))
	NewWithT(t).Expect(results(m2.LoadEq(ptrx.Ptr(5), equal))).To(Equal([]any{nil, false}))

	m1.Store(k1_, 105)
	m2.Store(k1_, 105)

	NewWithT(t).Expect(m1.LoadEqs(k1, equal)).To(ConsistOf(101, 105))
	NewWithT(t).Expect(m2.LoadEqs(k1, equal)).To(ConsistOf(101, 105))
	NewWithT(t).Expect(m1.LoadEqs(ptrx.Ptr(1), equal)).To(ConsistOf(101, 105))
	NewWithT(t).Expect(m2.LoadEqs(ptrx.Ptr(1), equal)).To(ConsistOf(101, 105))
	NewWithT(t).Expect(m1.LoadEqs(ptrx.Ptr(5), equal)).To(HaveLen(0))
	NewWithT(t).Expect(m2.LoadEqs(ptrx.Ptr(5), equal)).To(HaveLen(0))

	keys1 := mapx.Keys(m1)
	keys2 := mapx.Keys(m2)
	NewWithT(t).Expect(keys1).To(ConsistOf(keys2...))

	values1 := mapx.Values(m1)
	values2 := mapx.Values(m2)
	NewWithT(t).Expect(values1).To(ConsistOf(values2...))

	NewWithT(t).Expect(mapx.Len(m1)).To(Equal(mapx.Len(m2)))

	m1.Clear()
	m2.Clear()
	NewWithT(t).Expect(mapx.Len(m1)).To(Equal(0))
	NewWithT(t).Expect(mapx.Len(m1)).To(Equal(mapx.Len(m2)))
}
