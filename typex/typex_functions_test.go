package typex_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestTypes_Functions(t *testing.T) {
	cases := []struct {
		name string
		c    *CaseAssertion
	}{
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
	}

	rtype := reflect.TypeOf(testdata.Functions{})
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
