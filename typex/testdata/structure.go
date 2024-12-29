package testdata

import "fmt"

type SimpleStruct struct {
	A              string
	B              String
	unexported     any
	Name           fmt.Stringer
	HasTag         Int `tag:"tagKey,otherLabel"`
	EmptyInterface     // anonymous
}

type EmbedSimpleStruct struct {
	SimpleStruct
	Name string
}

type EmbedInterface struct{ fmt.Stringer }

type HasValReceiverMethods struct{}

func (v HasValReceiverMethods) Name() string { return "Name" }

func (v HasValReceiverMethods) String() string { return "String" }

func (v HasValReceiverMethods) Bytes() []byte { return nil }

type HasPtrReceiverMethods struct{ SimpleStruct }

func (v *HasPtrReceiverMethods) SetA(a string) { v.A = a }

func (v *HasPtrReceiverMethods) SetB(b String) { v.B = b }

func (v *HasPtrReceiverMethods) UnmarshalText([]byte) error { return nil }

type HasValAndPtrReceiverMethods struct{ some any }

func (v HasValAndPtrReceiverMethods) MarshalText() ([]byte, error) { return nil, nil }

func (v *HasValAndPtrReceiverMethods) UnmarshalText([]byte) error { return nil }

type hasNameMethod struct{}

func (v hasNameMethod) Name() string { return "hasNameMethod" }

type hasNamePtrMethod struct{}

func (v *hasNamePtrMethod) Name() string { return "hasNamePtrMethod" }

type hasNameField struct{ Name String }

type hasStringMethod struct{}

func (v hasStringMethod) String() string { return "hasStringMethod" }

type hasStringField struct{ String }

// EmbedFieldOverwritten1 hasNameField.Name was overwritten by self Name field
type EmbedFieldOverwritten1 struct {
	hasNameField
	Name string
}

// EmbedFieldOverwritten2 *hasNameField.Name was overwritten by self Name field
type EmbedFieldOverwritten2 struct {
	Name string
	*hasNameField
}

// EmbedFieldOverwritten3 hasNameField.Name was overwritten by self Name method
type EmbedFieldOverwritten3 struct {
	hasNameField
}

func (EmbedFieldOverwritten3) Name() String { return "EmbedFieldOverwritten3" }

// InheritedMethodOverwritten1 hasNameMethod.Name() was overwritten by self Name field
type InheritedMethodOverwritten1 struct {
	hasNameMethod
	Name string
}

// InheritedMethodOverwritten2 *hasNameMethod.Name() was overwritten by self Name field
type InheritedMethodOverwritten2 struct {
	Name String
	*hasNameMethod
}

// InheritedMethodOverwritten3 hasNamePtrMethod.Name() was overwritten by self Name field
type InheritedMethodOverwritten3 struct {
	hasNamePtrMethod
	Name Int
}

// InheritedMethodOverwritten4 *hasNamePtrMethod.Name() was overwritten by self Name method
type InheritedMethodOverwritten4 struct {
	*hasNamePtrMethod
}

func (InheritedMethodOverwritten4) Name() Int { return 0 }

// InheritedMethodOverwritten5 *hasNamePtrMethod.Name() was overwritten by self Name method
type InheritedMethodOverwritten5 struct {
	hasNamePtrMethod
}

func (*InheritedMethodOverwritten5) Name() float32 { return 0 }

type EmbedFieldAndMethodConflict struct {
	hasNameField
	hasNameMethod
}

type hasNameField2 struct{ Name string }

type EmbedFieldsConflict1 struct {
	hasNameField
	hasNameField2
}

type EmbedFieldsConflict2 struct {
	hasNameField
	*hasNameField2
}

type InheritedMethodsConflict1 struct {
	hasNameMethod
	hasNamePtrMethod
}

type InheritedMethodsConflict2 struct {
	hasNameMethod
	*hasNamePtrMethod
}

type hasNameMethod2 struct{}

func (hasNameMethod2) Name() String { return "hasNameMethod2" }

type InheritedMethodsConflict3 struct {
	hasNameMethod
	hasNameMethod2
}

type InheritedMethodsConflict4 struct {
	hasNameMethod
	hasNamePtrMethod
}

type InheritedMethodsConflict5 struct {
	*hasNameMethod
	*hasNamePtrMethod
}

type InheritedMethodsConflict6 struct {
	*hasNameMethod
	*hasNamePtrMethod
}

func (*InheritedMethodsConflict6) Name() any { return nil }

type UncomparableStruct struct{}

type Structures struct {
	SimpleStruct                   SimpleStruct
	EmbedSimpleStruct              EmbedSimpleStruct
	EmbedInterface                 EmbedInterface
	EmbedInterfacePtr              *EmbedInterface
	HasValReceiverMethods          HasValReceiverMethods
	HasValReceiverMethodsPtr       *HasValReceiverMethods
	HasPtrReceiverMethods          HasPtrReceiverMethods
	HasPtrReceiverMethodsPtr       *HasPtrReceiverMethods
	HasValAndPtrReceiverMethods    HasValAndPtrReceiverMethods
	HasValAndPtrReceiverMethodsPtr *HasValAndPtrReceiverMethods
	EmbedFieldOverwritten1         EmbedFieldOverwritten1
	EmbedFieldOverwritten2         EmbedFieldOverwritten2
	EmbedFieldOverwritten3         EmbedFieldOverwritten3
	InheritedMethodOverwritten1    InheritedMethodOverwritten1
	InheritedMethodOverwritten2    InheritedMethodOverwritten2
	InheritedMethodOverwritten3    InheritedMethodOverwritten3
	InheritedMethodOverwritten4    InheritedMethodOverwritten4
	InheritedMethodOverwritten5    *InheritedMethodOverwritten5
	EmbedFieldsConflict1           EmbedFieldsConflict1
	EmbedFieldsConflict2           EmbedFieldsConflict2
	EmbedFieldAndMethodConflict    EmbedFieldAndMethodConflict
	InheritedMethodsConflict1      InheritedMethodsConflict1
	InheritedMethodsConflict2      InheritedMethodsConflict2
	InheritedMethodsConflict3      InheritedMethodsConflict3
	InheritedMethodsConflict4      *InheritedMethodsConflict4
	InheritedMethodsConflict5      InheritedMethodsConflict5
	InheritedMethodsConflict6      *InheritedMethodsConflict6
}
