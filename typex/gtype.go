package typex

import (
	"bytes"
	"go/ast"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/typex/internal"
)

func NewGType(t any) (tt *GType) {
	switch x := t.(type) {
	case *GType:
		tt = x
	case *RType:
		tt = &GType{t: internal.NewTypesTypeFromReflectType(x.t)}
	case types.Type:
		switch x := x.(type) {
		case *types.Alias:
			tt = &GType{t: x.Rhs()}
		case *types.TypeParam:
			tt = &GType{t: x.Constraint()}
		}
		tt = &GType{t: x}
	case reflect.Type:
		tt = &GType{t: internal.NewTypesTypeFromReflectType(x)}
	default:
		return nil
	}

	var params *types.TypeParamList
	if x, ok := tt.t.(interface{ TypeParams() *types.TypeParamList }); ok {
		params = x.TypeParams()
	}

	var (
		fields  = internal.FieldNamesOf(tt.t, 0)
		methods = internal.MethodsOf(tt.t, params)
	)
	for _, mm := range methods {
		directed, inherited := -1, make([]int, 0)
		for i, m := range mm {
			if m.Ref == tt.t {
				directed = i
				break
			}
			inherited = append(inherited, i)
		}
		var m *internal.Method
		if directed == -1 && len(inherited) > 1 {
			continue
		}
		if directed >= 0 {
			m = mm[directed]
		} else {
			m = mm[inherited[0]]
		}
		gm := &GMethod{Method: m, recv: tt}
		tt.methods = append(tt.methods, gm)
		if !gm.IsPtr {
			tt.valRecvMethods = append(tt.valRecvMethods, gm)
		}
	}
	sort.SliceStable(tt.methods, func(i, j int) bool {
		return tt.methods[i].Name() < tt.methods[j].Name()
	})
	sort.SliceStable(tt.valRecvMethods, func(i, j int) bool {
		return tt.valRecvMethods[i].Name() < tt.valRecvMethods[j].Name()
	})
	tt.fields = fields
	return tt
}

type GType struct {
	t              types.Type
	methods        []*GMethod // methods sorted methods with value and pointer receiver
	valRecvMethods []*GMethod // valRecvMethods sorted methods with value receiver
	fields         map[string][]types.Type
}

func (t *GType) Unwrap() any { return t.t }

func (t *GType) PkgPath() string {
	switch x := t.t.(type) {
	case *types.Named:
		if pkg := x.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
		if x.String() == "error" {
			return ""
		}
	case *types.Basic:
		if strings.HasPrefix(x.String(), "unsafe.") {
			return "unsafe"
		}
	}
	return ""
}

func (t *GType) Name() string {
	switch x := t.t.(type) {
	case *types.Named:
		buf := bytes.NewBuffer(nil)
		buf.WriteString(x.Obj().Name())

		if params := x.TypeParams(); params.Len() > 0 {
			buf.WriteByte('[')
			for i := range params.Len() {
				if i > 0 {
					buf.WriteByte(',')
				}
				if args := x.TypeArgs(); args != nil {
					buf.WriteString(NewT(NewGType(args.At(i))).String())
				} else {
					param := params.At(i).Constraint()
					switch p := param.(type) {
					case *types.Interface:
						if p.NumEmbeddeds() > 0 {
							buf.WriteString(NewT(NewGType(p.EmbeddedType(0))).String())
						} else {
							buf.WriteString(NewT(NewGType(p)).String())
						}
					case *types.Named:
						buf.WriteString(p.String())
					}
				}
			}
			buf.WriteByte(']')
		}
		return buf.String()
	case *types.Basic:
		k := t.Kind()
		if k == reflect.UnsafePointer {
			return "Pointer"
		}
		return k.String()
	}
	return ""
}

func (t *GType) String() string { return NewT(t).String() }

func (t *GType) Kind() reflect.Kind {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).Kind()
	case *types.Basic:
		return internal.KindsG2R[x.Kind()]
	case *types.Interface:
		return reflect.Interface
	case *types.Struct:
		return reflect.Struct
	case *types.Pointer:
		return reflect.Pointer
	case *types.Slice:
		return reflect.Slice
	case *types.Array:
		return reflect.Array
	case *types.Map:
		return reflect.Map
	case *types.Chan:
		return reflect.Chan
	case *types.Signature:
		return reflect.Func
	case *types.TypeParam:
		switch xp := x.Underlying().(type) {
		case *types.Interface:
			if xp.NumEmbeddeds() > 0 {
				return NewGType(xp.EmbeddedType(0)).Kind()
			}
			return NewGType(xp).Kind()
		case *types.Named:
			return NewGType(xp).Kind()
		default:
			panic(errors.Errorf("unexpected Kind TypeParam type: %s", reflect.TypeOf(xp)))
		}
	case *types.Alias:
		return NewGType(x.Rhs()).Kind()
	default: // reflect.Invalid
		panic(errors.Errorf("unexpected Kind type: %s", reflect.TypeOf(x)))
	}
}

func (t *GType) Implements(u Type) bool {
	switch x := u.(type) {
	case *GType:
		if xt, ok := x.t.Underlying().(*types.Interface); ok {
			return types.Implements(t.t, xt)
		}
	case *RType:
		if x == nil {
			return false
		}
		v, ptr := Type(t), false
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
			ptr = true
		}

		vt := v.Unwrap().(types.Type)
		if ptr {
			vt = types.NewPointer(vt)
		}

		it := NewGType(x.Unwrap()).Unwrap().(types.Type).Underlying()
		if _, ok := it.(*types.Interface); !ok {
			return false
		}

		return types.Implements(vt, it.(*types.Interface))
	}
	return false
}

func (t *GType) AssignableTo(u Type) bool {
	if x, ok := u.(*GType); ok {
		return types.AssignableTo(t.t, x.t)
	}
	return false
}

func (t *GType) ConvertibleTo(u Type) bool {
	if x, ok := u.(*GType); ok {
		return types.ConvertibleTo(t.t, x.t)
	}
	return false
}

func (t *GType) Comparable() bool {
	if x, ok := t.t.(*types.Named); ok {
		xx := NewGType(internal.Constrain(x.Underlying(), x.TypeParams()))
		return xx.Kind() == reflect.Struct || types.Comparable(xx.t)
	}
	return t.Kind() == reflect.Struct || types.Comparable(t.t)
}

func (t *GType) Key() Type {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(internal.Constrain(x.Underlying(), x.TypeParams())).Key()
	case interface{ Key() types.Type }:
		return NewGType(x.Key())
	}
	return nil
}

func (t *GType) Elem() Type {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(internal.Constrain(x.Underlying(), x.TypeParams())).Elem()
	case interface{ Elem() types.Type }:
		return NewGType(x.Elem())
	}
	return nil
}

func (t *GType) Len() int {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).Len()
	case *types.Array:
		return int(x.Len())
	}
	return 0
}

func (t *GType) NumField() int {
	switch x := t.t.(type) {
	case *types.Pointer:
		// return NewGType(x.Elem()).NumField()
	case *types.Named:
		return NewGType(x.Underlying()).NumField()
	case *types.Struct:
		return x.NumFields()
	}
	return 0
}

func (t *GType) Field(i int) StructField {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(internal.Constrain(x.Underlying(), x.TypeParams())).Field(i)
	case *types.Struct:
		return &GStructField{f: x.Field(i), tag: x.Tag(i)}
	}
	return nil
}

func (t *GType) FieldByName(name string) (StructField, bool) {
	return t.FieldByNameFunc(func(s string) bool { return name == s })
}

func (t *GType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	var (
		directf   StructField
		directc   int
		indirectf StructField
		indirectc int
	)
	for i := range t.NumField() {
		f := t.Field(i)
		if match(f.Name()) {
			directc++
			directf = f
			// potential match or matched multiple times
			if directc > 1 {
				return nil, false
			}
		}
		if f.Anonymous() {
			ft := f.Type()
			for ft.Kind() == reflect.Pointer {
				ft = f.Type().Elem()
			}
			if ff, ok := ft.FieldByNameFunc(match); ok {
				indirectc++
				if indirectc > 1 {
					indirectf = nil
				} else {
					indirectf = ff
				}
			}
		}
	}
	if directf != nil {
		return directf, true
	}
	if indirectf != nil {
		return indirectf, true
	}
	return nil, false
}

func (t *GType) NumMethod() int {
	if t.Kind() == reflect.Interface {
		switch x := t.t.(type) {
		case *types.Named:
			return x.Underlying().(*types.Interface).NumMethods()
		case *types.Interface:
			return x.NumMethods()
		}
	}
	methods := t.methods
	if _, ok := t.t.(*types.Pointer); !ok {
		methods = t.valRecvMethods
	}
	num := len(methods)
	for _, m := range methods {
		if _, ok := t.fields[m.Name()]; ok && m.Ref != t.t {
			num--
		}
	}
	return num
}

func (t *GType) Method(i int) Method {
	if t.Kind() == reflect.Interface {
		var (
			f      *types.Func
			params *types.TypeParamList
		)
		switch x := t.t.(type) {
		case *types.Named:
			f = x.Underlying().(*types.Interface).Method(i)
			params = x.TypeParams()
		case *types.Interface:
			f = x.Method(i)
		default:
			panic(errors.Errorf("unexpected Method type %s", reflect.TypeOf(t)))
		}
		ss := internal.Constrain(f.Signature(), params).(*types.Signature)
		return &GMethod{
			Method: &internal.Method{
				Fn: types.NewFunc(f.Pos(), f.Pkg(), f.Name(), ss),
			},
		}
	}
	_, ok := t.t.(*types.Pointer)
	if ok {
		if i >= 0 && i < len(t.methods) {
			return t.methods[i]
		}
	} else {
		if i >= 0 && i < len(t.valRecvMethods) {
			return t.valRecvMethods[i]
		}
	}
	return nil
}

func (t *GType) MethodByName(name string) (Method, bool) {
	for i := range t.NumMethod() {
		if m := t.Method(i); m.Name() == name {
			return m, true
		}
	}
	return nil, false
}

func (t *GType) IsVariadic() bool {
	if f, ok := t.t.Underlying().(*types.Signature); ok {
		return f.Variadic()
	}
	return false
}

func (t *GType) NumIn() int {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).NumIn()
	case *types.Signature:
		return x.Params().Len()
	}
	return 0
}

func (t *GType) In(i int) Type {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).In(i)
	case *types.Signature:
		return NewGType(x.Params().At(i).Type())
	}
	return nil
}

func (t *GType) NumOut() int {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).NumOut()
	case *types.Signature:
		return x.Results().Len()
	}
	return 0
}

func (t *GType) Out(i int) Type {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).Out(i)
	case *types.Signature:
		return NewGType(x.Results().At(i).Type())
	}
	return nil
}

type GStructField struct {
	f   *types.Var
	tag string
}

func (f *GStructField) PkgPath() string {
	if ast.IsExported(f.f.Name()) {
		return ""
	}
	if pkg := f.f.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (f *GStructField) Name() string {
	return f.f.Name()
}

func (f *GStructField) Type() Type {
	return NewGType(f.f.Type())
}

func (f *GStructField) Tag() reflect.StructTag {
	return reflect.StructTag(f.tag)
}

func (f *GStructField) Anonymous() bool {
	return f.f.Anonymous()
}

type GMethod struct {
	*internal.Method
	recv *GType
}

func (m *GMethod) PkgPath() string {
	if ast.IsExported(m.Fn.Name()) {
		return ""
	}
	if pkg := m.Fn.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (m *GMethod) Name() string {
	return m.Fn.Name()
}

func (m *GMethod) Type() Type {
	sig := m.Fn.Type().(*types.Signature)
	if m.recv == nil {
		return NewGType(sig)
	}

	params := make([]*types.Var, sig.Params().Len()+1)
	params[0] = types.NewVar(0, nil, "", m.recv.t)
	for i := 0; i < sig.Params().Len(); i++ {
		pi := sig.Params().At(i)
		params[i+1] = types.NewVar(pi.Pos(), pi.Pkg(), pi.Name(), pi.Type())
	}

	return NewGType(
		types.NewSignatureType(
			nil, nil, nil,
			types.NewTuple(params...), sig.Results(), sig.Variadic(),
		),
	)
}

var (
	_ types.Type = (*types.Alias)(nil)
	_ types.Type = (*types.Array)(nil)
	_ types.Type = (*types.Basic)(nil)
	_ types.Type = (*types.Chan)(nil)
	_ types.Type = (*types.Interface)(nil)
	_ types.Type = (*types.Map)(nil)
	_ types.Type = (*types.Named)(nil)
	_ types.Type = (*types.Pointer)(nil)
	_ types.Type = (*types.Slice)(nil)
	_ types.Type = (*types.Struct)(nil)
	_ types.Type = (*types.Tuple)(nil)
	_ types.Type = (*types.TypeParam)(nil)
	_ types.Type = (*types.Signature)(nil)
	_ types.Type = (*types.Union)(nil)
)
