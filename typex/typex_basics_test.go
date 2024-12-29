package typex_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestTypes_Basics(t *testing.T) {
	cases := []struct {
		name string
		c    *CaseAssertion
	}{
		{
			"String",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "string",
				String:        "string",
				Kind:          reflect.String,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"Boolean",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "bool",
				String:        "bool",
				Kind:          reflect.Bool,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"Int",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "int",
				String:        "int",
				Kind:          reflect.Int,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"UnsafePointer",
			&CaseAssertion{
				PkgPath:       "unsafe",
				Name:          "Pointer",
				String:        "unsafe.Pointer",
				Kind:          reflect.UnsafePointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"Error",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "error",
				String:        "error",
				Kind:          reflect.Interface,
				Implements:    []bool{false, false, true, false, false, true},
				AssignableTo:  []bool{false, false, true, false, false, true},
				ConvertibleTo: []bool{false, false, true, false, false, true},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						PkgPath: "",
						Name:    "Error",
						Type:    "func() string",
					},
				},
			},
		},
		{
			"Chan",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "chan string",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "string",
			},
		},
		{
			"SendOnlyChan",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "chan<- string",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "string",
			},
		},
		{
			"RecvOnlyChan",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "<-chan string",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "string",
			},
		},
		{
			"IntPtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*int",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "int",
			},
		},
		{
			"StringArray",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[3]string",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "string",
				Len:           3,
			},
		},
		{
			"IntSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]int32",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "int32",
			},
		},
		{
			"IntPtrSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]*int64",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "*int64",
			},
		},
		{
			"IntStringMap",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "map[int]string",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "int",
				Elem:          "string",
			},
		},
		{
			"IntSet",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "map[int]struct {}",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "int",
				Elem:          "struct {}",
			},
		},
		{
			"EmptyStruct",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "struct {}",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, true, false},
				ConvertibleTo: []bool{false, false, true, false, true, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"EmptyInterface",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "interface {}",
				Kind:          reflect.Interface,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"Func",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "func(int, string, ...interface {}) (float32, error)",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    true,
				NumIn:         3,
				Ins:           []string{"int", "string", "[]interface {}"},
				NumOut:        2,
				Outs:          []string{"float32", "error"},
			},
		},
		{
			"Curry",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "func(interface {}) func() string",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    false,
				NumIn:         1,
				Ins:           []string{"interface {}"},
				NumOut:        1,
				Outs:          []string{"func() string"},
			},
		},
		{
			"NamedString",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "String",
				String:        "github.com/xoctopus/x/typex/testdata.String",
				Kind:          reflect.String,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"NamedBoolean",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Boolean",
				String:        "github.com/xoctopus/x/typex/testdata.Boolean",
				Kind:          reflect.Bool,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"NamedInt",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Int",
				String:        "github.com/xoctopus/x/typex/testdata.Int",
				Kind:          reflect.Int,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"NamedUnsafePointer",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "UnsafePointer",
				String:        "github.com/xoctopus/x/typex/testdata.UnsafePointer",
				Kind:          reflect.UnsafePointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"NamedError",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Error",
				String:        "github.com/xoctopus/x/typex/testdata.Error",
				Kind:          reflect.Interface,
				Implements:    []bool{false, false, true, false, false, true},
				AssignableTo:  []bool{false, false, true, false, false, true},
				ConvertibleTo: []bool{false, false, true, false, false, true},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						PkgPath: "",
						Name:    "Error",
						Type:    "func() string",
					},
				},
			},
		},
		{
			"NamedChan",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Chan",
				String:        "github.com/xoctopus/x/typex/testdata.Chan",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Int",
			},
		},
		{
			"NamedSendOnlyChan",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "SendOnlyChan",
				String:        "github.com/xoctopus/x/typex/testdata.SendOnlyChan",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.String",
			},
		},
		{
			"NamedRecvOnlyChan",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "RecvOnlyChan",
				String:        "github.com/xoctopus/x/typex/testdata.RecvOnlyChan",
				Kind:          reflect.Chan,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Error",
			},
		},
		{
			"NamedIntPtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.Int",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Int",
			},
		},
		{
			"NamedIntPtrDef",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "IntPtrDef",
				String:        "github.com/xoctopus/x/typex/testdata.IntPtrDef",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "int",
			},
		},
		{
			"NamedStringArray",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[10]github.com/xoctopus/x/typex/testdata.String",
				Kind:          reflect.Array,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.String",
				Len:           10,
			},
		},
		{
			"NamedIntSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]github.com/xoctopus/x/typex/testdata.Int",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.Int",
			},
		},
		{
			"NamedIntPtrSlice",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "[]*github.com/xoctopus/x/typex/testdata.Int",
				Kind:          reflect.Slice,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "*github.com/xoctopus/x/typex/testdata.Int",
			},
		},
		{
			"NamedStringIntMap",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "map[github.com/xoctopus/x/typex/testdata.String]github.com/xoctopus/x/typex/testdata.Int",
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
			"NamedStringSet",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "map[github.com/xoctopus/x/typex/testdata.String]github.com/xoctopus/x/typex/testdata.EmptyStruct",
				Kind:          reflect.Map,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "github.com/xoctopus/x/typex/testdata.String",
				Elem:          "github.com/xoctopus/x/typex/testdata.EmptyStruct",
			},
		},
		{
			"NamedEmptyStruct",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmptyStruct",
				String:        "github.com/xoctopus/x/typex/testdata.EmptyStruct",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, true, false},
				ConvertibleTo: []bool{false, false, true, false, true, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
		{
			"NamedEmptyInterface",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmptyInterface",
				String:        "github.com/xoctopus/x/typex/testdata.EmptyInterface",
				Kind:          reflect.Interface,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
			},
		},
	}

	rtype := reflect.TypeOf(testdata.Basics{})
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