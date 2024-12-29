package typex_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestTypes_CompositeStructures(t *testing.T) {
	cases := []struct {
		name string
		c    *CaseAssertion
	}{
		{
			"SimpleStruct",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "SimpleStruct",
				String:        "github.com/xoctopus/x/typex/testdata.SimpleStruct",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      6,
				Fields: []FieldAssertion{
					{
						Name: "A",
						Type: "string",
					},
					{
						Name: "B",
						Type: "github.com/xoctopus/x/typex/testdata.String",
					},
					{
						Name:    "unexported",
						Type:    "interface {}",
						PkgPath: "github.com/xoctopus/x/typex/testdata",
					},
					{
						Name: "Name",
						Type: "fmt.Stringer",
					},
					{
						Name: "HasTag",
						Type: "github.com/xoctopus/x/typex/testdata.Int",
						Tag:  `tag:"tagKey,otherLabel"`,
					},
					{
						Name:      "EmptyInterface",
						Type:      "github.com/xoctopus/x/typex/testdata.EmptyInterface",
						Anonymous: true,
					},
				},
			},
		},
		{
			"EmbedSimpleStruct",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedSimpleStruct",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedSimpleStruct",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						Name:      "SimpleStruct",
						Type:      "github.com/xoctopus/x/typex/testdata.SimpleStruct",
						Anonymous: true,
					},
					{
						Name: "Name",
						Type: "string",
					},
				},
			},
		},
		{
			"EmbedInterface",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedInterface",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedInterface",
				Kind:          reflect.Struct,
				Implements:    []bool{true, false, true, false, false, false},
				AssignableTo:  []bool{true, false, true, false, false, false},
				ConvertibleTo: []bool{true, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      1,
				Fields: []FieldAssertion{
					{
						Name:      "Stringer",
						Type:      "fmt.Stringer",
						Anonymous: true,
					},
				},
				NumMethod: 1,
				Methods: []MethodAssertion{
					{
						Name: "String",
						Type: "func(github.com/xoctopus/x/typex/testdata.EmbedInterface) string",
					},
				},
			},
		},
		{
			"EmbedInterfacePtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.EmbedInterface",
				Kind:          reflect.Pointer,
				Implements:    []bool{true, false, true, false, false, false},
				AssignableTo:  []bool{true, false, true, false, false, false},
				ConvertibleTo: []bool{true, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.EmbedInterface",
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						Name: "String",
						Type: "func(*github.com/xoctopus/x/typex/testdata.EmbedInterface) string",
					},
				},
			},
		},
		{
			"HasValReceiverMethods",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "HasValReceiverMethods",
				String:        "github.com/xoctopus/x/typex/testdata.HasValReceiverMethods",
				Kind:          reflect.Struct,
				Implements:    []bool{true, true, true, false, false, false},
				AssignableTo:  []bool{true, true, true, false, true, false},
				ConvertibleTo: []bool{true, true, true, false, true, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      0,
				NumMethod:     3,
				Methods: []MethodAssertion{
					{
						Name: "Bytes",
						Type: "func(github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) []uint8",
					},
					{
						Name: "Name",
						Type: "func(github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) string",
					},
					{
						Name: "String",
						Type: "func(github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) string",
					},
				},
			},
		},
		{
			"HasValReceiverMethodsPtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.HasValReceiverMethods",
				Kind:          reflect.Pointer,
				Implements:    []bool{true, true, true, false, false, false},
				AssignableTo:  []bool{true, true, true, false, false, false},
				ConvertibleTo: []bool{true, true, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.HasValReceiverMethods",
				NumField:      0,
				NumMethod:     3,
				Methods: []MethodAssertion{
					{
						Name: "Bytes",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) []uint8",
					},
					{
						Name: "Name",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) string",
					},
					{
						Name: "String",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasValReceiverMethods) string",
					},
				},
			},
		},
		{
			"HasPtrReceiverMethods",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "HasPtrReceiverMethods",
				String:        "github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      1,
				Fields: []FieldAssertion{
					{
						Name:      "SimpleStruct",
						Type:      "github.com/xoctopus/x/typex/testdata.SimpleStruct",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"HasPtrReceiverMethodsPtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods",
				NumField:      0,
				NumMethod:     3,
				Methods: []MethodAssertion{
					{
						Name: "SetA",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods, string)",
					},
					{
						Name: "SetB",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods, github.com/xoctopus/x/typex/testdata.String)",
					},
					{
						Name: "UnmarshalText",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasPtrReceiverMethods, []uint8) error",
					},
				},
			},
		},
		{
			"HasValAndPtrReceiverMethods",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "HasValAndPtrReceiverMethods",
				String:        "github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      1,
				Fields: []FieldAssertion{
					{
						PkgPath: "github.com/xoctopus/x/typex/testdata",
						Name:    "some",
						Type:    "interface {}",
					},
				},
				NumMethod: 1,
				Methods: []MethodAssertion{
					{
						Name: "MarshalText",
						Type: "func(github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods) ([]uint8, error)",
					},
				},
			},
		},
		{
			"HasValAndPtrReceiverMethodsPtr",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods",
				NumMethod:     2,
				Methods: []MethodAssertion{
					{
						Name: "MarshalText",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods) ([]uint8, error)",
					},
					{
						Name: "UnmarshalText",
						Type: "func(*github.com/xoctopus/x/typex/testdata.HasValAndPtrReceiverMethods, []uint8) error",
					},
				},
			},
		},
		{
			"EmbedFieldOverwritten1",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldOverwritten1",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldOverwritten1",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
					{
						Name: "Name",
						Type: "string",
					},
				},
				NumMethod: 0,
			},
		},
		{
			"EmbedFieldOverwritten2",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldOverwritten2",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldOverwritten2",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						Name: "Name",
						Type: "string",
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"EmbedFieldOverwritten3",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldOverwritten3",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldOverwritten3",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      1,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
				},
				NumMethod: 1,
				Methods: []MethodAssertion{
					{
						Name: "Name",
						Type: "func(github.com/xoctopus/x/typex/testdata.EmbedFieldOverwritten3) github.com/xoctopus/x/typex/testdata.String",
					},
				},
			},
		},
		{
			"InheritedMethodOverwritten1",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodOverwritten1",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten1",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
					{
						Name: "Name",
						Type: "string",
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodOverwritten2",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodOverwritten2",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten2",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						Name: "Name",
						Type: "github.com/xoctopus/x/typex/testdata.String",
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodOverwritten3",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodOverwritten3",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten3",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNamePtrMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNamePtrMethod",
						Anonymous: true,
					},
					{
						Name: "Name",
						Type: "github.com/xoctopus/x/typex/testdata.Int",
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodOverwritten4",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodOverwritten4",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten4",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      1,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNamePtrMethod",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNamePtrMethod",
						Anonymous: true,
					},
				},
				NumMethod: 1,
				Methods: []MethodAssertion{
					{
						Name: "Name",
						Type: "func(github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten4) github.com/xoctopus/x/typex/testdata.Int",
					},
				},
			},
		},
		{
			"InheritedMethodOverwritten5",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten5",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten5",
				NumField:      0,
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						Name: "Name",
						Type: "func(*github.com/xoctopus/x/typex/testdata.InheritedMethodOverwritten5) float32",
					},
				},
			},
		},
		{
			"EmbedFieldsConflict1",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldsConflict1",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldsConflict1",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField2",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField2",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"EmbedFieldsConflict2",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldsConflict2",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldsConflict2",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField2",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNameField2",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"EmbedFieldAndMethodConflict",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "EmbedFieldAndMethodConflict",
				String:        "github.com/xoctopus/x/typex/testdata.EmbedFieldAndMethodConflict",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameField",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameField",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodsConflict1",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodsConflict1",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict1",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNamePtrMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNamePtrMethod",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodsConflict2",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodsConflict2",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict2",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNamePtrMethod",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNamePtrMethod",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodsConflict3",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodsConflict3",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict3",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod2",
						Type:      "github.com/xoctopus/x/typex/testdata.hasNameMethod2",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodsConflict4",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict4",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict4",
				NumField:      0,
				NumMethod:     0,
			},
		},
		{
			"InheritedMethodsConflict5",
			&CaseAssertion{
				PkgPath:       "github.com/xoctopus/x/typex/testdata",
				Name:          "InheritedMethodsConflict5",
				String:        "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict5",
				Kind:          reflect.Struct,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "nil",
				NumField:      2,
				Fields: []FieldAssertion{
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNameMethod",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNameMethod",
						Anonymous: true,
					},
					{
						PkgPath:   "github.com/xoctopus/x/typex/testdata",
						Name:      "hasNamePtrMethod",
						Type:      "*github.com/xoctopus/x/typex/testdata.hasNamePtrMethod",
						Anonymous: true,
					},
				},
				NumMethod: 0,
			},
		},
		{
			"InheritedMethodsConflict6",
			&CaseAssertion{
				PkgPath:       "",
				Name:          "",
				String:        "*github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict6",
				Kind:          reflect.Pointer,
				Implements:    []bool{false, false, true, false, false, false},
				AssignableTo:  []bool{false, false, true, false, false, false},
				ConvertibleTo: []bool{false, false, true, false, false, false},
				Comparable:    true,
				Key:           "nil",
				Elem:          "github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict6",
				NumMethod:     1,
				Methods: []MethodAssertion{
					{
						Name: "Name",
						Type: "func(*github.com/xoctopus/x/typex/testdata.InheritedMethodsConflict6) interface {}",
					},
				},
			},
		},
	}

	rtype := reflect.TypeOf(testdata.Structures{})
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
