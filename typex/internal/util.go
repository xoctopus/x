package internal

import (
	"go/types"
	"reflect"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/xoctopus/x/misc/must"
)

var (
	gTypesCache    = sync.Map{}
	gBasicsCache   = sync.Map{}
	gPackagesCache = sync.Map{}
)

var (
	LoadFiles   = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	LoadImports = LoadFiles | packages.NeedImports
	LoadTypes   = LoadImports | packages.NeedTypes | packages.NeedTypesSizes
)

var (
	KindsR2G = map[reflect.Kind]types.BasicKind{
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
	KindsG2R = map[types.BasicKind]reflect.Kind{
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
)

func init() {
	for i := range types.Typ {
		basic := types.Typ[i]
		gBasicsCache.Store(
			types.TypeString(basic, nil),
			func() types.Type { return basic },
		)
	}
	gBasicsCache.Store("interface {}", func() types.Type {
		return types.NewInterfaceType(nil, nil)
	})
	gBasicsCache.Store("any", func() types.Type {
		return types.NewInterfaceType(nil, nil)
	})
	gBasicsCache.Store("comparable", func() types.Type {
		return types.NewInterfaceType(nil, nil)
	})
	gBasicsCache.Store("error", func() types.Type {
		return NewPackage("errors").
			Scope().
			Lookup("New").
			Type().
			Underlying().(*types.Signature).
			Results().
			At(0).
			Type()
	})
}

func NewPackage(path string) *types.Package {
	if path == "" {
		return nil
	}

	if v, ok := gPackagesCache.Load(path); ok {
		return v.(*types.Package)
	}

	pkgs, err := packages.Load(&packages.Config{
		Overlay: make(map[string][]byte),
		Tests:   true,
		Mode:    65535,
	}, path)
	must.NoErrorWrap(err, "failed to load packages from %s", path)

	// TODO enable check loading package error
	// if len(pkgs[0].Errors) != 0 {
	// 	return nil
	// }

	pkg := pkgs[0].Types
	gPackagesCache.Store(path, pkg)

	return pkg
}
