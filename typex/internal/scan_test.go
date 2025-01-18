package internal_test

import (
	"go/types"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestScanMethodsAndFields(t *testing.T) {
	cases := testdata.StructureCases
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tt := NewTypesType(c.Type)
			methods := ScanMethods(tt)
			NewWithT(t).Expect(methods.NumMethod()).To(Equal(c.Type.NumMethod()))
			for i := range c.Type.NumMethod() {
				m1 := methods.Method(i)
				m2 := c.Type.Method(i)
				NewWithT(t).Expect(m1.Name()).To(Equal(m2.Name))
			}
			NewWithT(t).Expect(methods.Method(methods.NumMethod() + 1)).To(BeNil())

			for _, name := range []string{"String", "Name", "name", "unexported"} {
				t.Run(name, func(t *testing.T) {
					m1, ok := c.Type.MethodByName(name)
					m2 := methods.MethodByName(name)
					if ok {
						NewWithT(t).Expect(m2.Name()).To(Equal(m1.Name))
					} else {
						NewWithT(t).Expect(m2).To(BeNil())
					}
					if c.Type.Kind() == reflect.Struct {
						f1, ok := c.Type.FieldByName(name)
						f2 := FieldByName(tt, func(v string) bool { return v == name })
						if ok {
							NewWithT(t).Expect(f2.Name()).To(Equal(f1.Name))
						} else {
							NewWithT(t).Expect(f2).To(BeNil())
						}
					}
				})
			}
		})
	}

	t.Run("PotentialMatchFunc", func(t *testing.T) {
		tt := NewPackage("github.com/xoctopus/x/typex/testdata").Scope().Lookup("SimpleStruct").Type()
		field := FieldByName(tt, func(s string) bool { return true })
		NewWithT(t).Expect(field).To(BeNil())
	})

	t.Run("PointerEntry0", func(t *testing.T) {
		tt := NewPackage("github.com/xoctopus/x/typex/testdata").Scope().Lookup("SimpleStruct").Type()
		tt = types.NewPointer(tt)
		field := FieldByName(tt, func(s string) bool { return s == "A" })
		NewWithT(t).Expect(field).To(BeNil())
	})

	t.Run("AliasType", func(t *testing.T) {
		tt := NewPackage("github.com/xoctopus/x/typex/testdata").Scope().Lookup("TypedSliceAlias").Type()
		methods := ScanMethods(tt)
		NewWithT(t).Expect(methods.NumMethod()).To(Equal(1))
		NewWithT(t).Expect(methods.Method(0).Name()).To(Equal("Len"))
		fields := InspectFields(tt, &NamedBacktrace{})
		NewWithT(t).Expect(len(fields)).To(Equal(0))
	})
}
