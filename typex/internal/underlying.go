package internal

import (
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

// NewTypesTypeByTypeID new types.Type by type id
func NewTypesTypeByTypeID(id string) (tt types.Type) {
	if t, ok := gTypesCache.Load(id); ok {
		return t.(types.Type)
	}

	defer func() {
		if tt == nil {
			tt = types.Typ[types.Invalid]
		}
		gTypesCache.Store(id, tt)
	}()

	if id == "" {
		return types.Typ[types.Invalid]
	}
	if t, ok := gBasicsCache.Load(id); ok {
		return t.(func() types.Type)()
	}

	// map[k]v
	if l := strings.Index(id, "map["); l == 0 {
		r := strings.Index(id, "]")
		k, v := id[4:r], id[r+1:]
		return types.NewMap(NewTypesTypeByTypeID(k), NewTypesTypeByTypeID(v))
	}

	if l := strings.Index(id, "["); l == 0 {
		// []t | [N]t
		r := strings.Index(id, "]")
		must.BeTrueWrap(r != -1, "failed to parse slice or array bracket: %s", id)
		t := NewTypesTypeByTypeID(id[r+1:])
		if r == l+1 {
			return types.NewSlice(t)
		}
		n, err := strconv.ParseInt(id[1:r], 10, 64)
		must.NoErrorWrap(err, "failed to parse array length: %s", id)
		return types.NewArray(t, n)
	} else if l == -1 {
		// path.type
		d := strings.LastIndex(id, ".")
		must.BeTrueWrap(d > 0, "failed to parse type id: %s", id)
		path, typename := id[0:d], id[d+1:]
		pkg := NewPackage(path)
		must.NotNilWrap(pkg, "failed to parse import path: %s", id)
		typ, ok := pkg.Scope().Lookup(typename).(*types.TypeName)
		must.BeTrueWrap(ok, "failed to lookup typename: %s", id)
		return types.NewTypeName(token.NoPos, pkg, typ.Name(), typ.Type()).Type()
	} else {
		// path.type[generic parameters]
		r := strings.LastIndex(id, "]")
		must.BeTrueWrap(r != -1, "failed to parse generic parameters bracket: %s", id)
		d := strings.LastIndex(id[0:l], ".")
		must.BeTrueWrap(d > 0, "failed to parse type id: %s", id)
		path, typename := id[0:d], id[d+1:l]
		pkg := NewPackage(path)
		must.NotNilWrap(pkg, "failed to parse import path: %s", id)
		typ, _ := pkg.Scope().Lookup(typename).(*types.TypeName)
		must.NotNilWrap(typ, "failed to lookup typename: %s", id)

		named, ok := typ.Type().(*types.Named)
		must.BeTrueWrap(ok, "failed to lookup typename, is not a named type: %s", id)

		clone := types.NewNamed(typ, typ.Type().Underlying(), nil)
		for i := 0; i < named.NumMethods(); i++ {
			clone.AddMethod(named.Method(i))
		}
		names := strings.Split(id[l+1:r], ",")
		params := named.TypeParams()
		must.BeTrueWrap(
			len(names) > 0 && params != nil && params.Len() == len(names),
			"invalid generic parameter list: %s", id,
		)
		pts := make([]*types.TypeParam, params.Len())
		for i := 0; i < params.Len(); i++ {
			pt := NewTypesTypeByTypeID(names[i])
			must.BeTrueWrap(
				pt != nil && pt != types.Typ[types.Invalid],
				"invalid generic parameter: %s", id,
			)
			pts[i] = types.NewTypeParam(params.At(i).Obj(), pt)
		}
		clone.SetTypeParams(pts)
		return clone
	}
}

// NewTypesTypeFromReflectType new types.Type from reflect.Type
func NewTypesTypeFromReflectType(t reflect.Type) types.Type {
	if t == nil {
		return nil
	}

	refs := 0
	if t.Name() == "" {
		for t.Kind() == reflect.Pointer {
			refs++
			t = t.Elem()
		}
	}

	tpe, pkg := t.Name(), t.PkgPath()
	if tpe == "error" && pkg == "" {
		return must.BeTrueV(gBasicsCache.Load("error")).(func() types.Type)()
	}

	var tt types.Type
	if pkg != "" {
		tt = NewTypesTypeByTypeID(t.PkgPath() + "." + t.Name())
	} else {
		tt = Underlying(t)
	}

	for refs > 0 {
		tt = types.NewPointer(tt)
		refs--
	}
	return tt
}

func Underlying(t reflect.Type) types.Type {
	if k, ok := KindsR2G[t.Kind()]; ok {
		return types.Typ[k]
	}

	switch t.Kind() {
	case reflect.Array:
		return types.NewArray(NewTypesTypeFromReflectType(t.Elem()), int64(t.Len()))
	case reflect.Slice:
		return types.NewSlice(NewTypesTypeFromReflectType(t.Elem()))
	case reflect.Map:
		return types.NewMap(
			NewTypesTypeFromReflectType(t.Key()),
			NewTypesTypeFromReflectType(t.Elem()),
		)
	case reflect.Chan:
		dir := types.ChanDir(0)
		switch t.ChanDir() {
		case reflect.RecvDir:
			dir = types.RecvOnly
		case reflect.SendDir:
			dir = types.SendOnly
		case reflect.BothDir:
			dir = types.SendRecv
		}
		return types.NewChan(dir, NewTypesTypeFromReflectType(t.Elem()))
	case reflect.Func:
		ins := make([]*types.Var, t.NumIn())
		for i := range ins {
			ti := t.In(i)
			ins[i] = types.NewParam(
				0, NewPackage(ti.PkgPath()), "",
				NewTypesTypeFromReflectType(ti),
			)
		}
		outs := make([]*types.Var, t.NumOut())
		for i := range outs {
			ti := t.Out(i)
			outs[i] = types.NewParam(0, NewPackage(ti.PkgPath()), "", NewTypesTypeFromReflectType(ti))
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
				NewTypesTypeFromReflectType(f.Type).(*types.Signature),
			)
		}
		return types.NewInterfaceType(fns, nil).Complete()
	case reflect.Struct:
		n := t.NumField()
		fields, tags := make([]*types.Var, n), make([]string, n)
		for i := range fields {
			f := t.Field(i)
			fields[i] = types.NewField(
				0,
				NewPackage(f.PkgPath),
				f.Name,
				NewTypesTypeFromReflectType(f.Type),
				f.Anonymous,
			)
			tags[i] = string(f.Tag)
		}
		return types.NewStruct(fields, tags)
	default:
		return nil
	}
}
