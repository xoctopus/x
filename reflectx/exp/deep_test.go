package exp_test

import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/reflectx/exp"
	. "github.com/xoctopus/x/testx"
)

type Struct struct {
	cannotSet  int
	CanSet     int
	Any        any
	Any2       any
	Ptr        *int
	Map        map[string]int
	Slice      []any
	Array      [5]any
	Ref        *Struct
	Func       func()
	Chan       chan int
	unexported struct {
		str string
	}
}

func TestDeepCopy(t *testing.T) {
	t.Run("InvalidDstValue", func(t *testing.T) {
		dst := reflect.ValueOf(1)
		Expect(t, dst.CanSet(), BeFalse())

		ExpectPanic(
			t,
			func() { DeepCopy(reflect.ValueOf(1), reflect.ValueOf(10)) },
			ErrorContains("invalid dst value cannot set"),
		)
	})
	t.Run("InvalidSrcValue", func(t *testing.T) {
		var i any
		src := reflect.ValueOf(i)

		Expect(t, src.IsValid(), BeFalse())
		ExpectPanic(
			t,
			func() { DeepCopy(reflect.ValueOf(&Struct{}).Elem().Field(1), reflect.ValueOf(i)) },
			ErrorContains("invalid src value"),
		)

	})
	t.Run("SameValues", func(t *testing.T) {
		v := reflect.ValueOf(&Struct{}).Elem().Field(1)
		DeepCopy(v, v)
	})
	t.Run("CannotAssignableTo", func(t *testing.T) {
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Array")
		v2 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		ExpectPanic(
			t,
			func() { DeepCopy(v1, v2) },
			ErrorContains("src type cannot assign to dst"),
		)
	})
	t.Run("SetZeroValue", func(t *testing.T) {
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		v2 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Slice")
		DeepCopy(v1, v2)
		Expect(t, v1.IsZero(), BeTrue())
	})
	t.Run("UnsupportedTypeChan", func(t *testing.T) {
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Any")
		v2 := reflect.ValueOf(&Struct{Chan: make(chan int)}).Elem().FieldByName("Chan")
		ExpectPanic(
			t,
			func() { DeepCopy(v1, v2) },
			ErrorContains("chan type cannot be copied"),
		)
	})
	t.Run("UnsupportedTypeFunc", func(t *testing.T) {
		v1 := reflect.ValueOf(&Struct{}).Elem().FieldByName("Any")
		v2 := reflect.ValueOf(&Struct{Func: func() {}}).Elem().FieldByName("Func")
		ExpectPanic(
			t,
			func() { DeepCopy(v1, v2) },
			ErrorContains("func type cannot be copied"),
		)
	})
	t.Run("Hack", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			{
				v1 := &struct{ str string }{str: "Any"}
				v2 := Clone(v1)
				Expect(t, v1, Equal(v2))
			}

			v1 := &Struct{Any: &struct{ str string }{str: "Any"}}
			v2 := Clone(v1)
			Expect(t, reflect.DeepEqual(v1, v2), BeTrue())
		})
		t.Run("Failed", func(t *testing.T) {
			v1 := &Struct{Any: struct{ str string }{str: "Any"}}
			ExpectPanic(
				t,
				func() { _ = Clone(v1) },
				ErrorContains("unexported"),
				ErrorContains("cannot be copied"),
			)
		})
	})

	src := []any{
		1,
		ptrx.Ptr(100),
		[1]any{1},
		struct{ A string }{A: "100"},
		nil,
		map[int]int{1: 1},
	}
	dst := Clone(src)
	Expect(t, reflect.DeepEqual(src, dst), BeTrue())
}

func TestHackFieldByName(t *testing.T) {
	v := &struct {
		a string
		b struct {
			b0 string
		}
	}{"a", struct {
		b0 string
	}{"b0"}}

	rv := HackFieldByName(v, "a")
	Expect(t, v.a, Equal("a"))
	rv.Set(reflect.ValueOf("abc"))
	Expect(t, v.a, Equal("abc"))

	rv = HackField(v, 0)
	rv.Set(reflect.ValueOf("bbb"))
	Expect(t, v.a, Equal("bbb"))

	rv = HackFieldByName(reflect.ValueOf(v), "a")
	rv.Set(reflect.ValueOf("aaa"))
	Expect(t, v.a, Equal("aaa"))

	rv = HackFieldByName(&v.b, "b0")
	rv.Set(reflect.ValueOf("def"))
	Expect(t, v.b.b0, Equal("def"))

	t.Run("NotStructValue", func(t *testing.T) {
		ExpectPanic(
			t,
			func() { HackFieldByName(1, "any") },
			ErrorContains("not a struct value"),
		)

	})

	t.Run("FieldNotFound", func(t *testing.T) {
		ExpectPanic(
			t,
			func() { HackFieldByName(v, "any") },
			ErrorContains("not found"),
		)
	})
	t.Run("CannotAddr", func(t *testing.T) {
		ExpectPanic(
			t,
			func() { HackFieldByName(*v, "a") },
			ErrorContains("cannot addr"),
		)
	})
}
