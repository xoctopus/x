package typex

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/typex/internal"
)

func NewGType(t any) (tt *GType) {
	switch x := t.(type) {
	case reflect.Type:
		tt = NewGType(internal.NewTypesType(x))
	case types.Type:
		switch x.(type) {
		case *types.Union, *types.Tuple, *types.TypeParam:
			panic(errors.Errorf("NewGType with unexpected types.Type: %T", x))
		default:
			if alias, ok := x.(*types.Alias); ok {
				x = alias.Rhs()
			}
			tt = &GType{t: x}
		}
	default:
		panic(errors.Errorf("NewGType with unknown type: %T", t))
	}
	tt.methods = internal.ScanMethods(tt.t)
	return tt
}

type GType struct {
	t       types.Type
	args    []types.Type
	methods *internal.Methods
}

func (t *GType) Unwrap() any { return t.t }

func (t *GType) PkgPath() string {
	switch x := t.t.(type) {
	case internal.HasObj:
		if pkg := x.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
	case *types.Basic:
		if x.Kind() == types.UnsafePointer {
			return "unsafe"
		}
	}
	return ""
}

func (t *GType) Name() string {
	switch x := t.t.(type) {
	case *types.Named:
		name := internal.TypesTypeID(t.t)
		params := ""
		if idx := strings.Index(name, "["); idx != -1 {
			params = name[idx:]
			name = name[0:idx]
		}
		if idx := strings.LastIndex(name, "."); idx != -1 {
			name = name[idx+1:]
		}
		return name + params
	case *types.Basic:
		return x.Name()
	}
	return ""
}

func (t *GType) String() string {
	return internal.TypesTypeID(t.t)
}

func (t *GType) Kind() reflect.Kind {
	switch x := t.t.(type) {
	case *types.Named:
		return NewGType(x.Underlying()).Kind()
	case *types.Basic:
		return internal.BasicKindsG2R[x.Kind()]
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
	}
	panic(errors.Errorf("unexpected kind type: %T", t.t))
}

func (t *GType) Implements(u Type) bool {
	switch x := u.(type) {
	case *GType:
		if ut, ok := x.t.Underlying().(*types.Interface); ok {
			return types.Implements(t.t, ut)
		}
	case *RType:
		return t.Implements(NewGType(x.t))
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
	return types.Comparable(internal.InstantiateUnderlying(t.t))
}

func (t *GType) Key() Type {
	switch x := t.t.(type) {
	case internal.HasKey:
		return NewGType(x.Key())
	case *types.Named:
		return NewGType(internal.InstantiateUnderlying(x)).Key()
	default:
		return nil
	}
}

func (t *GType) Elem() Type {
	switch x := t.t.(type) {
	case internal.HasElem:
		return NewGType(x.Elem())
	case *types.Named:
		return NewGType(internal.InstantiateUnderlying(x)).Elem()
	default:
		return nil
	}
}

func (t *GType) Len() int {
	switch x := t.t.(type) {
	case internal.HasLen64:
		return int(x.Len())
	case *types.Named:
		return NewGType(x.Underlying()).Len()
	default:
		return 0
	}
}

func (t *GType) NumField() int {
	switch x := t.t.(type) {
	case internal.HasFields:
		return x.NumFields()
	case *types.Named:
		return NewGType(x.Underlying()).NumField()
	default:
		return 0
	}
}

func (t *GType) Field(i int) StructField {
	switch x := t.t.(type) {
	case internal.HasFields:
		if i > x.NumFields() {
			return nil
		}
		return &GStructField{
			t: t.t,
			f: &internal.Field{Var: x.Field(i), Tag: x.Tag(i)},
		}
	case *types.Named:
		f := NewGType(internal.InstantiateUnderlying(x)).Field(i)
		if ff, ok := f.(*GStructField); ok {
			ff.t = t.t
			return f
		}
		return nil
	default:
		return nil
	}
}

func (t *GType) FieldByName(name string) (StructField, bool) {
	ff := internal.FieldByName(t.t, func(v string) bool { return v == name })
	if ff != nil {
		return &GStructField{t: t.t, f: ff}, true
	}
	return nil, false
}

func (t *GType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	ff := internal.FieldByName(t.t, match)
	if ff != nil {
		return &GStructField{t: t.t, f: ff}, true
	}
	return nil, false
}

func (t *GType) NumMethod() int {
	return t.methods.NumMethod()
}

func (t *GType) Method(i int) Method {
	if f := t.methods.Method(i); f != nil {
		return &GMethod{f: f, t: t.t}
	}
	return nil
}

func (t *GType) MethodByName(name string) (Method, bool) {
	if f := t.methods.MethodByName(name); f != nil {
		return &GMethod{f: f, t: t.t}, true
	}
	return nil, false
}

func (t *GType) IsVariadic() bool {
	switch x := t.t.(type) {
	case internal.Function:
		return x.Variadic()
	case *types.Named:
		return NewGType(x.Underlying()).IsVariadic()
	default:
		return false
	}
}

func (t *GType) NumIn() int {
	switch x := t.t.(type) {
	case internal.Function:
		return x.Params().Len()
	case *types.Named:
		return NewGType(x.Underlying()).NumIn()
	default:
		return 0
	}
}

func (t *GType) In(i int) Type {
	switch x := t.t.(type) {
	case internal.Function:
		if i < x.Params().Len() {
			return NewGType(x.Params().At(i).Type())
		}
		return nil
	case *types.Named:
		return NewGType(internal.InstantiateUnderlying(x)).In(i)
	default:
		return nil
	}
}

func (t *GType) NumOut() int {
	switch x := t.t.(type) {
	case internal.Function:
		return x.Results().Len()
	case *types.Named:
		return NewGType(x.Underlying()).NumOut()
	default:
		return 0
	}
}

func (t *GType) Out(i int) Type {
	switch x := t.t.(type) {
	case internal.Function:
		if i < x.Results().Len() {
			return NewGType(x.Results().At(i).Type())
		}
		return nil
	case *types.Named:
		return NewGType(internal.InstantiateUnderlying(x)).Out(i)
	default:
		return nil
	}
}

type GStructField struct {
	t types.Type
	f *internal.Field
}

func (f *GStructField) PkgPath() string {
	if pkg := f.f.Pkg(); pkg != nil && !f.f.Exported() {
		return pkg.Path()
	}
	return ""
}

func (f *GStructField) Name() string {
	return f.f.Var.Name()
}

func (f *GStructField) Type() Type {
	return NewGType(f.f.Var.Type())
}

func (f *GStructField) Tag() reflect.StructTag {
	return reflect.StructTag(f.f.Tag)
}

func (f *GStructField) Anonymous() bool {
	return f.f.Var.Anonymous()
}

type GMethod struct {
	t types.Type
	f *types.Func
}

func (m *GMethod) PkgPath() string {
	// unexported methods was hidden in static analysis
	// if pkg := m.f.Pkg(); pkg != nil && !m.f.Exported() {
	// 	return pkg.Path()
	// }
	return ""
}

func (m *GMethod) Name() string {
	return m.f.Name()
}

func (m *GMethod) Type() Type {
	s := m.f.Signature()
	params := make([]*types.Var, 0, s.Params().Len()+1)
	if _, ok := m.t.Underlying().(*types.Interface); !ok {
		params = append(params, types.NewParam(0, nil, "", m.t))
	}
	for i := range s.Params().Len() {
		params = append(params, s.Params().At(i))
	}
	s = types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(params...),
		s.Results(),
		s.Variadic(),
	)
	return NewGType(s)
}
