package reflectx_test

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/reflectx"
)

func TestIndirect(t *testing.T) {
	v := &struct{ Any any }{Any: 100.1}

	type IntPtr *int

	cases := []*struct {
		input  any
		expect reflect.Value
	}{
		{1, reflect.ValueOf(1)},
		{ptrx.Ptr(1), reflect.ValueOf(1)},
		{(*int)(nil), InvalidValue},
		{ptrx.Ptr(ptrx.Ptr(0.2)), reflect.ValueOf(0.2)},
		{reflect.ValueOf(v).Elem().Field(0), reflect.ValueOf(100.1)},
		{ptrx.Ptr(IntPtr(ptrx.Ptr(1))), reflect.ValueOf(IntPtr(ptrx.Ptr(1)))},
	}

	for _, c := range cases {
		result := Indirect(c.input)
		if c.expect == InvalidValue {
			NewWithT(t).Expect(result).To(Equal(c.expect))
			continue
		}
		NewWithT(t).Expect(result.Type()).To(Equal(c.expect.Type()))
		NewWithT(t).Expect(result.Interface()).To(Equal(c.expect.Interface()))
	}
}

func TestIndirectNew(t *testing.T) {
	v1 := &struct{ Int int }{Int: 1}
	v2 := &struct{ Int ***int }{}
	v3 := &struct{ Any any }{Any: new(****int)}
	v4 := &struct{ v ***int }{}

	cases := []*struct {
		Input  any
		expect reflect.Value
	}{
		{1, reflect.ValueOf(1)},
		{ptrx.Ptr(1), reflect.ValueOf(1)},
		{nil, InvalidValue},
		{reflect.ValueOf(v1).Elem().Field(0), reflect.ValueOf(1)},
		{reflect.ValueOf(v2).Elem().Field(0), reflect.ValueOf(0)},
		{reflect.ValueOf(v3).Elem().Field(0), reflect.ValueOf(0)},
		{reflect.ValueOf(v4).Elem().Field(0), InvalidValue},
	}

	for _, c := range cases {
		result := IndirectNew(c.Input)
		if c.expect == InvalidValue {
			NewWithT(t).Expect(result).To(Equal(c.expect))
			continue
		}
		NewWithT(t).Expect(result.Type()).To(Equal(c.expect.Type()))
		NewWithT(t).Expect(result.Interface()).To(Equal(c.expect.Interface()))
	}

	v2f0rv := IndirectNew(reflect.ValueOf(v2).Elem().Field(0))
	v2f0rv.Set(reflect.ValueOf(100))
	NewWithT(t).Expect(***v2.Int).To(Equal(100))
}

func TestDeref(t *testing.T) {
	cases := []*struct {
		input  any
		expect reflect.Type
	}{
		{1, reflect.TypeOf(1)},
		{ptrx.Ptr(1), reflect.TypeOf(1)},
		{(*int)(nil), reflect.TypeOf(1)},
		{ptrx.Ptr(ptrx.Ptr(0.2)), reflect.TypeOf(float64(0))},
		{InvalidType, InvalidType},
		{any(1), reflect.TypeOf(int(0))},
	}

	for _, c := range cases {
		result := Deref(reflect.TypeOf(c.input))
		if c.expect == InvalidType {
			NewWithT(t).Expect(result).To(BeNil())
		} else {
			NewWithT(t).Expect(result).To(Equal(c.expect))
		}
	}
}

func TestNew(t *testing.T) {
	cases := []*struct {
		value any
		orig  reflect.Type
		elem  reflect.Type
	}{
		{1, reflect.TypeOf(int(0)), reflect.TypeOf(int(0))},
		{ptrx.Ptr(1), reflect.TypeOf(new(int)), reflect.TypeOf(int(0))},
		{ptrx.Ptr(ptrx.Ptr(0.2)), reflect.TypeOf(new(*float64)), reflect.TypeOf(new(float64))},
		{[]int{}, reflect.TypeOf([]int{}), reflect.TypeOf([]int{})},
	}

	for _, c := range cases {
		rt := reflect.TypeOf(c.value)

		rv := New(rt)
		rve := NewElem(rt)

		NewWithT(t).Expect(rv.Type()).To(Equal(c.orig))
		NewWithT(t).Expect(rve.Type()).To(Equal(c.elem))
	}
}

type Endpoint struct {
	Host string
	Port int16
}

func (v Endpoint) IsZero() bool {
	return v.Host == ""
}

type Address struct {
	Path string
}

func (v *Address) IsZero() bool {
	return v.Path == ""
}

type AlwaysZero struct{ V int }

func (AlwaysZero) IsZero() bool { return true }

type (
	Int   int
	Float float32
)

func TestIsZero(t *testing.T) {
	cases := []*struct {
		value any
		zero  bool
	}{
		{reflect.ValueOf(1), false},
		{InvalidValue, true},
		{InvalidType, true},
		{(*int)(nil), true},
		{ptrx.Ptr(1), false},
		{any(nil), true},
		{new(any), true},
		{[3]int{}, true},
		{[3]int{1}, false},
		{[0]int{}, true},
		{[]int{1}, false},
		{[]int{}, true},
		{map[string]struct{}{}, true},
		{map[string]struct{}{"1": {}}, false},
		{"", true},
		{"abc", false},
		{true, false},
		{false, true},
		{1, false},
		{0, true},
		{Endpoint{Port: 80}, true},
		{Endpoint{Host: "example.com"}, false},
		{Address{}, true},
		{&Address{Path: "/path/to/file"}, false},
		{AlwaysZero{}, true},
		{AlwaysZero{1}, true},
		{struct{}{}, true},
		{struct{ A int }{}, true},
	}

	for _, c := range cases {
		NewWithT(t).Expect(IsZero(c.value)).To(Equal(c.zero))
	}
}

func TestTypename(t *testing.T) {
	cases := []*struct {
		value    any
		typename string
	}{
		{1, "int"},
		{bytes.Buffer{}, "bytes.Buffer"},
		{bytes.NewBuffer(nil), "*bytes.Buffer"},
		{struct{}{}, "struct {}"},
		{struct{ int }{}, "struct { int }"},
	}

	for _, c := range cases {
		typename := Typename(reflect.TypeOf(c.value))
		NewWithT(t).Expect(typename).To(Equal(c.typename))
	}
}

type Byte byte

func TestIsBytes(t *testing.T) {
	cases := []*struct {
		value  any
		expect bool
	}{
		{1, false},
		{[]byte{}, true},
		{reflect.TypeOf([]byte{}), true},
		{[]Byte{}, false},
		{reflect.ValueOf([]Byte{}), false},
	}

	for _, c := range cases {
		NewWithT(t).Expect(IsBytes(c.value)).To(Equal(c.expect))
	}
}

func TestIsInteger(t *testing.T) {
	cases := []*struct {
		value  any
		expect bool
	}{
		{0.1, false},
		{1, true},
		{"1", false},
		{Int(0), true},
		{any(100), true},
		{reflect.TypeOf(1), true},
	}

	for _, c := range cases {
		NewWithT(t).Expect(IsInteger(c.value)).To(Equal(c.expect))
	}
}

func TestIsFloat(t *testing.T) {
	cases := []*struct {
		value  any
		expect bool
	}{
		{1, false},
		{0.1, true},
		{"1", false},
		{Int(0), false},
		{any(100), false},
		{any(100.0), true},
		{reflect.TypeOf(1.0), true},
	}

	for _, c := range cases {
		NewWithT(t).Expect(IsFloat(c.value)).To(Equal(c.expect))
	}
}

func TestIsNumeric(t *testing.T) {
	cases := []*struct {
		value  any
		expect bool
	}{
		{1, true},
		{0.1, true},
		{"1", false},
		{Int(0), true},
		{Float(0), true},
		{any(100), true},
		{any(100.0), true},
		{reflect.TypeOf(1), true},
		{reflect.TypeOf(1.0), true},
	}

	for _, c := range cases {
		NewWithT(t).Expect(IsNumeric(c.value)).To(Equal(c.expect))
	}
}

func TestCanElemType(t *testing.T) {
	cases := []struct {
		value  any
		expect bool
	}{
		{reflect.TypeOf([]int{}), true},
		{reflect.TypeOf(map[string]int{}), true},
		{reflect.TypeOf(1), false},
		{reflect.TypeOf(new(int)), true},
		{reflect.TypeOf(make(chan int)), true},
		{reflect.ValueOf([]int{}), true},
		{reflect.ValueOf(1), false},
		{reflect.Slice, true},
		{reflect.Int, false},
		{[]int{}, true},
		{100, false},
	}

	for _, c := range cases {
		NewWithT(t).Expect(CanElemType(c.value)).To(Equal(c.expect))
	}
}

func TestCanNilValue(t *testing.T) {
	cases := []struct {
		value  any
		expect bool
	}{
		{nil, false},
		{reflect.TypeOf(nil), false},
		{reflect.Value{}, false},
		{make(chan int), true},
		{(func())(nil), true},
		{interface{}(nil), false}, // type nil trap
		{map[string]int{}, true},
		{[]int{}, true},
		{new(int), true},
		{unsafe.Pointer(nil), true},
		{100, false},
		{struct{}{}, false},
	}

	for _, c := range cases {
		NewWithT(t).Expect(CanNilValue(c.value)).To(Equal(c.expect))
	}

	// type nil trap
	var v interface{} = []int(nil)
	NewWithT(t).Expect(CanNilValue(v)).To(Equal(true))
}
