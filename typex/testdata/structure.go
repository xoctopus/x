package testdata

import (
	"fmt"
	"net"
)

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

// EmbedFieldsConflict1 both hasNameField and hasNameField has `Name` Field
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

type InheritedMethodAndFieldsConflict1 struct {
	hasNameField
	hasNameMethod
}

type InheritedMethodAndFieldsConflict2 struct {
	hasNameField
	hasNameMethod
}

func (InheritedMethodAndFieldsConflict2) Name() string { return "InheritedMethodAndFieldsConflict2" }

type InheritedMethodAndFieldsConflict3 struct {
	hasNameField
	hasNameMethod
	Name string
}

type InheritedMethodAndFieldsConflict4 struct {
	hasNameField
	hasNameField2
	hasNameMethod
}

type HasUnexportedField struct {
	str string
}

type HasUnexportedMethod struct{}

func (HasUnexportedMethod) str() string    { return "HasUnexportedMethod" }
func (HasUnexportedMethod) String() string { return "HasUnexportedMethod" }

type UncomparableStruct struct {
	v map[string]any
}

type Serialized[T CanBeSerialized] struct {
	data T
}

func (v Serialized[T]) String() string {
	switch data := any(v.data).(type) {
	case string:
		return data
	case []byte:
		return string(data)
	}
	return ""
}

func (v Serialized[T]) Data() T {
	return v.data
}

func (v Serialized[T]) Bytes() []byte {
	switch data := any(v.data).(type) {
	case string:
		return []byte(data)
	case []byte:
		return data
	}
	return nil
}

func NewNode[T any](v T) *BTreeNode[T] {
	return &BTreeNode[T]{
		Value: v,
		Left:  nil,
		Right: nil,
	}
}

type BTreeNode[T any] struct {
	Value  T
	Left   *BTreeNode[T]
	Right  *BTreeNode[T]
	Parent *BTreeNode[T]
}

func (n *BTreeNode[T]) InsertLeft(v T) *BTreeNode[T] {
	l := NewNode(v)
	l.Parent = n
	old := n.Left
	old.Parent = nil
	n.Left = l
	return old
}

type TypeParamPass[T1 any, T2 fmt.Stringer] struct {
	t1 T1
	t2 T2
	*BTreeNode[T2]
}

func (v *TypeParamPass[T1, T2]) Deal(t1 T1) T2 {
	v.t1 = t1
	return v.t2
}

type CircleEmbedsA struct {
	CircleEmbedsB
}

type CircleEmbedsB struct {
	*CircleEmbedsC
}

type CircleEmbedsC struct {
	CircleEmbedsA
	CircleEmbedsB
}

type Valuer[T any] struct {
	v T
}

func (v Valuer[T]) String() T {
	return v.v
}

type MapBTreeNode[V comparable] map[V]*BTreeNode[V]

type TypedSliceAlias = TypedSlice[net.Addr]

type Structures struct {
	SimpleStruct                         SimpleStruct
	SimpleStructPtr                      *SimpleStruct
	EmbedSimpleStruct                    EmbedSimpleStruct
	EmbedSimpleStructPtr                 *EmbedSimpleStruct
	EmbedInterface                       EmbedInterface
	EmbedInterfacePtr                    *EmbedInterface
	HasValReceiverMethods                HasValReceiverMethods
	HasValReceiverMethodsPtr             *HasValReceiverMethods
	HasPtrReceiverMethods                HasPtrReceiverMethods
	HasPtrReceiverMethodsPtr             *HasPtrReceiverMethods
	HasValAndPtrReceiverMethods          HasValAndPtrReceiverMethods
	HasValAndPtrReceiverMethodsPtr       *HasValAndPtrReceiverMethods
	EmbedFieldOverwritten1               EmbedFieldOverwritten1
	EmbedFieldOverwritten1Ptr            *EmbedFieldOverwritten1
	EmbedFieldOverwritten2               EmbedFieldOverwritten2
	EmbedFieldOverwritten2Ptr            *EmbedFieldOverwritten2
	EmbedFieldOverwritten3               EmbedFieldOverwritten3
	EmbedFieldOverwritten3Ptr            *EmbedFieldOverwritten3
	InheritedMethodOverwritten1          InheritedMethodOverwritten1
	InheritedMethodOverwritten1Ptr       *InheritedMethodOverwritten1
	InheritedMethodOverwritten2          InheritedMethodOverwritten2
	InheritedMethodOverwritten2Ptr       *InheritedMethodOverwritten2
	InheritedMethodOverwritten3          InheritedMethodOverwritten3
	InheritedMethodOverwritten3Ptr       *InheritedMethodOverwritten3
	InheritedMethodOverwritten4          InheritedMethodOverwritten4
	InheritedMethodOverwritten4Ptr       *InheritedMethodOverwritten4
	InheritedMethodOverwritten5          InheritedMethodOverwritten5
	InheritedMethodOverwritten5Ptr       *InheritedMethodOverwritten5
	EmbedFieldAndMethodConflict          EmbedFieldAndMethodConflict
	EmbedFieldAndMethodConflictPtr       *EmbedFieldAndMethodConflict
	EmbedFieldsConflict1                 EmbedFieldsConflict1
	EmbedFieldsConflict1Ptr              *EmbedFieldsConflict1
	EmbedFieldsConflict2                 EmbedFieldsConflict2
	EmbedFieldsConflict2Ptr              *EmbedFieldsConflict2
	InheritedMethodsConflict1            InheritedMethodsConflict1
	InheritedMethodsConflict1Ptr         *InheritedMethodsConflict1
	InheritedMethodsConflict2            InheritedMethodsConflict2
	InheritedMethodsConflict2Ptr         *InheritedMethodsConflict2
	InheritedMethodsConflict3            InheritedMethodsConflict3
	InheritedMethodsConflict3Ptr         *InheritedMethodsConflict3
	InheritedMethodsConflict4            InheritedMethodsConflict4
	InheritedMethodsConflict4Ptr         *InheritedMethodsConflict4
	InheritedMethodsConflict5            InheritedMethodsConflict5
	InheritedMethodsConflict5Ptr         *InheritedMethodsConflict5
	InheritedMethodsConflict6            InheritedMethodsConflict6
	InheritedMethodsConflict6Ptr         *InheritedMethodsConflict6
	InheritedMethodAndFieldsConflict1    InheritedMethodAndFieldsConflict1
	InheritedMethodAndFieldsConflict1Ptr *InheritedMethodAndFieldsConflict1
	InheritedMethodAndFieldsConflict2    InheritedMethodAndFieldsConflict2
	InheritedMethodAndFieldsConflict2Ptr *InheritedMethodAndFieldsConflict2
	InheritedMethodAndFieldsConflict3    InheritedMethodAndFieldsConflict3
	InheritedMethodAndFieldsConflict3Ptr *InheritedMethodAndFieldsConflict3
	InheritedMethodAndFieldsConflict4    InheritedMethodAndFieldsConflict4
	InheritedMethodAndFieldsConflict4Ptr *InheritedMethodAndFieldsConflict4
	HasUnexportedField                   HasUnexportedField
	HasUnexportedFieldPtr                *HasUnexportedField
	HasUnexportedMethod                  HasUnexportedMethod
	HasUnexportedMethodPtr               *HasUnexportedMethod
	SerializedString                     Serialized[string]
	SerializedStringPtr                  *Serialized[string]
	SerializedBytes                      Serialized[[]uint8]
	SerializedBytesPtr                   *Serialized[[]uint8]
	IntBTreeNode                         BTreeNode[int]
	IntBTreeNodePtr                      *BTreeNode[int]
	SerializedBytesBTreeNode             BTreeNode[Serialized[[]byte]]
	SerializedBytesBTreeNodePtr          *BTreeNode[Serialized[[]byte]]
	CircleEmbedsA                        CircleEmbedsA
	CircleEmbedsAPtr                     *CircleEmbedsA
	CircleEmbedsB                        CircleEmbedsB
	CircleEmbedsBPtr                     *CircleEmbedsB
	CircleEmbedsC                        CircleEmbedsC
	CircleEmbedsCPtr                     *CircleEmbedsC
	TypeParamPass1                       TypeParamPass[string, net.Addr]
	TypeParamPass1Ptr                    *TypeParamPass[string, net.Addr]
	TypeParamPass2                       TypeParamPass[Serialized[string], Serialized[[]byte]]
	TypeParamPass2Ptr                    *TypeParamPass[Serialized[string], Serialized[[]byte]]
	UncomparableStruct                   UncomparableStruct
	UncomparableStructPtr                *UncomparableStruct
	UnameStruct                          struct{ Int }
	UnameStructPtr                       *struct{ Int }
	ValuerString                         Valuer[string]
	ValuerStringPtr                      *Valuer[string]
	MapBTreeNodeInt                      MapBTreeNode[int]
	MapBTreeNodeIntPtr                   *MapBTreeNode[int]
	TypedSliceAlias                      TypedSliceAlias
	TypedSliceAliasPtr                   *TypedSliceAlias
}

type T[E comparable] struct {
	_ int
	_ TypedSliceAlias
	_ [3]int
	_ chan int
	_ []int
	_ map[string]int
	_ *int
	_ error
	_ func(...E) E
	_ TypedArray[E]
	_ Chan
	_ Max[E]
	_ TypeParamPass[E, net.Addr]
	_ interface{ Value() E }
}
