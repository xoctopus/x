package typex_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestTypes_Composites(t *testing.T) {
	cases := []struct {
		name string
		c    *CaseAssertion
	}{
		{
			"EmptyArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmptyArray",
				String:        "github.com/xoctopus/x/typex/testdata.EmptyArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "interface {}",
				Len:           0,
			},
		},
		{
			"Array",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Array",
				String:        "github.com/xoctopus/x/typex/testdata.Array",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "string",
				Len:           1,
			},
		},
		{
			"AnyArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "AnyArray",
				String:        "github.com/xoctopus/x/typex/testdata.AnyArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "interface {}",
				Len:           2,
			},
		},
		{
			"StringArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "StringArray",
				String:        "github.com/xoctopus/x/typex/testdata.StringArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.String",
				Len:           4,
			},
		},
		{
			"StringPtrArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "StringPtrArray",
				String:        "github.com/xoctopus/x/typex/testdata.StringPtrArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "*github.com/xoctopus/x/typex/testdata.String",
				Len:           8,
			},
		},
		{
			"IntPtrDefArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "IntPtrDefArray",
				String:        "github.com/xoctopus/x/typex/testdata.IntPtrDefArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.IntPtrDef",
				Len:           16,
			},
		},
		{
			"SizedArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "SizedArray",
				String:        "github.com/xoctopus/x/typex/testdata.SizedArray",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Boolean",
				Len:           testdata.SIZE,
			},
		},
		{
			"Map",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Map",
				String:        "github.com/xoctopus/x/typex/testdata.Map",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "string",
				Elem:          "string",
			},
		},
		{
			"StringIntMap",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "StringIntMap",
				String:        "github.com/xoctopus/x/typex/testdata.StringIntMap",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "github.com/xoctopus/x/typex/testdata.String",
				Elem:          "github.com/xoctopus/x/typex/testdata.Int",
			},
		},
		{
			"Slice",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Slice",
				String:        "github.com/xoctopus/x/typex/testdata.Slice",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "string",
			},
		},
		{
			"ErrorSlice",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "ErrorSlice",
				String:        "github.com/xoctopus/x/typex/testdata.ErrorSlice",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Error",
			},
		},
		{
			"TypedArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "TypedArray[int]",
				String:        "github.com/xoctopus/x/typex/testdata.TypedArray[int]",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "int",
				Len:           1,
			},
		},
		{
			"IntegerArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "IntegerArray[int8]",
				String:        "github.com/xoctopus/x/typex/testdata.IntegerArray[int8]",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "int8",
				Len:           10,
			},
		},
		{
			"TypedSizedArray",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "TypedSizedArray[float32]",
				String:        "github.com/xoctopus/x/typex/testdata.TypedSizedArray[float32]",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "float32",
				Len:           testdata.SIZE,
			},
		},
		{
			"TypedSlice",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "TypedSlice[net.Addr]",
				String:        "github.com/xoctopus/x/typex/testdata.TypedSlice[net.Addr]",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "net.Addr",
			},
		},
		{
			"TypedMap",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "TypedMap[string,int]",
				String:        "github.com/xoctopus/x/typex/testdata.TypedMap[string,int]",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "string",
				Elem:          "int",
			},
		},
		{
			"TypedMap2",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "TypedMap[github.com/xoctopus/x/typex/testdata.String,int]",
				String:        "github.com/xoctopus/x/typex/testdata.TypedMap[github.com/xoctopus/x/typex/testdata.String,int]",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "github.com/xoctopus/x/typex/testdata.String",
				Elem:          "int",
			},
		},
		{
			"StructSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]github.com/xoctopus/x/typex/testdata.SimpleStruct",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.SimpleStruct",
			},
		},
		{
			"StructPtrSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]*github.com/xoctopus/x/typex/testdata.SimpleStruct",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "*github.com/xoctopus/x/typex/testdata.SimpleStruct",
			},
		},
	}

	rtype := reflect.TypeOf(testdata.CompositeBasics{})
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tt := rtype.Field(i).Type
			NewWithT(t).Expect(rtype.Field(i).Name).To(Equal(c.name))

			rt := NewRType(tt)
			gt := NewGType(tt)

			NewWithT(t).Expect(rt.Unwrap()).To(Equal(tt))
			NewWithT(t).Expect(gt.Unwrap()).To(Equal(internal.NewTypesTypeFromReflectType(tt)))

			c.c.Check(t, rt, gt)
		})
	}
}
