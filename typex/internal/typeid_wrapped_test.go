package internal_test

import (
	"fmt"
	"go/parser"
	"io"
	"net"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestTypeID(t *testing.T) {
	type (
		MyStringer interface{ String() testdata.String }
		MyError    = error
	)
	r := reflect.TypeOf(struct {
		Error         error
		AliasError    MyError
		String        string
		UnsafePointer unsafe.Pointer
		NamedError    testdata.Error
		Array         [3]float32
		Slice         []string
		Map0          map[string]int
		Map1          map[testdata.String]testdata.Int
		Map2          map[testdata.Serialized[string]]testdata.Serialized[[]byte]
		GenericType0  testdata.TypeParamPass[string, net.Addr]
		GenericType1  testdata.TypeParamPass[testdata.TypeParamPass[string, testdata.Serialized[[]byte]], testdata.Serialized[string]]
		GenericType2  testdata.BTreeNode[testdata.Serialized[string]]
		Chan0         chan int
		Chan1         <-chan int
		Chan2         chan<- int
		Pointer0      *int
		Pointer1      **int
		NamedPointer  *testdata.Int
		Func0         func()
		Func1         func(...testdata.Func)
		Func2         func(string, bool, ...any)
		Func3         func(func() bool) (bool, int)
		Structure0    struct{}
		Structure1    struct {
			A string `json:"a"`
		}
		Structure2 struct {
			A int `json:"a"`
			B testdata.AnyArray
		}
		Interface0 interface{}
		Interface1 interface {
			fmt.Stringer
			io.ReadCloser
		}
		Interface2 interface{ MyStringer }
	}{})
	cases := []struct {
		name     string
		wrapped  string
		typename string
	}{
		{"Error", `error`, "error"},
		{"AliasError", `error`, "error"},
		{"String", `string`, "string"},
		{"UnsafePointer", `unsafe.Pointer`, "unsafe.Pointer"},
		{
			"NamedError",
			`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Error`,
			"github.com/xoctopus/x/typex/testdata.Error",
		},
		{"Array", `[3]float32`, "[3]float32"},
		{"Slice", `[]string`, "[]string"},
		{"Map0", `map[string]int`, "map[string]int"},
		{
			"Map1",
			`map[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.String]xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Int`,
			`map[github.com/xoctopus/x/typex/testdata.String]github.com/xoctopus/x/typex/testdata.Int`,
		},
		{
			"Map2",
			`map[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]]xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]`,
			"map[github.com/xoctopus/x/typex/testdata.Serialized[string]]github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]",
		},
		{
			"GenericType0",
			`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypeParamPass[string,net.Addr]`,
			"github.com/xoctopus/x/typex/testdata.TypeParamPass[string,net.Addr]",
		},
		{
			"GenericType1",
			`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypeParamPass[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypeParamPass[string,xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]],xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]]`,
			"github.com/xoctopus/x/typex/testdata.TypeParamPass[github.com/xoctopus/x/typex/testdata.TypeParamPass[string,github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]],github.com/xoctopus/x/typex/testdata.Serialized[string]]",
		},
		{
			"GenericType2",
			"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.BTreeNode[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]]",
			"github.com/xoctopus/x/typex/testdata.BTreeNode[github.com/xoctopus/x/typex/testdata.Serialized[string]]",
		},
		{"Chan0", `chan int`, "chan int"},
		{"Chan1", `<-chan int`, "<-chan int"},
		{"Chan2", `chan<- int`, "chan<- int"},
		{"Pointer0", `*int`, "*int"},
		{"Pointer1", `**int`, "**int"},
		{"NamedPointer", "*xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Int", "*github.com/xoctopus/x/typex/testdata.Int"},
		{"Func0", `func()`, "func()"},
		{
			"Func1",
			`func(...xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Func)`,
			"func(...github.com/xoctopus/x/typex/testdata.Func)",
		},
		{"Func2", `func(string, bool, ...interface {})`, "func(string, bool, ...interface {})"},
		{"Func3", `func(func() bool) (bool, int)`, "func(func() bool) (bool, int)"},
		{"Structure0", `struct {}`, "struct {}"},
		{
			"Structure1",
			`struct { A string "json:\"a\"" }`,
			`struct { A string "json:\"a\"" }`,
		},
		{
			"Structure2",
			`struct { A int "json:\"a\""; B xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.AnyArray }`,
			`struct { A int "json:\"a\""; B github.com/xoctopus/x/typex/testdata.AnyArray }`,
		},
		{"Interface0", `interface {}`, "interface {}"},
		{
			"Interface1",
			`interface { Close() error; Read([]uint8) (int, error); String() string }`,
			`interface { Close() error; Read([]uint8) (int, error); String() string }`,
		},
		{
			"Interface2",
			`interface { String() xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.String }`,
			`interface { String() github.com/xoctopus/x/typex/testdata.String }`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, ok := r.FieldByName(c.name)
			NewWithT(t).Expect(ok).To(BeTrue())
			rt := f.Type

			typename := ReflectTypeID(rt)
			NewWithT(t).Expect(typename).To(Equal(c.typename))

			wrapped := ReflectTypeID(rt, true)
			NewWithT(t).Expect(wrapped).To(Equal(c.wrapped))

			wrap := NewWrappedType(wrapped)
			NewWithT(t).Expect(wrap.String()).To(Equal(c.wrapped))

			tt := NewTypesType(rt)
			NewWithT(t).Expect(TypesTypeID(tt)).To(Equal(c.typename))
			NewWithT(t).Expect(TypesTypeID(tt, true)).To(Equal(c.wrapped))
		})
	}

	t.Run("WrappedTypeID", func(t *testing.T) {
		t.Run("Ident", func(t *testing.T) {
			tt := NewWrappedType("Typename")
			NewWithT(t).Expect(tt.String()).To(Equal("Typename"))
		})
		t.Run("SelectorExpr", func(t *testing.T) {
			tt := NewWrappedType("PackagePath.Typename")
			NewWithT(t).Expect(tt.String()).To(Equal("PackagePath.Typename"))
		})
		t.Run("IndexExpr", func(t *testing.T) {
			t.Run("IndexListExpr", func(t *testing.T) {
				tt := NewWrappedType("PackagePath.Typename[TypeParam1]")
				NewWithT(t).Expect(tt.String()).To(Equal("PackagePath.Typename[TypeParam1]"))
			})
		})
		t.Run("IndexListExpr", func(t *testing.T) {
			tt := NewWrappedType("PackagePath.Typename[TypeParam1,TypeParams2]")
			NewWithT(t).Expect(tt.String()).To(Equal("PackagePath.Typename[TypeParam1,TypeParams2]"))
		})
		t.Run("InvalidExpr", func(t *testing.T) {
			_, err := parser.ParseExpr("method(a)") // CallExpr
			NewWithT(t).Expect(err).To(BeNil())
			defer func() {
				v := recover()
				NewWithT(t).Expect(v).NotTo(BeNil())
			}()
			NewWrappedType("int(a)")
		})
	})

	t.Run("InvalidReflectType", func(t *testing.T) {
		defer func() {
			v := recover()
			NewWithT(t).Expect(v).NotTo(BeNil())
		}()
		ReflectTypeID(&MockInvalidType{})
	})
}

type MockInvalidType struct {
	reflect.Type
}

func (*MockInvalidType) Name() string       { return "" }
func (*MockInvalidType) PkgPath() string    { return "" }
func (*MockInvalidType) Kind() reflect.Kind { return reflect.Invalid }
