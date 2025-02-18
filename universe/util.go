package universe

import (
	"go/ast"
	"go/types"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/resultx"
)

func ChanDir(dir any) resultx.Result[types.ChanDir] {
	switch d := dir.(type) {
	case ast.ChanDir:
		switch d {
		case ast.SEND:
			return resultx.Succeed(types.SendOnly)
		case ast.RECV:
			return resultx.Succeed(types.RecvOnly)
		default:
			return resultx.Succeed(types.SendRecv)
		}
	case reflect.ChanDir:
		switch d {
		case reflect.RecvDir:
			return resultx.Succeed(types.RecvOnly)
		case reflect.SendDir:
			return resultx.Succeed(types.SendOnly)
		case reflect.BothDir:
			return resultx.Succeed(types.SendRecv)
		default:
			return resultx.Err[types.ChanDir](errors.Errorf("invalid reflect.ChanDir %d", d))
		}
	case types.ChanDir:
		return resultx.Succeed(d)
	default:
		return resultx.Err[types.ChanDir](errors.Errorf("invalid dir type: %T", d))
	}
}

var (
	gBasicName2Kind = map[string]reflect.Kind{
		"bool":           reflect.Bool,
		"int":            reflect.Int,
		"int8":           reflect.Int8,
		"int16":          reflect.Int16,
		"int32":          reflect.Int32,
		"int64":          reflect.Int64,
		"uint":           reflect.Uint,
		"uint8":          reflect.Uint8,
		"uint16":         reflect.Uint16,
		"uint32":         reflect.Uint32,
		"uint64":         reflect.Uint64,
		"uintptr":        reflect.Uintptr,
		"float32":        reflect.Float32,
		"float64":        reflect.Float64,
		"complex64":      reflect.Complex64,
		"complex128":     reflect.Complex128,
		"unsafe.Pointer": reflect.UnsafePointer,
		"string":         reflect.String,
		"rune":           reflect.Int32,
		"byte":           reflect.Uint8,
	}
	gChanDirPrefix = map[any]string{
		types.SendOnly:  "chan<- ",
		types.RecvOnly:  "<-chan ",
		types.SendRecv:  "chan ",
		reflect.SendDir: "chan<- ",
		reflect.RecvDir: "<-chan ",
		reflect.BothDir: "chan ",
	}
	gWrapID   sync.Map
	gUniverse sync.Map
	gPackages sync.Map
)

var (
	errNotNamedType   = errors.New("")
	errNotUnnamedType = errors.New("")
)
