package testdata

import (
	"context"
	"encoding"
	"fmt"

	"github.com/xoctopus/x/typex/testdata/xoxo"
)

type (
	String     string
	Boolean    bool
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uintptr    uintptr
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
)

type (
	Array [1]string
	Map   map[string]string
	Slice []string
	Chan  chan string
	Func  func(x, y string) bool
)

func F() {}

type Interface interface {
	String() string
}

type Struct struct {
	Interface
	a       string
	A       string `json:"a"`
	B       string `json:"b"`
	Boolean `json:"bool,omitempty"`
	xoxo.Part
	Part2 Part `json:",omitempty"`
}

func (Struct) String() string {
	return ""
}

func (Struct) WithArg(arg any) {}

type Part struct {
	C string `json:"c"`
}

func (Part) Value() string {
	return ""
}

func (*Part) PtrValue() string {
	return ""
}

type Compose struct {
	Struct
}

type Enum int

const (
	ENUM__ONE Enum = iota + 1 // 1
	ENUM__TWO                 // 2
)

func (e *Enum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "ONE":
		*e = ENUM__ONE
	case "TWO":
		*e = ENUM__TWO
	}
	return fmt.Errorf("unknown")
}

func (e Enum) MarshalText() ([]byte, error) {
	switch e {
	case ENUM__ONE:
		return []byte("ONE"), nil
	case ENUM__TWO:
		return []byte("TWO"), nil
	}
	return nil, fmt.Errorf("unknown")
}

func (e Enum) String() string {
	switch e {
	case ENUM__ONE:
		return "ONE"
	case ENUM__TWO:
		return "TWO"
	}
	return ""
}

type MixedInterface interface {
	encoding.TextMarshaler
	Stringify(ctx context.Context, vs ...any) string
	Add(a, b string) string
	Bytes() []byte
	str() string
}

type AnySlice[V any] []V

func (s AnySlice[V]) Each() {}

type AnyArray[V any] [2]V

type AnyMap[K comparable, V any] map[K]V

type EnumMap = AnyMap[Enum, any]

type AnyStruct[V any] struct {
	Struct
	Name V
}

var Var string
