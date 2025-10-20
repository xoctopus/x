package mapx_test

/*
import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/testx"
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

	Expect(t, k1, Equal(ptrx.Ptr(1)))
	Expect(t, k1 == ptrx.Ptr(1), BeFalse())

	for _, m := range []mapx.Map[any, any]{
		mapx.NewXmap[any, any](),
		mapx.NewSafeXmap[any, any](),
		mapx.Wrap(map[any]any{}),
		mapx.SafeWrap(map[any]any{}),
		mapx.NewSmap[any, any](),
	} {
		Expect(t, m.Exists(k1), BeFalse())
		Expect(t, results(m.Load(k1)), Equal([]any{nil, false}))

		m.Store(k1, 100)
		Expect(t, m.BatchLoad(k1, k2), Equal([]mapx.Result[any]{{100, true}, {}}))

		m.BatchStore([]any{k1, k2}, []any{101, 102})
		Expect(t, results(m.LoadOrStore(k1, 100)), Equal([]any{101, true}))
		Expect(t, results(m.LoadOrStore(k3, 103)), Equal([]any{103, false}))
		Expect(t, results(m.LoadAndDelete(k3)), Equal([]any{103, true}))
		Expect(t, results(m.LoadAndDelete(ptrx.Ptr(1))), Equal([]any{nil, false}))

		m.Delete(ptrx.Ptr(1))
		m.BatchDelete(ptrx.Ptr(1), k2)
		Expect(t, results(m.Swap(k1, 100)), Equal([]any{101, true}))
		Expect(t, results(m.Swap(k1, 102)), Equal([]any{100, true}))

		Expect(t, m.CompareAndSwap(k1, 101, 102), BeFalse())
		Expect(t, m.CompareAndSwap(k1, 102, 101), BeTrue())
		Expect(t, m.CompareAndDelete(k1, 100), BeFalse())
		Expect(t, m.CompareAndDelete(k1, 101), BeTrue())

		m.BatchStore([]any{k1, k2}, []any{101, 102})

		Expect(t, results(m.LoadEq(equal(k1))), Equal([]any{101, true}))
		Expect(t, results(m.LoadEq(equal(ptrx.Ptr(1)))), Equal([]any{101, true}))
		Expect(t, results(m.LoadEq(equal(ptrx.Ptr(5)))), Equal([]any{nil, false}))

		k := ptrx.Ptr(1)
		m.Store(k, 105)
		Expect(t, m.LoadEqs(equal(k1)), ConsistOf(101, 105))
		Expect(t, m.LoadEqs(equal(k)), ConsistOf(101, 105))
		Expect(t, m.LoadEqs(equal(ptrx.Ptr(5))), HaveLen[[]any](0))

		Expect(t, mapx.Keys(m), ConsistOf(k1, k2, k))
		Expect(t, mapx.Values(m), ConsistOf(101, 102, 105))

		Expect(t, mapx.Len(m), Equal(3))

		mm := m.Clone()
		Expect(t, mapx.Len(m), Equal(mapx.Len(mm)))
		Expect(t, mapx.Equal(m, mm), BeTrue())

		m.Clear()
		Expect(t, mapx.Len(m), Equal(0))
	}

	t.Run("Equal", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			Expect(t, mapx.Equal[any, any](nil, nil), BeTrue())
		})

		m1 := mapx.NewSafeXmap[int, int]()
		m1.BatchStore([]int{1, 2, 3}, []int{1, 2, 3})
		m2 := mapx.NewSmap[int, int]()
		m2.BatchStore([]int{1, 2, 3, 4}, []int{1, 2, 3, 4})

		Expect(t, mapx.Equal(m1, m2), BeFalse())
		Expect(t, mapx.Equal(m2, m1), BeFalse())

		m1.Store(4, 4)
		Expect(t, mapx.Equal(m1, m2), BeTrue())
		Expect(t, mapx.Equal(m2, m1), BeTrue())
	})
}

func TestSet(t *testing.T) {
	for _, set := range []mapx.Set[int]{mapx.NewSet[int](), mapx.NewSafeSet[int]()} {
		set.Store(1, 2, 3)
		Expect(t, set.Exists(1), BeTrue())
		set.Delete(1)
		Expect(t, set.Exists(1), BeFalse())
		// Expect(t, set.Keys(), ConsistOf(2, 3))
		Expect(t, set.Len(), Equal(2))

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
		Expect(t, has, BeTrue())

		has = false
		expect = 1
		set.Range(f(expect, &has))
		Expect(t, has, BeFalse())

		set2 := set.Clone()
		Expect(t, set.Equal(set2), BeTrue())

		set.Clear()
		Expect(t, set.Len(), Equal(0))
	}
}

*/
