package testdata

import (
	"unsafe"
)

type (
	String         string
	Boolean        bool
	Int            int
	UnsafePointer  unsafe.Pointer
	Error          error
	Chan           chan Int
	SendOnlyChan   chan<- String
	RecvOnlyChan   <-chan Error
	EmptyStruct    struct{}
	EmptyInterface interface{}
	IntPtrDef      *int
)

type Basics struct {
	String              string
	Boolean             bool
	Int                 int
	UnsafePointer       unsafe.Pointer
	Error               error
	Chan                chan string
	SendOnlyChan        chan<- string
	RecvOnlyChan        <-chan string
	IntPtr              *int
	StringArray         [3]string
	IntSlice            []int32
	IntPtrSlice         []*int64
	IntStringMap        map[int]string
	IntSet              map[int]struct{}
	EmptyStruct         struct{}
	EmptyInterface      any
	Func                func(x int, y string, z ...any) (float32, error)
	Curry               func(input any) func() string
	NamedString         String
	NamedBoolean        Boolean
	NamedInt            Int
	NamedUnsafePointer  UnsafePointer
	NamedError          Error
	NamedChan           Chan
	NamedSendOnlyChan   SendOnlyChan
	NamedRecvOnlyChan   RecvOnlyChan
	NamedIntPtr         *Int
	NamedIntPtrDef      IntPtrDef
	NamedStringArray    [10]String
	NamedIntSlice       []Int
	NamedIntPtrSlice    []*Int
	NamedStringIntMap   map[String]Int
	NamedStringSet      map[String]EmptyStruct
	NamedEmptyStruct    EmptyStruct
	NamedEmptyInterface EmptyInterface
}
