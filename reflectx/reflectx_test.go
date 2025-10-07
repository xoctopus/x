package reflectx_test

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/reflectx"
	. "github.com/xoctopus/x/testx"
)

type (
	IntPtr    *int
	AnyHolder struct{ Any any }
)

type Interface interface{ Foo() }

type Impl struct{}

func (Impl) Foo() {}

func TestIndirect(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		input  any
		expect reflect.Value
		check  func(t *testing.T, result reflect.Value)
	}

	cases := []testCase{
		{
			name:   "NilInput",
			input:  nil,
			expect: InvalidValue,
			check: func(t *testing.T, result reflect.Value) {
				Expect(t, result.IsValid(), BeFalse())
			},
		}, {
			name:   "BasicInt",
			input:  1,
			expect: reflect.ValueOf(1),
		}, {
			name:   "PointerToInt",
			input:  ptrx.Ptr(1),
			expect: reflect.ValueOf(1),
		}, {
			name:   "NamedPointerNotDereferenced",
			input:  ptrx.Ptr(IntPtr(ptrx.Ptr(1))),
			expect: reflect.ValueOf(IntPtr(ptrx.Ptr(1))),
		}, {
			name:   "InterfaceWrappingPointer",
			input:  reflect.ValueOf(&AnyHolder{Any: 100.1}).Elem().Field(0),
			expect: reflect.ValueOf(100.1),
		}, {
			name:   "PointerToPointer",
			input:  ptrx.Ptr(ptrx.Ptr(0.2)),
			expect: reflect.ValueOf(0.2),
		}, {
			name:   "ReflectValueInput",
			input:  reflect.ValueOf(123),
			expect: reflect.ValueOf(123),
		}, {
			name:   "NestedInterfaces",
			input:  any(any(42)),
			expect: reflect.ValueOf(42),
		}, {
			name:   "NamedInterfaceNotDereferenced",
			input:  Interface(Impl{}),
			expect: reflect.ValueOf(Interface(Impl{})),
		}, {
			name:  "InterfaceNil",
			input: Interface(nil),
			check: func(t *testing.T, result reflect.Value) {
				Expect(t, result.IsValid(), BeFalse())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := Indirect(c.input)

			if c.check != nil {
				c.check(t, result)
				return
			}

			Expect(t, result.IsValid(), BeTrue())
			Expect(t, result.Type(), Equal(c.expect.Type()))
			Expect(t, result.Interface(), Equal(c.expect.Interface()))
		})
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
			Expect(t, result, Equal(c.expect))
			continue
		}
		Expect(t, result.Type(), Equal(c.expect.Type()))
		Expect(t, result.Interface(), Equal(c.expect.Interface()))
	}

	v2f0rv := IndirectNew(reflect.ValueOf(v2).Elem().Field(0))
	v2f0rv.Set(reflect.ValueOf(100))
	Expect(t, ***v2.Int, Equal(100))
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
			Expect(t, result, BeNil[reflect.Type]())
		} else {
			Expect(t, result, Equal(c.expect))
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

		Expect(t, rv.Type(), Equal(c.orig))
		Expect(t, rve.Type(), Equal(c.elem))
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
		Expect(t, IsZero(c.value), Equal(c.zero))
	}
}

type TypedMap[K comparable, V any] map[K]V

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
		{TypedMap[int, string]{}, "github.com/xoctopus/x/reflectx_test.TypedMap[int,string]"},
	}

	for _, c := range cases {
		typename := Typename(reflect.TypeOf(c.value))
		Expect(t, typename, Equal(c.typename))
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
		Expect(t, IsBytes(c.value), Equal(c.expect))
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
		{reflect.Int, true},
	}

	for _, c := range cases {
		Expect(t, IsInteger(c.value), Equal(c.expect))
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
		Expect(t, IsFloat(c.value), Equal(c.expect))
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
		Expect(t, IsNumeric(c.value), Equal(c.expect))
	}
}

func TestCanElemType(t *testing.T) {
	cases := []struct {
		typ    reflect.Type
		expect bool
	}{
		{reflect.TypeOf([]int{}), true},
		{reflect.TypeOf(map[string]int{}), true},
		{reflect.TypeOf(1), false},
		{reflect.TypeOf(new(int)), true},
		{reflect.TypeOf(make(chan int)), true},
		{reflect.ValueOf([]int{}).Type(), true},
		{reflect.ValueOf(1).Type(), false},
		{reflect.TypeOf(nil), false},
	}

	for _, c := range cases {
		Expect(t, CanElem(c.typ), Equal(c.expect))
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
		Expect(t, CanNilValue(c.value), Equal(c.expect))
	}

	// type nil trap
	var v interface{} = []int(nil)
	Expect(t, CanNilValue(v), BeTrue())
}
