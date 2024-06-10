package typex_test

import (
	"bytes"
	"encoding"
	"go/types"
	"io"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/testdata"
)

func TestType(t *testing.T) {
	cases := []struct {
		name     string
		v        any
		typename string
	}{
		{"GenericStruct", testdata.AnyStruct[string]{Name: "any"}, "github.com/xoctopus/x/typex/testdata.AnyStruct[string]"},
		{"GenericSlice", testdata.AnySlice[string]{}, "github.com/xoctopus/x/typex/testdata.AnySlice[string]"},
		{"GenericMap", testdata.AnyMap[int, string]{}, "github.com/xoctopus/x/typex/testdata.AnyMap[int,string]"},
		{"GenericArray", testdata.AnyArray[string]{}, "github.com/xoctopus/x/typex/testdata.AnyArray[string]"},
		{"InstancedGeneric", testdata.EnumMap{}, "github.com/xoctopus/x/typex/testdata.AnyMap[github.com/xoctopus/x/typex/testdata.Enum,interface {}]"},
		{"ComposeStruct", testdata.Compose{}, "github.com/xoctopus/x/typex/testdata.Compose"},
		{
			"Func",
			func() *testdata.Enum {
				v := testdata.ENUM__ONE
				return &v
			},
			"func() *github.com/xoctopus/x/typex/testdata.Enum",
		},
		{"Enum", testdata.ENUM__ONE, "github.com/xoctopus/x/typex/testdata.Enum"},
		{"RtEncodingTextMarshaler", reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem(), "encoding.TextMarshaler"},
		{"RtMixedInterface", reflect.TypeOf((*testdata.MixedInterface)(nil)).Elem(), "github.com/xoctopus/x/typex/testdata.MixedInterface"},
		{"UnsafePointer", unsafe.Pointer(t), "unsafe.Pointer"},
		{"NamedChan", make(testdata.Chan), "github.com/xoctopus/x/typex/testdata.Chan"},
		{"Chan", make(chan int, 100), "chan int"},
		{"GlobalFunc", testdata.F, "func()"},
		{"NamedFunc", testdata.Func(func(string, string) bool { return true }), "github.com/xoctopus/x/typex/testdata.Func"},
		{"NamedStringArray", testdata.Array{}, "github.com/xoctopus/x/typex/testdata.Array"},
		{"StringArray", [1]string{}, "[1]string"},
		{"NamedStringSlice", testdata.Slice{}, "github.com/xoctopus/x/typex/testdata.Slice"},
		{"StringSlice", []string{}, "[]string"},
		{"NamedMap", testdata.Map{}, "github.com/xoctopus/x/typex/testdata.Map"},
		{"Map", map[string]string{}, "map[string]string"},
		{"Struct", testdata.Struct{}, "github.com/xoctopus/x/typex/testdata.Struct"},
		{"StructPtr", &testdata.Struct{}, "*github.com/xoctopus/x/typex/testdata.Struct"},
		{
			"UnameStruct",
			struct {
				A string         `json:"a"`
				B *testdata.Part `json:"part"`
			}{},
			"struct { A string `json:\"a\"`; B *github.com/xoctopus/x/typex/testdata.Part `json:\"part\"` }",
		},
		{"EmptyUnnameStruct", struct{}{}, "struct {}"},
		{
			"UnameComposedInterfaceSlice",
			[]interface {
				io.Reader
				io.Closer
				shouldConstrained()
			}{},
			"[]interface { Close() error; Read([]uint8) (int, error); .shouldConstrained() }",
		},
		{"FuncWithInAndOut", func(string, string) bool { return true }, "func(string, string) bool"},
		{"NamedString", testdata.String(""), "github.com/xoctopus/x/typex/testdata.String"},
		{"String", "", "string"},
		{"NamedBoolean", testdata.Boolean(true), "github.com/xoctopus/x/typex/testdata.Boolean"},
		{"Boolean", false, "bool"},
		{"NamedInt", testdata.Int(0), "github.com/xoctopus/x/typex/testdata.Int"},
		{"NamedIntPtr", ptrx.Ptr(testdata.Int(0)), "*github.com/xoctopus/x/typex/testdata.Int"},
		{"Int", int(0), "int"},
		{"IntPtr", ptrx.Ptr(int(0)), "*int"},
		{"NamedInt8", testdata.Int8(0), "github.com/xoctopus/x/typex/testdata.Int8"},
		{"NamedInt16", testdata.Int16(0), "github.com/xoctopus/x/typex/testdata.Int16"},
		{"NamedInt32", testdata.Int32(0), "github.com/xoctopus/x/typex/testdata.Int32"},
		{"NamedInt64", testdata.Int64(0), "github.com/xoctopus/x/typex/testdata.Int64"},
		{"NamedUint", testdata.Uint(0), "github.com/xoctopus/x/typex/testdata.Uint"},
		{"NamedUint8", testdata.Uint8(0), "github.com/xoctopus/x/typex/testdata.Uint8"},
		{"NamedUint16", testdata.Uint16(0), "github.com/xoctopus/x/typex/testdata.Uint16"},
		{"NamedUint32", testdata.Uint32(0), "github.com/xoctopus/x/typex/testdata.Uint32"},
		{"NamedUint64", testdata.Uint64(0), "github.com/xoctopus/x/typex/testdata.Uint64"},
		{"NamedUintptr", testdata.Uintptr(0), "github.com/xoctopus/x/typex/testdata.Uintptr"},
		{"NamedFloat32", testdata.Float32(0), "github.com/xoctopus/x/typex/testdata.Float32"},
		{"NamedFloat64", testdata.Float64(0), "github.com/xoctopus/x/typex/testdata.Float64"},
		{"NamedComplex64", testdata.Complex64(0), "github.com/xoctopus/x/typex/testdata.Complex64"},
		{"NamedComplex128", testdata.Complex128(0), "github.com/xoctopus/x/typex/testdata.Complex128"},
		{"Nil", nil, "nil"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rt, ok := c.v.(reflect.Type)
			if !ok {
				rt = reflect.TypeOf(c.v)
			}
			gt := NewGoTypeFromReflectType(rt)

			rtt := NewReflectType(rt)
			gtt := NewGoType(gt)

			tt := NewWithT(t)
			tt.Expect(Typename(rtt)).To(Equal(c.typename))
			if rtt.Unwrap() == nil {
				tt.Expect(rtt.Unwrap()).To(BeNil())
				tt.Expect(gtt.Unwrap()).To(BeNil())
				return
			}
			tt.Expect(rtt.Unwrap()).To(Equal(rt))
			tt.Expect(gtt.Unwrap()).To(Equal(gt))
			tt.Expect(rtt.String()).To(Equal(gtt.String()))
			tt.Expect(rtt.Kind().String()).To(Equal(gtt.Kind().String()))
			tt.Expect(rtt.Name()).To(Equal(gtt.Name()))
			tt.Expect(rtt.PkgPath()).To(Equal(gtt.PkgPath()))
			tt.Expect(rtt.Comparable()).To(Equal(gtt.Comparable()))
			tt.Expect(rtt.AssignableTo(NewReflectType(reflect.TypeOf("")))).
				To(Equal(gtt.AssignableTo(NewGoType(types.Typ[types.String]))))
			tt.Expect(rtt.AssignableTo(NewReflectType(reflect.TypeOf(0)))).
				To(Equal(gtt.AssignableTo(NewGoType(types.Typ[types.Int]))))
			tt.Expect(rtt.ConvertibleTo(NewReflectType(reflect.TypeOf("")))).
				To(Equal(gtt.ConvertibleTo(NewGoType(types.Typ[types.String]))))
			tt.Expect(rtt.ConvertibleTo(NewReflectType(reflect.TypeOf(0)))).
				To(Equal(gtt.ConvertibleTo(NewGoType(types.Typ[types.Int]))))
			tt.Expect(rtt.AssignableTo(gtt)).To(BeFalse())
			tt.Expect(gtt.AssignableTo(rtt)).To(BeFalse())
			tt.Expect(rtt.ConvertibleTo(gtt)).To(BeFalse())
			tt.Expect(gtt.ConvertibleTo(rtt)).To(BeFalse())
			if rtt.Kind() == reflect.Struct {
				tt.Expect(rtt.NumField()).To(Equal(gtt.NumField()))
			}
			if rtt.Kind() == reflect.Pointer && Deref(rtt).Kind() == reflect.Struct {
				tt.Expect(Deref(rtt).NumField()).To(Equal(gtt.NumField()))
			}
			tt.Expect(rtt.NumMethod()).To(Equal(gtt.NumMethod()))
			for i := 0; i < rtt.NumMethod(); i++ {
				rm := rtt.Method(i)
				gm, exists := gtt.MethodByName(rm.Name())
				tt.Expect(exists).To(BeTrue())
				tt.Expect(rm.Name()).To(Equal(gm.Name()))
				tt.Expect(rm.PkgPath()).To(Equal(gm.PkgPath()))
				tt.Expect(rm.Type().String()).To(Equal(gm.Type().String()))
			}

			_, rStringMethodExists := rtt.MethodByName("String")
			_, gStringMethodExists := gtt.MethodByName("String")
			tt.Expect(rStringMethodExists).To(Equal(gStringMethodExists))

			if rtt.Kind() == reflect.Array {
				tt.Expect(rtt.Len()).To(Equal(gtt.Len()))
				tt.Expect(Typename(rtt.Elem())).To(Equal(Typename(gtt.Elem())))
			}

			if rtt.Kind() == reflect.Map {
				tt.Expect(Typename(rtt.Key())).To(Equal(Typename(gtt.Key())))
				tt.Expect(Typename(rtt.Elem())).To(Equal(Typename(gtt.Elem())))
			}

			if rtt.Kind() == reflect.Slice {
				tt.Expect(Typename(rtt.Elem())).To(Equal(Typename(gtt.Elem())))
			}

			if rtt.Kind() == reflect.Struct {
				tt.Expect(rtt.NumField()).To(Equal(gtt.NumField()))
				for i := 0; i < rtt.NumField(); i++ {
					rsf := rtt.Field(i)
					gsf := gtt.Field(i)

					tt.Expect(rsf.Anonymous()).To(Equal(gsf.Anonymous()))
					tt.Expect(rsf.Tag()).To(Equal(gsf.Tag()))
					tt.Expect(rsf.Name()).To(Equal(gsf.Name()))
					tt.Expect(rsf.PkgPath()).To(Equal(gsf.PkgPath()))
					tt.Expect(Typename(rsf.Type())).To(Equal(Typename(gsf.Type())))
				}

				rsf, exists1 := rtt.FieldByName("A")
				gsf, exists2 := gtt.FieldByName("A")
				tt.Expect(exists1).To(Equal(exists2))
				if exists2 {
					tt.Expect(rsf.Anonymous()).To(Equal(gsf.Anonymous()))
					tt.Expect(rsf.Tag()).To(Equal(gsf.Tag()))
					tt.Expect(rsf.Name()).To(Equal(gsf.Name()))
					tt.Expect(rsf.PkgPath()).To(Equal(gsf.PkgPath()))
					tt.Expect(Typename(rsf.Type())).To(Equal(Typename(gsf.Type())))
				}

				_, exists1 = rtt.FieldByName("_")
				_, exists2 = gtt.FieldByName("_")
				tt.Expect(exists1).To(Equal(exists2))
				tt.Expect(exists1).To(BeFalse())

				match := func(s string) bool {
					return s == "A"
				}
				rsf, exists1 = rtt.FieldByNameFunc(match)
				gsf, exists2 = gtt.FieldByNameFunc(match)
				tt.Expect(exists1).To(Equal(exists2))
				if exists1 {
					tt.Expect(rsf.Anonymous()).To(Equal(gsf.Anonymous()))
					tt.Expect(rsf.Tag()).To(Equal(gsf.Tag()))
					tt.Expect(rsf.Name()).To(Equal(gsf.Name()))
					tt.Expect(rsf.PkgPath()).To(Equal(gsf.PkgPath()))
					tt.Expect(Typename(rsf.Type())).To(Equal(Typename(gsf.Type())))
				}

				match = func(s string) bool { return false }
				_, exists1 = rtt.FieldByNameFunc(match)
				_, exists2 = gtt.FieldByNameFunc(match)
				tt.Expect(exists1).To(Equal(exists2))
				tt.Expect(exists1).To(BeFalse())
			}

			if rtt.Kind() == reflect.Func {
				tt.Expect(rtt.NumIn()).To(Equal(gtt.NumIn()))
				tt.Expect(rtt.NumOut()).To(Equal(gtt.NumOut()))
				for i := 0; i < rtt.NumIn(); i++ {
					tt.Expect(rtt.In(i).String()).To(Equal(gtt.In(i).String()))
				}
				for i := 0; i < rtt.NumOut(); i++ {
					tt.Expect(rtt.Out(i).String()).To(Equal(gtt.Out(i).String()))
				}
			}

			tt.Expect(Deref(rtt).String()).To(Equal(Deref(gtt).String()))
		})
	}
}

func TestGoTypeExt(t *testing.T) {
	t.Run("BuildUnsafePointerGoType", func(t *testing.T) {
		tpe := types.NewTypeName(0, NewPackage("unsafe"), "Pointer", types.Typ[types.UnsafePointer])
		gt := types.NewNamed(tpe, tpe.Type().Underlying(), nil)
		NewWithT(t).Expect(NewGoType(gt).Kind()).To(Equal(reflect.UnsafePointer))
	})

	t.Run("InvalidGoType", func(t *testing.T) {
		gt := NewGoType(types.Typ[types.Invalid])
		NewWithT(t).Expect(gt.Key()).To(BeNil())
		NewWithT(t).Expect(gt.Elem()).To(BeNil())
		NewWithT(t).Expect(gt.Len()).To(BeZero())
		NewWithT(t).Expect(gt.Kind()).To(Equal(reflect.Invalid))
		NewWithT(t).Expect(gt.Field(1)).To(BeNil())
		NewWithT(t).Expect(gt.Method(1)).To(BeNil())
		NewWithT(t).Expect(gt.IsVariadic()).To(BeFalse())
		NewWithT(t).Expect(gt.NumIn()).To(BeZero())
		NewWithT(t).Expect(gt.In(0)).To(BeNil())
		NewWithT(t).Expect(gt.NumOut()).To(BeZero())
		NewWithT(t).Expect(gt.Out(0)).To(BeNil())

		gt = NewGoType(nil)
		NewWithT(t).Expect(gt.Kind()).To(Equal(reflect.Invalid))
	})

	t.Run("ConstrainedInterface", func(t *testing.T) {
		gt := NewGoType(types.NewInterfaceType([]*types.Func{
			types.NewFunc(
				0, nil, "shouldBeConstrained",
				types.NewSignatureType(
					nil, nil, nil,
					types.NewTuple(types.NewVar(0, nil, "v", types.NewInterfaceType(nil, nil))),
					nil, false,
				)),
		}, nil))
		NewWithT(t).Expect(gt.String()).To(Equal("interface { shouldBeConstrained(interface {}) }"))
		NewWithT(t).Expect(gt.Method(0).Name()).To(Equal("shouldBeConstrained"))
		NewWithT(t).Expect(gt.Method(0).Type().String()).To(Equal("func(interface {})"))
		NewWithT(t).Expect(gt.Method(0).PkgPath()).To(BeEmpty())
	})

	t.Run("UnexportedStructField", func(t *testing.T) {
		gt := NewGoType(NewGoTypeFromReflectType(reflect.TypeOf(testdata.Struct{})))
		f, exists := gt.FieldByName("a")
		NewWithT(t).Expect(exists).To(BeTrue())
		NewWithT(t).Expect(f.PkgPath()).To(Equal(gt.PkgPath()))

		gt = NewGoType(NewGoTypeFromReflectType(reflect.TypeOf(struct{ a string }{})))
		f, exists = gt.FieldByName("a")
		NewWithT(t).Expect(exists).To(BeTrue())
		NewWithT(t).Expect(f.PkgPath()).To(Equal("github.com/xoctopus/x/typex_test"))
	})

	t.Run("NewGoTypeFromTypeParam", func(t *testing.T) {
		gt := NewGoType(NewGoTypeFromReflectType(reflect.TypeOf(testdata.AnyArray[string]{})))
		gt = NewGoType(gt.Type.(*types.Named).TypeParams().At(0))
		// TODO if use the underlying type
		NewWithT(t).Expect(gt.Type.String()).To(Equal("string"))
		NewWithT(t).Expect(gt.String()).To(Equal("interface {}"))
	})

	t.Run("Implements", func(t *testing.T) {
		gt := NewGoType(NewGoTypeFromReflectType(reflect.TypeOf(ptrx.Ptr(testdata.Enum(1)))))
		t.Run("GoType", func(t *testing.T) {
			xgt := NewGoType(types.NewInterfaceType(
				[]*types.Func{
					types.NewFunc(0, nil, "String", types.NewSignatureType(
						nil, nil, nil, nil,
						types.NewTuple(
							types.NewVar(0, nil, "", types.Typ[types.String]),
						), false,
					)),
				},
				nil,
			))
			NewWithT(t).Expect(gt.Implements(xgt)).To(BeTrue())
		})
		t.Run("ReflectType", func(t *testing.T) {
			rt := reflect.TypeOf([]encoding.TextMarshaler{}).Elem()
			xrt := NewReflectType(rt)
			NewWithT(t).Expect(gt.Implements(xrt)).To(BeTrue())
		})
		t.Run("Invalid", func(t *testing.T) {
			t.Run("Nil", func(t *testing.T) {
				NewWithT(t).Expect(gt.Implements(nil)).To(BeFalse())
			})
			t.Run("NotInterface", func(t *testing.T) {
				NewWithT(t).Expect(gt.Implements(NewGoType(types.Typ[types.Int]))).To(BeFalse())
				NewWithT(t).Expect(gt.Implements(NewReflectType(reflect.TypeOf(1)))).To(BeFalse())
				NewWithT(t).Expect(gt.Implements(NewReflectType(reflect.TypeOf(bytes.Buffer{})))).To(BeFalse())
			})
		})
	})
}

func TestReflectTypeExt(t *testing.T) {
	rtt := NewReflectType(reflect.TypeOf(testdata.Enum(1)))
	xgt := NewGoType(types.NewInterfaceType(
		[]*types.Func{
			types.NewFunc(0, nil, "String", types.NewSignatureType(
				nil, nil, nil, nil,
				types.NewTuple(
					types.NewVar(0, nil, "", types.Typ[types.String]),
				), false,
			)),
		},
		nil,
	))
	xrt := NewReflectType(reflect.TypeOf([]encoding.TextMarshaler{}).Elem())
	NewWithT(t).Expect(rtt.Implements(xgt)).To(BeTrue())
	NewWithT(t).Expect(rtt.Implements(xrt)).To(BeTrue())
	NewWithT(t).Expect(rtt.Implements(nil)).To(BeFalse())

	rtt = NewReflectType(reflect.TypeOf(1))
	NewWithT(t).Expect(rtt.Implements(xgt)).To(BeFalse())
}

func TestComparable(t *testing.T) {
	gt := NewGoType(NewGoTypeFromReflectType(reflect.TypeOf([2]string{})))
	// true
	t.Log(gt.Comparable())

	rt := reflect.TypeOf(testdata.AnyArray[string]{})
	gt = NewGoType(NewGoTypeFromReflectType(rt))
	// true
	t.Log(rt.Comparable())
	// false
	t.Log(types.Comparable(gt.Type))
	// false
	t.Log(types.Comparable(gt.Type.Underlying()))
	// true
	t.Log(types.Comparable(ConstrainUnderlying(gt.Type.(*types.Named).TypeParams(), gt.Type.Underlying())))
}
