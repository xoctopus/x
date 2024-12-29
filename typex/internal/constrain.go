package internal

import (
	"go/types"
	"reflect"

	"github.com/pkg/errors"
)

type HasTypeParams interface {
	TypeParams() *types.TypeParamList
}

func Constrain(t types.Type, params *types.TypeParamList) types.Type {
	if params.Len() == 0 {
		return t
	}

	switch x := t.(type) {
	case *types.Basic:
		return t
	case *types.Map:
		return types.NewMap(Constrain(x.Key(), params), Constrain(x.Elem(), params))
	case *types.Array:
		return types.NewArray(Constrain(x.Elem(), params), x.Len())
	case *types.Slice:
		return types.NewSlice(Constrain(x.Elem(), params))
	case *types.Struct:
		n := x.NumFields()
		fields, tags := make([]*types.Var, n), make([]string, n)
		for i := 0; i < n; i++ {
			f := x.Field(i)
			fields[i] = types.NewField(f.Pos(), f.Pkg(), f.Name(), Constrain(f.Type(), params), f.Embedded())
			tags[i] = x.Tag(i)
		}
		return types.NewStruct(fields, tags)
	case *types.Named:
		methods := make([]*types.Func, x.NumMethods())
		for i := 0; i < x.NumMethods(); i++ {
			mi := x.Method(i)
			methods[i] = types.NewFunc(mi.Pos(), mi.Pkg(), mi.Name(), mi.Signature())
		}
		tparams := make([]*types.TypeParam, x.TypeParams().Len())
		for i := 0; i < x.TypeParams().Len(); i++ {
			pi := params.At(x.TypeParams().At(i).Index())
			tparams[i] = types.NewTypeParam(pi.Obj(), pi.Constraint())
		}
		cloned := types.NewNamed(x.Obj(), x.Underlying(), methods)
		cloned.SetTypeParams(tparams)
		return cloned
	case *types.TypeParam:
		switch xx := params.At(x.Index()).Constraint().(type) {
		case *types.Interface:
			if xx.NumEmbeddeds() > 0 {
				return Constrain(xx.EmbeddedType(0), params)
			}
			return Constrain(xx, params)
		case *types.Named:
			return Constrain(xx, params)
		default:
			panic(errors.Errorf("unexpected Constrain TypeParam type: %s", reflect.TypeOf(x)))
		}
	case *types.Signature:
		var ps, rs *types.Tuple
		if _params := x.Params(); _params.Len() > 0 {
			sparams := make([]*types.Var, 0, _params.Len())
			for i := 0; i < _params.Len(); i++ {
				vt := _params.At(i)
				v := types.NewVar(0, nil, vt.Name(), Constrain(vt.Type(), params))
				sparams = append(sparams, v)
			}
			ps = types.NewTuple(sparams...)
		}

		if _results := x.Results(); _results.Len() > 0 {
			results := make([]*types.Var, 0, _results.Len())
			for i := 0; i < _results.Len(); i++ {
				vt := _results.At(i)
				v := types.NewVar(0, nil, vt.Name(), Constrain(vt.Type(), params))
				results = append(results, v)
			}
			rs = types.NewTuple(results...)
		}
		return types.NewSignatureType(nil, nil, nil, ps, rs, x.Variadic())
	case *types.Interface:
		var (
			embeddeds []types.Type
			methods   []*types.Func
		)

		if x.NumEmbeddeds() > 0 {
			embeddeds = make([]types.Type, x.NumEmbeddeds())
			for i := 0; i < x.NumEmbeddeds(); i++ {
				embeddeds[i] = Constrain(x.EmbeddedType(i), params)
			}
		}
		if x.NumMethods() > 0 {
			methods = make([]*types.Func, x.NumMethods())
			for i := 0; i < x.NumMethods(); i++ {
				m := x.Method(i)
				methods[i] = types.NewFunc(m.Pos(), m.Pkg(), m.Name(), Constrain(m.Signature(), params).(*types.Signature))
			}
		}
		return types.NewInterfaceType(methods, embeddeds)
	case *types.Alias:
		panic(errors.Errorf("unimplemented Constrain of *types.Alias"))
	default:
		panic(errors.Errorf("unexpected Constrain type: %s", reflect.TypeOf(x)))
	}
}
