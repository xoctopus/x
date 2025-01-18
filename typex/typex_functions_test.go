package typex_test

import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/typex/testdata"
)

func TestTypes_Functions(t *testing.T) {
	cases := []Case{
		{
			"Func",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Func",
				String:        "github.com/xoctopus/x/typex/testdata.Func",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    true,
				NumIn:         3,
				Ins: []string{
					"string",
					"string",
					"[]int",
				},
				NumOut: 2,
				Outs: []string{
					"github.com/xoctopus/x/typex/testdata.Boolean",
					"error",
				},
			},
		},
		{
			"Curry",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Curry",
				String:        "github.com/xoctopus/x/typex/testdata.Curry",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    true,
				NumIn:         2,
				Ins: []string{
					"github.com/xoctopus/x/typex/testdata.String",
					"[]fmt.Stringer",
				},
				NumOut: 1,
				Outs:   []string{"func() string"},
			},
		},
		{
			"Uname",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "func() func() func() string",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    false,
				NumIn:         0,
				NumOut:        1,
				Outs:          []string{"func() func() string"},
			},
		},
		{
			"Max",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "Max[int]",
				String:        "github.com/xoctopus/x/typex/testdata.Max[int]",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						Name: "Compute",
						Type: "func(github.com/xoctopus/x/typex/testdata.Max[int], ...int) int",
					},
				},
				IsVariadic: true,
				NumIn:      1,
				Ins:        []string{"[]int"},
				NumOut:     1,
				Outs:       []string{"int"},
			},
		},
		{
			"CompareInt",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "func(int, int) int",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    false,
				NumIn:         2,
				Ins:           []string{"int", "int"},
				NumOut:        1,
				Outs:          []string{"int"},
			},
		},
		{
			"CompareNamedString",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "func(github.com/xoctopus/x/typex/testdata.String, github.com/xoctopus/x/typex/testdata.String) int",
				Kind:          reflect.Func,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    false,
				Key:           "nil",
				Elem:          "nil",
				IsVariadic:    false,
				NumIn:         2,
				Ins:           []string{"github.com/xoctopus/x/typex/testdata.String", "github.com/xoctopus/x/typex/testdata.String"},
				NumOut:        1,
				Outs:          []string{"int"},
			},
		},
	}

	RunCase(t, cases, testdata.FunctionCases)
}
