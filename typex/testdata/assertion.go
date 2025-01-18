package testdata

import "reflect"

type Case struct {
	Name string
	Type reflect.Type
}

var (
	BasicCases     []*Case
	CompositeCases []*Case
	FunctionCases  []*Case
	StructureCases []*Case
)

func init() {
	BasicCases = GenerateCases(Basics{})
	CompositeCases = GenerateCases(Composites{})
	FunctionCases = GenerateCases(Functions{})
	StructureCases = GenerateCases(Structures{})

	// add external cases here
	FunctionCases = append(
		FunctionCases,
		&Case{"CompareInt", reflect.TypeOf(CompareInt)},
		&Case{"CompareNamedString", reflect.TypeOf(CompareNamedString)},
	)
}

func GenerateCases(v any) []*Case {
	rt := reflect.TypeOf(v)
	cases := make([]*Case, rt.NumField())
	for i := range len(cases) {
		f := rt.Field(i)
		cases[i] = &Case{Name: f.Name, Type: f.Type}
	}
	return cases
}
