package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/testx"
)

type Struct struct {
	cannotSet int
	CanSet    int
	Any       any
	Any2      any
	Ptr       *int
	Map       map[string]int
	Slice     []any
	Array     [5]any
	Ref       *Struct
	Func      func()
	Chan      chan int
}

type String string

func TestDeepCopy(t *testing.T) {
	t.Run("InvalidDstValue", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "invalid dst value cannot set")
		}()
		dst := reflect.ValueOf(1)
		NewWithT(t).Expect(dst.CanSet()).To(BeFalse())
		DeepCopy(reflect.ValueOf(1), reflect.ValueOf(10))
	})
	t.Run("InvalidSrcValue", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "invalid src value")
		}()
		var i any
		src := reflect.ValueOf(i)
		NewWithT(t).Expect(src.IsValid()).To(BeFalse())
		DeepCopy(reflect.ValueOf(&Struct{}).Elem().Field(1), reflect.ValueOf(i))
	})
	t.Run("SameValues", func(t *testing.T) {
		v := reflect.ValueOf(&Struct{}).Elem().Field(1)
		DeepCopy(v, v)
	})
	t.Run("CannotAssignableTo", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "src type cannot assign to dst")
		}()
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Array")
		v2 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		DeepCopy(v1, v2)
	})
	t.Run("SetZeroValue", func(t *testing.T) {
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		v2 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		DeepCopy(v1, v2)
		NewWithT(t).Expect(v1.IsZero()).To(BeTrue())
	})
	t.Run("UnsupportedTypeChan", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "chan type cannot be copied")
		}()
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Any")
		v2 := reflect.ValueOf(&Struct{Chan: make(chan int)}).Elem().FieldByName("Chan")
		DeepCopy(v1, v2)
	})
	t.Run("UnsupportedTypeFunc", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "func type cannot be copied")
		}()
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Any")
		v2 := reflect.ValueOf(&Struct{Func: func() {}}).Elem().FieldByName("Func")
		DeepCopy(v1, v2)
	})
	t.Run("StructContainsUnexportedField", func(t *testing.T) {
		defer func() {
			NewWithT(t).Expect(testx.Recover(recover())).To(ContainSubstring("cannot be copied"))
		}()
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Any")
		v2 := reflect.ValueOf(&Struct{Any: struct{ a string }{a: "a"}}).Elem().FieldByName("Any")
		DeepCopy(v1, v2)
	})

	src := []any{
		1,
		ptrx.Ptr(100),
		[1]any{1},
		struct{ A string }{A: "100"},
		nil,
		map[int]int{1: 1},
	}
	// src = append(src, src)
	dst := Clone(src)
	NewWithT(t).Expect(reflect.DeepEqual(src, dst)).To(BeTrue())
}
