package internal

import (
	"go/token"
	"go/types"
	"reflect"

	"github.com/pkg/errors"
)

type Method struct {
	IsPtr bool
	Fn    *types.Func
	Ref   types.Type
}

func (m *Method) Name() string {
	return m.Fn.Name()
}

type HasMethods interface {
	NumMethods() int
	Method(int) *types.Func
}

func MethodsOf(u types.Type, params *types.TypeParamList) map[string][]*Method {
	combine := func(m1, m2 map[string][]*Method) map[string][]*Method {
		for _, mm := range m2 {
			for _, m := range mm {
				m1[m.Name()] = append(m1[m.Name()], m)
			}
		}
		return m1
	}
	switch t := u.(type) {
	case HasMethods: // *types.Named *types.Interface
		methods := make(map[string][]*Method)
		tt := t.(types.Type)
		for i := 0; i < t.NumMethods(); i++ {
			m := t.Method(i)
			name := m.Name()
			if !m.Exported() {
				continue
			}
			sig := Constrain(m.Signature(), params).(*types.Signature)
			fn := types.NewFunc(token.NoPos, m.Pkg(), m.Name(), sig)
			_, ptr := m.Type().(*types.Signature).Recv().Type().(*types.Pointer)
			methods[name] = append(methods[name], &Method{Fn: fn, IsPtr: ptr, Ref: t.(types.Type)})
		}
		if underlying := tt.Underlying(); underlying != tt {
			return combine(methods, MethodsOf(underlying, params))
		}
		return methods
	case *types.Struct:
		methods := make(map[string][]*Method)
		for i := 0; i < t.NumFields(); i++ {
			if f := t.Field(i); f.Anonymous() {
				methods = combine(methods, MethodsOf(f.Type(), params))
			}
		}
		return methods
	case *types.Pointer:
		methods := MethodsOf(t.Elem(), params)
		for _, mm := range methods {
			for _, m := range mm {
				if m.Ref == t.Elem() {
					m.Ref = t
				}
			}
		}
		return methods
	case *types.Alias:
		return MethodsOf(t.Rhs(), params)
	case *types.Basic, *types.Signature, *types.Slice, *types.Array, *types.Map, *types.Chan:
		return nil
	default:
		panic(errors.Errorf("unexpected MethodsOf type: %s", reflect.TypeOf(u)))
	}
}

func FieldNamesOf(u types.Type, lv int) map[string][]types.Type {
	switch t := u.(type) {
	case *types.Named:
		return FieldNamesOf(t.Underlying(), lv)
	case *types.Struct:
		fields := map[string][]types.Type{}
		for i := 0; i < t.NumFields(); i++ {
			name := t.Field(i).Name()
			fields[name] = append(fields[name], t)
			if f := t.Field(i); f.Anonymous() {
				_names := FieldNamesOf(f.Type(), lv+1)
				for name := range _names {
					fields[name] = append(fields[name], f.Type())
				}
			}
		}
		return fields
	case *types.Pointer:
		if lv > 0 {
			return FieldNamesOf(t.Elem(), lv)
		}
		return nil
	case *types.Basic, *types.Signature, *types.Slice, *types.Array, *types.Map, *types.Chan:
		return nil
	case *types.Interface:
		return nil
	case *types.Alias:
		return FieldNamesOf(t.Rhs(), lv)
	default:
		panic(errors.Errorf("unexpect FieldNamesOf %s", reflect.TypeOf(t)))
	}
}
