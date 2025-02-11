package internal

import (
	"go/ast"
	"go/types"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

var (
	gTypesCache    = sync.Map{}
	gPackagesCache = sync.Map{}
)

var (
	BasicKindsR2G = map[reflect.Kind]types.BasicKind{
		reflect.Invalid:       types.Invalid,
		reflect.Bool:          types.Bool,
		reflect.Int:           types.Int,
		reflect.Int8:          types.Int8,
		reflect.Int16:         types.Int16,
		reflect.Int32:         types.Int32,
		reflect.Int64:         types.Int64,
		reflect.Uint:          types.Uint,
		reflect.Uint8:         types.Uint8,
		reflect.Uint16:        types.Uint16,
		reflect.Uint32:        types.Uint32,
		reflect.Uint64:        types.Uint64,
		reflect.Uintptr:       types.Uintptr,
		reflect.Float32:       types.Float32,
		reflect.Float64:       types.Float64,
		reflect.Complex64:     types.Complex64,
		reflect.Complex128:    types.Complex128,
		reflect.String:        types.String,
		reflect.UnsafePointer: types.UnsafePointer,
	}
	BasicKindsG2R = map[types.BasicKind]reflect.Kind{
		types.Invalid:        reflect.Invalid,
		types.Bool:           reflect.Bool,
		types.UntypedBool:    reflect.Bool,
		types.Int:            reflect.Int,
		types.UntypedInt:     reflect.Int,
		types.Int8:           reflect.Int8,
		types.Int16:          reflect.Int16,
		types.Int32:          reflect.Int32, // includes types.Rune
		types.UntypedRune:    reflect.Int32,
		types.Int64:          reflect.Int64,
		types.Uint:           reflect.Uint,
		types.Uint8:          reflect.Uint8, // includes types.Byte
		types.Uint16:         reflect.Uint16,
		types.Uint32:         reflect.Uint32,
		types.Uint64:         reflect.Uint64,
		types.Uintptr:        reflect.Uintptr,
		types.Float32:        reflect.Float32,
		types.UntypedFloat:   reflect.Float32,
		types.Float64:        reflect.Float64,
		types.Complex64:      reflect.Complex64,
		types.UntypedComplex: reflect.Complex64,
		types.Complex128:     reflect.Complex128,
		types.String:         reflect.String,
		types.UntypedString:  reflect.String,
		types.UnsafePointer:  reflect.UnsafePointer,
	}
	gBasicKindNames = map[string]reflect.Kind{
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
		"string":         reflect.String,
		"unsafe.Pointer": reflect.UnsafePointer,
		"rune":           reflect.Int32,
		"byte":           reflect.Uint8,
		"error":          reflect.Interface,
	}
	gChanDirAst2Types = map[ast.ChanDir]types.ChanDir{
		ast.RECV: types.RecvOnly,
		ast.SEND: types.SendOnly,
		0:        types.SendRecv,
	}
	gChanDirPrefix = map[any]string{
		types.SendOnly:  "chan<- ",
		types.RecvOnly:  "<-chan ",
		types.SendRecv:  "chan ",
		reflect.SendDir: "chan<- ",
		reflect.RecvDir: "<-chan ",
		reflect.BothDir: "chan ",
	}
)

func init() {
	for i := range types.Typ {
		basic := types.Typ[i]
		gTypesCache.Store(basic.String(), basic)
	}
	gTypesCache.Store("interface {}", types.NewInterfaceType(nil, nil))
	gTypesCache.Store("any", types.NewInterfaceType(nil, nil))
	gTypesCache.Store("rune", types.Typ[types.Rune])
	gTypesCache.Store("byte", types.Typ[types.Byte])
	gTypesCache.Store("error", NewPackage("errors").Scope().Lookup("New").Type().Underlying().(*types.Signature).Results().At(0).Type())
}

func ChanDirPrefix(t any) string {
	switch x := t.(type) {
	case interface{ ChanDir() reflect.ChanDir }:
		return gChanDirPrefix[x.ChanDir()]
	case interface{ Dir() types.ChanDir }:
		return gChanDirPrefix[x.Dir()]
	case types.ChanDir, reflect.ChanDir:
		return gChanDirPrefix[x]
	default:
		panic(errors.Errorf("unexpected ChanDir type: %T", t))
	}
}

type HasTypeParams interface {
	TypeParams() *types.TypeParamList
}

type HasTypeArgs interface {
	TypeArgs() *types.TypeList
}

type CanBeInstantiated interface {
	TypeParams() *types.TypeParamList
	TypeArgs() *types.TypeList
}

type HasPkg interface {
	Pkg() *types.Package
}

type HasPkgPath interface {
	PkgPath() string
}

type HasObj interface {
	Obj() *types.TypeName
}

type HasKey interface {
	Key() types.Type
}

type HasElem interface {
	Elem() types.Type
}

type HasLen interface {
	Len() int
}

type HasLen64 interface {
	Len() int64
}

type HasFields interface {
	NumFields() int
	Field(int) *types.Var
	Tag(int) string
}

type HasMethods interface {
	NumMethods() int
	Method(int) *types.Func
}

type Function interface {
	Params() *types.Tuple
	Results() *types.Tuple
	Variadic() bool
	Recv() *types.Var
}

type NamedBacktrace struct {
	names []*types.Named
}

func (b *NamedBacktrace) Push(v *types.Named) bool {
	for _, named := range b.names {
		if types.Identical(named, v) {
			return false
		}
	}
	b.names = append(b.names, v)
	return true
}

func (b *NamedBacktrace) Top() *types.Named {
	if len(b.names) == 0 {
		return nil
	}
	return b.names[len(b.names)-1]
}

func (b *NamedBacktrace) Len() int {
	return len(b.names)
}
