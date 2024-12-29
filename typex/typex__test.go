package typex_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex"
)

func typename(t Type) string {
	if t == nil {
		return "nil"
	}

	buf := bytes.NewBuffer(nil)
	if t.Name() == "" {
		for t.Kind() == reflect.Pointer {
			buf.WriteByte('*')
			t = t.Elem()
		}
	}

	if name := t.Name(); name != "" {
		if path := t.PkgPath(); path != "" {
			buf.WriteString(path)
			buf.WriteRune('.')
		}
		buf.WriteString(name)
		return buf.String()
	}

	buf.WriteString(t.String())
	return buf.String()
}

func NewBases() *Bases {
	b := &Bases{}
	t := reflect.TypeOf(Bases{})
	for i := 0; i < t.NumField()-1; i++ {
		b.fields = append(b.fields, t.Field(i))
	}
	return b
}

var bases = NewBases()

type Bases struct {
	FmtStringer    fmt.Stringer
	Bytes          interface{ Bytes() []byte }
	EmptyInterface any
	Struct         struct{ some any }
	EmptyStruct    struct{}
	Error          error
	fields         []reflect.StructField
}

func (b *Bases) CheckImplements(t *testing.T, rt *RType, gt *GType, assertions []bool) {
	NewWithT(t).Expect(len(assertions)).To(Equal(len(b.fields)))
	for i, f := range b.fields {
		brt := NewRType(f.Type)
		bgt := NewGType(f.Type)
		t.Run(f.Name+"#"+strconv.Itoa(i), func(t *testing.T) {
			NewWithT(t).Expect(rt.Implements(brt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(rt.Implements(bgt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(gt.Implements(brt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(gt.Implements(bgt)).To(Equal(assertions[i]))
		})
		if f.Name == "EmptyInterface" {
			NewWithT(t).Expect(assertions[i]).To(BeTrue())
		}
		if f.Name == "SomeStruct" {
			NewWithT(t).Expect(assertions[i]).To(BeFalse())
		}
	}
}

func (b *Bases) CheckAssignableTo(t *testing.T, rt *RType, gt *GType, assertions []bool) {
	NewWithT(t).Expect(len(assertions)).To(Equal(len(b.fields)))
	for i, f := range b.fields {
		brt := NewRType(f.Type)
		bgt := NewGType(f.Type)
		t.Run(f.Name+"#"+strconv.Itoa(i), func(t *testing.T) {
			NewWithT(t).Expect(rt.AssignableTo(brt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(rt.AssignableTo(bgt)).To(BeFalse())
			NewWithT(t).Expect(gt.AssignableTo(bgt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(gt.AssignableTo(brt)).To(BeFalse())
		})
		if f.Name == "EmptyInterface" {
			NewWithT(t).Expect(assertions[i]).To(BeTrue())
		}
		if f.Name == "SomeStruct" {
			NewWithT(t).Expect(assertions[i]).To(BeFalse())
		}
	}
}

func (b *Bases) CheckConvertibleTo(t *testing.T, rt *RType, gt *GType, assertions []bool) {
	NewWithT(t).Expect(len(assertions)).To(Equal(len(b.fields)))
	for i, f := range b.fields {
		brt := NewRType(f.Type)
		bgt := NewGType(f.Type)
		t.Run(f.Name+"#"+strconv.Itoa(i), func(t *testing.T) {
			NewWithT(t).Expect(rt.ConvertibleTo(brt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(rt.ConvertibleTo(bgt)).To(BeFalse())
			NewWithT(t).Expect(gt.ConvertibleTo(bgt)).To(Equal(assertions[i]))
			NewWithT(t).Expect(gt.ConvertibleTo(brt)).To(BeFalse())
		})
		if f.Name == "EmptyInterface" {
			NewWithT(t).Expect(assertions[i]).To(BeTrue())
		}
		if f.Name == "SomeStruct" {
			NewWithT(t).Expect(assertions[i]).To(BeFalse())
		}
	}
}

type FieldAssertion struct {
	PkgPath   string
	Name      string
	Type      string
	Tag       string
	Anonymous bool
}

func (c *FieldAssertion) Check(t *testing.T, f StructField) {
	NewWithT(t).Expect(f.Name()).To(Equal(c.Name))
	NewWithT(t).Expect(f.PkgPath()).To(Equal(c.PkgPath))
	NewWithT(t).Expect(typename(f.Type())).To(Equal(c.Type))
	NewWithT(t).Expect(string(f.Tag())).To(Equal(c.Tag))
	NewWithT(t).Expect(f.Anonymous()).To(Equal(c.Anonymous))
}

type MethodAssertion struct {
	PkgPath string
	Name    string
	Type    string
}

func (c *MethodAssertion) Check(t *testing.T, m Method) {
	NewWithT(t).Expect(m.PkgPath()).To(Equal(c.PkgPath))
	NewWithT(t).Expect(m.Name()).To(Equal(c.Name))
	NewWithT(t).Expect(typename(m.Type())).To(Equal(c.Type))
}

type CaseAssertion struct {
	PkgPath string
	Name    string
	String  string
	Kind    reflect.Kind

	Implements    []bool
	AssignableTo  []bool
	ConvertibleTo []bool
	Comparable    bool

	Key  string
	Elem string
	Len  int

	NumField int
	Fields   []FieldAssertion

	NumMethod int
	Methods   []MethodAssertion

	IsVariadic bool
	NumIn      int
	Ins        []string
	NumOut     int
	Outs       []string
}

func (c *CaseAssertion) Check(t *testing.T, rt *RType, gt *GType) {
	tt := rt.Unwrap().(reflect.Type)
	t.Run("PkgPath", func(t *testing.T) {
		NewWithT(t).Expect(rt.PkgPath()).To(Equal(c.PkgPath))
		NewWithT(t).Expect(gt.PkgPath()).To(Equal(c.PkgPath))
	})
	t.Run("Name", func(t *testing.T) {
		NewWithT(t).Expect(rt.Name()).To(Equal(c.Name))
		NewWithT(t).Expect(gt.Name()).To(Equal(c.Name))
	})
	t.Run("String", func(t *testing.T) {
		NewWithT(t).Expect(rt.String()).To(Equal(c.String))
		NewWithT(t).Expect(gt.String()).To(Equal(c.String))
	})
	t.Run("Kind", func(t *testing.T) {
		NewWithT(t).Expect(rt.Kind()).To(Equal(c.Kind))
		NewWithT(t).Expect(gt.Kind()).To(Equal(c.Kind))
	})
	t.Run("Implements", func(t *testing.T) {
		bases.CheckImplements(t, rt, gt, c.Implements)
	})
	t.Run("AssignableTo", func(t *testing.T) {
		bases.CheckAssignableTo(t, rt, gt, c.AssignableTo)
	})
	t.Run("ConvertibleTo", func(t *testing.T) {
		bases.CheckConvertibleTo(t, rt, gt, c.ConvertibleTo)
	})
	t.Run("Comparable", func(t *testing.T) {
		NewWithT(t).Expect(rt.Comparable()).To(Equal(c.Comparable))
		NewWithT(t).Expect(gt.Comparable()).To(Equal(c.Comparable))
	})
	t.Run("Key", func(t *testing.T) {
		NewWithT(t).Expect(typename(rt.Key())).To(Equal(c.Key))
		NewWithT(t).Expect(typename(gt.Key())).To(Equal(c.Key))
	})
	t.Run("Elem", func(t *testing.T) {
		NewWithT(t).Expect(typename(rt.Elem())).To(Equal(c.Elem))
		NewWithT(t).Expect(typename(gt.Elem())).To(Equal(c.Elem))
	})
	t.Run("Len", func(t *testing.T) {
		NewWithT(t).Expect(rt.Len()).To(Equal(c.Len))
		NewWithT(t).Expect(gt.Len()).To(Equal(c.Len))
	})
	t.Run("Fields", func(t *testing.T) {
		if tt.Kind() == reflect.Struct {
			NewWithT(t).Expect(tt.NumField()).To(Equal(c.NumField))
		} else {
			NewWithT(t).Expect(0).To(Equal(c.NumField))
		}
		NewWithT(t).Expect(rt.NumField()).To(Equal(c.NumField))
		NewWithT(t).Expect(gt.NumField()).To(Equal(c.NumField))
		NewWithT(t).Expect(c.NumField).To(Equal(len(c.Fields)))
		for i := 0; i < c.NumField; i++ {
			t.Run(rt.Field(i).Name(), func(t *testing.T) {
				c.Fields[i].Check(t, rt.Field(i))
				c.Fields[i].Check(t, gt.Field(i))
			})
		}
		t.Run("FieldByName", func(t *testing.T) {
			for i, name := range []string{"unexported", "String", "Name"} {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					expect := false
					if tt.Kind() == reflect.Struct {
						_, expect = tt.FieldByName(name)
					}
					_, exists1 := rt.FieldByName(name)
					_, exists2 := gt.FieldByName(name)
					NewWithT(t).Expect(exists1).To(Equal(expect))
					NewWithT(t).Expect(exists2).To(Equal(expect))
				})
			}
		})
		t.Run("FieldByNameFunc", func(t *testing.T) {
			for i, matcher := range []func(string) bool{
				func(string) bool { return false },
				func(string) bool { return true },
				func(x string) bool { return x == "Name" },
				func(x string) bool { return x == "unexported" },
				func(x string) bool { return x == "String" },
			} {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					expect := false
					if tt.Kind() == reflect.Struct {
						_, expect = tt.FieldByNameFunc(matcher)
					}
					_, exists1 := rt.FieldByNameFunc(matcher)
					_, exists2 := gt.FieldByNameFunc(matcher)
					NewWithT(t).Expect(exists1).To(Equal(expect))
					NewWithT(t).Expect(exists2).To(Equal(expect))
				})
			}
		})
	})
	t.Run("Methods", func(t *testing.T) {
		NewWithT(t).Expect(tt.NumMethod()).To(Equal(c.NumMethod))
		NewWithT(t).Expect(rt.NumMethod()).To(Equal(c.NumMethod))
		NewWithT(t).Expect(gt.NumMethod()).To(Equal(c.NumMethod))
		NewWithT(t).Expect(c.NumMethod).To(Equal(len(c.Methods)))
		for i := 0; i < c.NumMethod; i++ {
			t.Run(rt.Method(i).Name(), func(t *testing.T) {
				c.Methods[i].Check(t, rt.Method(i))
				c.Methods[i].Check(t, gt.Method(i))
			})
		}
		t.Run("MethodByName", func(t *testing.T) {
			for _, name := range []string{"String", "_", "Name"} {
				_, expect := tt.MethodByName(name)
				_, exists0 := rt.MethodByName(name)
				_, exists1 := gt.MethodByName(name)
				NewWithT(t).Expect(exists0).To(Equal(expect))
				NewWithT(t).Expect(exists1).To(Equal(expect))
			}
		})
	})
	t.Run("IsVariadic", func(t *testing.T) {
		NewWithT(t).Expect(rt.IsVariadic()).To(Equal(c.IsVariadic))
		NewWithT(t).Expect(gt.IsVariadic()).To(Equal(c.IsVariadic))
	})
	t.Run("Ins", func(t *testing.T) {
		NewWithT(t).Expect(rt.NumIn()).To(Equal(c.NumIn))
		NewWithT(t).Expect(gt.NumIn()).To(Equal(c.NumIn))
		NewWithT(t).Expect(c.NumIn).To(Equal(len(c.Ins)))
		for i := 0; i < rt.NumIn(); i++ {
			NewWithT(t).Expect(typename(rt.In(i))).To(Equal(c.Ins[i]))
			NewWithT(t).Expect(typename(gt.In(i))).To(Equal(c.Ins[i]))
		}
	})
	t.Run("Outs", func(t *testing.T) {
		NewWithT(t).Expect(rt.NumOut()).To(Equal(c.NumOut))
		NewWithT(t).Expect(gt.NumOut()).To(Equal(c.NumOut))
		NewWithT(t).Expect(c.NumOut).To(Equal(len(c.Outs)))
		for i := 0; i < rt.NumOut(); i++ {
			NewWithT(t).Expect(typename(rt.Out(i))).To(Equal(c.Outs[i]))
			NewWithT(t).Expect(typename(gt.Out(i))).To(Equal(c.Outs[i]))
		}
	})
}

/*
type Types struct {
	Array                    [3]int
	ArrayPtr                 *[3]int
	NamedArray               testdata.Array
	NamedArrayPtr            *testdata.Array
	Map                      map[string]int
	MapPtr                   *map[string]int
	NamedMap                 testdata.Map
	NamedMapPtr              *testdata.Map
	Slice                    []string
	SlicePtr                 *[]string
	NamedSlice               testdata.Slice
	NamedSlicePtr            *testdata.Slice
	NamedChan                testdata.Chan
	Func                     func(x, y string) int
	NamedFunc                testdata.Func
	Interface                any
	NamedInterface           testdata.Interface
	NamedInterfaceImpl       testdata.InterfaceImpl
	NamedInterfaceImplPtr    *testdata.InterfaceImpl
	NamedInterfacePtrImpl    testdata.InterfacePtrImpl
	NamedInterfacePtrImplPtr *testdata.InterfacePtrImpl
	Struct                   testdata.Struct
	StructPtr                *testdata.Struct
	Enum                     testdata.Enum
	EnumPtr                  *testdata.Enum
	MixedInterface           testdata.MixedInterface
	MixedInterfaceImpl       testdata.MixedInterfaceImpl
	MixedInterfaceImplPtr    *testdata.MixedInterfaceImpl
	AnySliceInt              testdata.AnySlice[int]
	AnySliceEnum             testdata.AnySlice[testdata.Enum]
	AnyArrayString           testdata.AnyArray[string]
	AnyArrayInterface        testdata.AnyArray[testdata.Interface]
	AnyMapIntString          testdata.AnyMap[int, string]
	EnumMap                  testdata.EnumMap
	AnyStructAny             testdata.AnyStruct[any]
	AnyStructAnyPtr          *testdata.AnyStruct[any]
	AnyStructNamedInterface  testdata.AnyStruct[testdata.Interface]
	AnyComposeStringInt      testdata.AnyCompose[string, int]
	AnyComposeStringerFloat  testdata.AnyCompose[fmt.Stringer, float32]
	AnonymousStruct          struct{ testdata.AnyStruct[any] }
	NamedAnyInterface        testdata.AnyInterface[int, fmt.Stringer]
	UnameStructPtrVar        *struct {
		testdata.AnyCompose[testdata.AnyStruct[fmt.Stringer], testdata.Int]
	}
	Error error
}
*/
