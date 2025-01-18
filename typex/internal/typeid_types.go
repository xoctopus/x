package internal

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

func InstTypeID(t types.Type, args []types.Type, wraps ...bool) string {
	var wrap = len(wraps) > 0 && wraps[0]
	switch x := t.(type) {
	case *types.Basic:
		return types.Typ[x.Kind()].String()
	case *types.Alias:
		return InstTypeID(x.Rhs(), args, wrap)
	case *types.Array:
		return fmt.Sprintf("[%d]%s", x.Len(), InstTypeID(x.Elem(), args, wrap))
	case *types.Chan:
		return fmt.Sprintf("%s%s", ChanDirPrefix(x), InstTypeID(x.Elem(), args, wrap))
	case *types.Map:
		return fmt.Sprintf("map[%s]%s", InstTypeID(x.Key(), args, wrap), InstTypeID(x.Elem(), args, wrap))
	case *types.Pointer:
		return fmt.Sprintf("*%s", InstTypeID(x.Elem(), args, wrap))
	case *types.Slice:
		return fmt.Sprintf("[]%s", InstTypeID(x.Elem(), args, wrap))
	case *types.TypeParam:
		return InstTypeID(args[x.Index()], args, wrap)
	case *types.Interface:
		if x.NumMethods() == 0 {
			return "interface {}"
		}
		b := strings.Builder{}
		b.WriteString("interface { ")
		for i := range x.NumMethods() {
			if i > 0 {
				b.WriteString("; ")
			}
			name := x.Method(i).Name()
			sig := InstTypeID(x.Method(i).Signature(), args, wrap)
			b.WriteString(name + sig[4:])
		}
		b.WriteString(" }")
		return b.String()
	case *types.Named:
		id := x.Obj().Name()
		if pkg := x.Obj().Pkg(); pkg != nil {
			id = pkg.Path() + "." + id
		}
		if wrap {
			id = WrapNamedTypeID(id)
		}
		b := strings.Builder{}
		b.WriteString(id)
		if x.TypeParams().Len() == 0 {
			return b.String()
		}
		b.WriteString("[")
		_args := args
		if x.TypeArgs().Len() > 0 {
			must.BeTrue(x.TypeArgs().Len() == x.TypeParams().Len())
			_args = make([]types.Type, x.TypeArgs().Len())
			for i := range x.TypeArgs().Len() {
				_args[i] = x.TypeArgs().At(i)
			}
		}
		for i, arg := range _args {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(InstTypeID(arg, args, wrap))
		}
		b.WriteString("]")
		return b.String()
	case *types.Signature:
		b := strings.Builder{}
		b.WriteString("func(")
		for i := range x.Params().Len() {
			if i > 0 {
				b.WriteString(", ")
			}
			p := x.Params().At(i).Type()
			if i == x.Params().Len()-1 && x.Variadic() {
				b.WriteString("...")
				p = p.(*types.Slice).Elem()
			}
			b.WriteString(InstTypeID(p, args, wrap))
		}
		b.WriteString(")")
		if x.Results().Len() > 0 {
			b.WriteString(" ")
		}
		if x.Results().Len() > 1 {
			b.WriteString("(")
		}
		for i := range x.Results().Len() {
			if i > 0 {
				b.WriteString(", ")
			}
			r := x.Results().At(i).Type()
			b.WriteString(InstTypeID(r, args, wrap))
		}
		if x.Results().Len() > 1 {
			b.WriteString(")")
		}
		return b.String()
	case *types.Struct:
		if x.NumFields() == 0 {
			return "struct {}"
		}
		b := strings.Builder{}
		b.WriteString("struct { ")
		for i := range x.NumFields() {
			if i > 0 {
				b.WriteString("; ")
			}
			f := x.Field(i)
			if !f.Anonymous() {
				b.WriteString(f.Name())
				b.WriteString(" ")
			}
			b.WriteString(InstTypeID(f.Type(), args, wrap))
			if tag := x.Tag(i); tag != "" {
				b.WriteString(" ")
				b.WriteString(tag)
			}
		}
		b.WriteString(" }")
		return b.String()
	default: // *types.Union *types.Tuple
		panic(errors.Errorf("unknown instantiate type: %T", x))
	}
}

func TypesTypeID(t types.Type, wraps ...bool) string {
	if tps, ok := t.(HasTypeParams); ok {
		if tas, ok := t.(HasTypeArgs); ok {
			must.BeTrue(tps.TypeParams().Len() == tas.TypeArgs().Len())
			if tas.TypeArgs().Len() > 0 {
				args := make([]types.Type, tas.TypeArgs().Len())
				for i := range tas.TypeArgs().Len() {
					args[i] = tas.TypeArgs().At(i)
				}
				InstTypeID(t, args, wraps...)
			}
		}
	}
	return InstTypeID(t, nil, wraps...)
}

func NewTypesType(t reflect.Type) types.Type {
	return NewTypesTypeByID(ReflectTypeID(t, true))
}

func NewTypesTypeByID(wrapped string) (tt types.Type) {
	if t, ok := gTypesCache.Load(wrapped); ok {
		return t.(types.Type)
	}

	defer func() {
		gTypesCache.Store(wrapped, tt)
	}()

	if wrapped == "" {
		return types.Typ[types.Invalid]
	}

	t := NewWrappedType(wrapped)
	switch t.kind {
	case reflect.Array:
		return types.NewArray(NewTypesTypeByID(t.elem.String()), int64(t.len))
	case reflect.Chan:
		return types.NewChan(t.dir, NewTypesTypeByID(t.elem.String()))
	case reflect.Func:
		return types.NewSignatureType(nil, nil, nil, t.Params(), t.Results(), t.variadic)
	case reflect.Interface:
		return types.NewInterfaceType(t.Methods(), nil).Complete()
	case reflect.Map:
		return types.NewMap(NewTypesTypeByID(t.key.String()), NewTypesTypeByID(t.elem.String()))
	case reflect.Pointer:
		return types.NewPointer(NewTypesTypeByID(t.elem.String()))
	case reflect.Slice:
		return types.NewSlice(NewTypesTypeByID(t.elem.String()))
	case reflect.Struct:
		return types.NewStruct(t.Fields())
	default:
		pkg := NewPackage(t.pkg)
		must.BeTrueWrap(pkg != nil, "failed to new package: %s(%s)", wrapped, t.pkg)
		obj, _ := pkg.Scope().Lookup(t.typename).(*types.TypeName)
		must.BeTrueWrap(obj != nil, "failed to lookup typename: %s(%s)", wrapped, t.typename)
		typ, _ := obj.Type().(*types.Named)
		must.BeTrueWrap(obj != nil, "failed to lookup typename: %s(%s), should be `Named`", wrapped, t.typename)
		args := make([]types.Type, len(t.args))
		for i, arg := range t.args {
			args[i] = NewTypesTypeByID(arg.String())
		}
		must.BeTrue(len(args) == typ.TypeParams().Len())
		if len(args) > 0 {
			typ = Instantiate(typ, args).(*types.Named)
		}
		return typ
	}
}

func InstantiateUnderlying(t types.Type) types.Type {
	if _t, ok := t.(CanBeInstantiated); ok {
		if _t.TypeParams().Len() > 0 {
			must.BeTrue(_t.TypeArgs().Len() == _t.TypeParams().Len())
			args := make([]types.Type, _t.TypeArgs().Len())
			for i := range _t.TypeArgs().Len() {
				args[i] = _t.TypeArgs().At(i)
			}
			return Instantiate(t.Underlying(), args)
		}
	}
	return t.Underlying()
}

func Instantiate(t types.Type, args []types.Type) types.Type {
	if len(args) == 0 {
		return t
	}
	switch x := t.(type) {
	case *types.Basic:
		return t
	case *types.TypeParam:
		return args[x.Index()]
	case *types.Alias:
		return Instantiate(x.Rhs(), args)
	case *types.Array:
		return types.NewArray(Instantiate(x.Elem(), args), x.Len())
	case *types.Chan:
		return types.NewChan(x.Dir(), Instantiate(x.Elem(), args))
	case *types.Slice:
		return types.NewSlice(Instantiate(x.Elem(), args))
	case *types.Map:
		return types.NewMap(Instantiate(x.Key(), args), Instantiate(x.Elem(), args))
	case *types.Pointer:
		return types.NewPointer(Instantiate(x.Elem(), args))
	case *types.Interface:
		methods := make([]*types.Func, x.NumMethods())
		for i := range x.NumMethods() {
			m := x.Method(i)
			s := Instantiate(m.Signature(), args).(*types.Signature)
			methods[i] = types.NewFunc(0, m.Pkg(), m.Name(), s)
		}
		return types.NewInterfaceType(methods, nil)
	case *types.Tuple:
		vars := make([]*types.Var, x.Len())
		for i := range x.Len() {
			v := x.At(i)
			vars[i] = types.NewField(0, v.Pkg(), v.Name(), Instantiate(v.Type(), args), v.Anonymous())
		}
		return types.NewTuple(vars...)
	case *types.Signature:
		return types.NewSignatureType(
			nil, nil, nil,
			Instantiate(x.Params(), args).(*types.Tuple),
			Instantiate(x.Results(), args).(*types.Tuple),
			x.Variadic(),
		)
	case *types.Struct:
		fields := make([]*types.Var, x.NumFields())
		tags := make([]string, x.NumFields())
		for i := range x.NumFields() {
			v := x.Field(i)
			fields[i] = types.NewField(0, v.Pkg(), v.Name(), Instantiate(v.Type(), args), v.Anonymous())
			tags[i] = x.Tag(i)
		}
		return types.NewStruct(fields, tags)
	case *types.Named:
		if x.TypeParams().Len() == 0 {
			return x
		}
		_args := make([]types.Type, x.TypeParams().Len())
		if x.TypeArgs().Len() > 0 {
			must.BeTrue(x.TypeArgs().Len() == x.TypeParams().Len())
			for i := range x.TypeArgs().Len() {
				if p, ok := x.TypeArgs().At(i).(*types.TypeParam); ok {
					_args[i] = Instantiate(p, args)
				} else {
					_args[i] = x.TypeArgs().At(i)
				}
			}
		} else {
			must.BeTrue(x.TypeParams().Len() == len(args))
			_args = args
		}
		ins, err := types.Instantiate(nil, x, _args, true)
		must.NoError(err)
		return ins
	default:
		panic(errors.Errorf("unknown instantiate type: %T", x))
	}
}
