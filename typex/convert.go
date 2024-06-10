package typex

import (
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/xoctopus/x/misc/must"
)

var (
	gGoTypeCache   = sync.Map{}
	gBasicsCache   = sync.Map{}
	gPackagesCache = sync.Map{}
)

var (
	LoadFiles   = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	LoadImports = LoadFiles | packages.NeedImports
	LoadTypes   = LoadImports | packages.NeedTypes | packages.NeedTypesSizes
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

func TypeByImportAndName(path, name string) types.Type {
	if path == "" {
		return TypeByID(name)
	}
	return TypeByID(path + "." + name)
}

func TypeByID(id string) (t types.Type) {
	if v, ok := gGoTypeCache.Load(id); ok {
		return v.(types.Type)
	}

	defer func() {
		if t == nil {
			t = types.Typ[types.Invalid]
		}
		gGoTypeCache.Store(id, t)
	}()

	if id == "" {
		return nil
	}

	if v, ok := gBasicsCache.Load(id); ok {
		t = v.(func() types.Type)()
		return
	}

	// map[x]
	if idxMapLsb := strings.Index(id, "map["); idxMapLsb == 0 {
		idxMapRsb := strings.Index(id, "]")
		if idxMapRsb == -1 {
			return nil
		}
		kt := TypeByID(id[4:idxMapRsb])
		if kt == nil || kt == types.Typ[types.Invalid] {
			return nil
		}
		vt := TypeByID(id[idxMapRsb+1:])
		if vt == nil || vt == types.Typ[types.Invalid] {
			return nil
		}
		t = types.NewMap(kt, vt)
		return t
	}

	// []x [n]x
	idxLsb := strings.Index(id, "[")
	if idxLsb == 0 {
		idxRsb := strings.Index(id, "]")
		if idxRsb == -1 {
			return nil
		}
		elem := TypeByID(id[idxRsb+1:])
		if elem == nil || elem == types.Typ[types.Invalid] {
			return nil
		}
		if idxRsb == idxLsb+1 {
			t = types.NewSlice(elem)
			return t
		}
		n, err := strconv.ParseInt(id[1:idxRsb], 10, 64)
		if err != nil {
			return t
		}
		t = types.NewArray(elem, n)
		return t
	}

	// path.Type
	if idxLsb == -1 {
		if idxDot := strings.LastIndex(id, "."); idxDot > 0 {
			pkg := NewPackage(id[0:idxDot])
			if pkg == nil {
				return nil
			}
			tpe, ok := pkg.Scope().Lookup(id[idxDot+1:]).(*types.TypeName)
			if !ok || tpe == nil {
				return nil
			}
			t = types.NewTypeName(token.NoPos, pkg, tpe.Name(), tpe.Type()).Type()
			// t = tpe.Type()
			return t
		}
		return t
	}

	// path.Type[generics parameters...]
	idxRsb := strings.Index(id, "]")
	if idxRsb == -1 {
		return nil
	}
	typeid := id[0:idxLsb]
	idxDot := strings.LastIndex(typeid, ".")
	if idxDot == -1 {
		return nil
	}
	pkg := NewPackage(typeid[:idxDot])
	if pkg == nil {
		return nil
	}
	tpe, ok := pkg.Scope().Lookup(typeid[idxDot+1:]).(*types.TypeName)
	if !ok || tpe == nil {
		return nil
	}
	named, ok := tpe.Type().(*types.Named)
	if !ok {
		return nil
	}
	cloned := types.NewNamed(tpe, tpe.Type().Underlying(), nil)
	for i := 0; i < named.NumMethods(); i++ {
		cloned.AddMethod(named.Method(i))
	}
	paramNames := strings.Split(id[idxLsb+1:idxRsb], ",")
	paramList := named.TypeParams()
	if len(paramNames) == 0 || paramList == nil || paramList.Len() != len(paramNames) {
		return nil
	}
	paramTypes := make([]*types.TypeParam, paramList.Len())
	for i := 0; i < paramList.Len(); i++ {
		paramType := TypeByID(paramNames[i])
		if paramType == nil || paramType == types.Typ[types.Invalid] {
			return nil
		}
		paramTypes[i] = types.NewTypeParam(paramList.At(i).Obj(), paramType)
	}
	cloned.SetTypeParams(paramTypes)
	t = cloned
	return t
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
		Mode:    LoadTypes,
	}, path)
	must.NoErrorWrap(err, "failed to load packages from %s", path)

	// if len(pkgs[0].Errors) != 0 {
	// 	return nil
	// }

	pkg := pkgs[0].Types
	gPackagesCache.Store(path, pkg)

	return pkg
}

func NewGoTypeFromReflectType(t reflect.Type) types.Type {
	if t == nil {
		return nil
	}

	refs := 0
	for t.Kind() == reflect.Pointer {
		refs++
		t = t.Elem()
	}

	tpe, pkg := t.Name(), t.PkgPath()
	if tpe == "error" && pkg == "" {
		return must.BeTrueV(gBasicsCache.Load("error")).(func() types.Type)()
	}

	var gt types.Type
	if pkg != "" {
		gt = TypeByID(pkg + "." + tpe)
	} else {
		gt = underlying(t)
	}

	for refs > 0 {
		gt = types.NewPointer(gt)
		refs--
	}
	return gt
}

func underlying(t reflect.Type) types.Type {
	if tt, ok := ReflectKindToTypesKind[t.Kind()]; ok {
		return types.Typ[tt]
	}

	switch t.Kind() {
	case reflect.Array:
		return types.NewArray(
			NewGoTypeFromReflectType(t.Elem()),
			int64(t.Len()),
		)
	case reflect.Chan:
		var dir = types.SendRecv
		switch t.ChanDir() {
		case reflect.RecvDir: // 1
			dir = types.RecvOnly // 2
		case reflect.SendDir: // 2
			dir = types.SendOnly // 1
		case reflect.BothDir: // 3
			dir = types.SendRecv // 0
		}
		return types.NewChan(dir, NewGoTypeFromReflectType(t.Elem()))
	case reflect.Slice:
		return types.NewSlice(
			NewGoTypeFromReflectType(t.Elem()),
		)
	case reflect.Map:
		return types.NewMap(
			NewGoTypeFromReflectType(t.Key()),
			NewGoTypeFromReflectType(t.Elem()),
		)
	case reflect.Func:
		ins := make([]*types.Var, t.NumIn())
		outs := make([]*types.Var, t.NumOut())
		for i := range ins {
			in := t.In(i)
			ins[i] = types.NewParam(
				0,
				NewPackage(in.PkgPath()),
				"",
				NewGoTypeFromReflectType(in),
			)
		}
		for i := range outs {
			out := t.Out(i)
			outs[i] = types.NewParam(
				0,
				NewPackage(out.PkgPath()),
				"",
				NewGoTypeFromReflectType(out),
			)
		}
		return types.NewSignatureType(
			nil, nil, nil,
			types.NewTuple(ins...),
			types.NewTuple(outs...),
			t.IsVariadic(),
		)
	case reflect.Interface:
		fns := make([]*types.Func, t.NumMethod())
		for i := range fns {
			f := t.Method(i)
			fns[i] = types.NewFunc(
				0,
				NewPackage(f.PkgPath),
				f.Name,
				NewGoTypeFromReflectType(f.Type).(*types.Signature),
			)
		}
		return types.NewInterfaceType(fns, nil).Complete()
	case reflect.Struct:
		fields := make([]*types.Var, t.NumField())
		tags := make([]string, len(fields))
		for i := range fields {
			f := t.Field(i)
			fields[i] = types.NewField(
				0,
				NewPackage(f.PkgPath),
				f.Name,
				NewGoTypeFromReflectType(f.Type),
				f.Anonymous,
			)
			tags[i] = string(f.Tag)
		}
		return types.NewStruct(fields, tags)
	default: // [never entered this case] reflect.Invalid, reflect.Pointer
		return nil
	}
}
