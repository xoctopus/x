package testdata

import (
	"fmt"
	"net"
)

const SIZE = 32

type (
	EmptyArray                    [0]interface{}
	Array                         [1]string
	AnyArray                      [2]any
	StringArray                   [4]String
	StringPtrArray                [8]*String
	IntPtrDefArray                [16]IntPtrDef
	SizedArray                    [SIZE]Boolean
	Map                           map[string]string
	StringIntMap                  map[String]Int
	Slice                         []string
	ErrorSlice                    []Error
	TypedArray[T any]             [1]T
	IntegerArray[T Integer]       [10]T
	TypedSizedArray[T Numeric]    [SIZE]T
	TypedSlice[T fmt.Stringer]    []T
	TypedMap[K comparable, V any] map[K]V
)

type CompositeBasics struct {
	EmptyArray      EmptyArray
	Array           Array
	AnyArray        AnyArray
	StringArray     StringArray
	StringPtrArray  StringPtrArray
	IntPtrDefArray  IntPtrDefArray
	SizedArray      SizedArray
	Map             Map
	StringIntMap    StringIntMap
	Slice           Slice
	ErrorSlice      ErrorSlice
	TypedArray      TypedArray[int]
	IntegerArray    IntegerArray[int8]
	TypedSizedArray TypedSizedArray[float32]
	TypedSlice      TypedSlice[net.Addr]
	TypedMap        TypedMap[string, int]
	TypedMap2       TypedMap[String, int]
	StructSlice     []SimpleStruct
	StructPtrSlice  []*SimpleStruct
}
