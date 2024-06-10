package typex

import (
	"bytes"
	"go/ast"
	"go/types"
	"reflect"
	"strings"
)

func NewGoType(t types.Type) *GoType {
	gt := &GoType{Type: t}
	if p, ok := t.(*types.TypeParam); ok {
		gt.Type = p.Constraint()
	}

	var scan func(t types.Type) []*GoMethod

	scan = func(t types.Type) []*GoMethod {
		methods := make([]*GoMethod, 0)
		switch x := t.(type) {
		case *types.Named:
			for i := range x.NumMethods() {
				m := &GoMethod{fn: x.Method(i)}
				if _, ok := m.fn.Type().(*types.Signature).Recv().Type().(*types.Pointer); ok {
					m.ptr = true
				}
				methods = append(methods, m)
			}
			return append(methods, scan(x.Underlying())...)
		case *types.Pointer:
			return append(methods, scan(x.Elem())...)
		case *types.Struct:
			for i := range x.NumFields() {
				if f := x.Field(i); f.Anonymous() {
					methods = append(methods, scan(f.Type())...)
				}
			}
			return methods
		}
		return nil
	}

	methods := scan(gt.Type)
	for _, m := range methods {
		if !m.ptr {
			gt.methods = append(gt.methods, m.fn)
		}
		gt.ptrMethods = append(gt.ptrMethods, m.fn)
	}
	return gt
}

type GoType struct {
	Type       types.Type
	methods    []*types.Func
	ptrMethods []*types.Func
}

var _ Type = (*GoType)(nil)

func (t *GoType) Unwrap() any {
	return t.Type
}

func (t *GoType) PkgPath() string {
	switch x := t.Type.(type) {
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
		return ""
	}
	return ""
}

func (t *GoType) Name() string {
	switch x := t.Type.(type) {
	case *types.Named:
		buf := bytes.NewBuffer(nil)
		buf.WriteString(x.Obj().Name())
		if params := x.TypeParams(); params.Len() > 0 {
			buf.WriteString("[")
			for i := range params.Len() {
				if i > 0 {
					buf.WriteByte(',')
				}
				p := params.At(i).Constraint().(*types.Interface)
				if p.NumEmbeddeds() > 0 {
					buf.WriteString(typename(NewGoType(p.EmbeddedType(0))))
				} else {
					buf.WriteString(typename(NewGoType(p)))
				}
			}
			buf.WriteString("]")
		}
		return buf.String()
	case *types.Basic:
		return x.Name()
	}
	return ""
}

func (t *GoType) String() string {
	return typename(t)
}

func (t *GoType) Kind() reflect.Kind {
	switch x := t.Type.(type) {
	case *types.Named:
		if pkg := x.Obj().Pkg(); pkg != nil &&
			pkg.Name() == "unsafe" &&
			x.Obj().Name() == "Pointer" {
			return reflect.UnsafePointer
		}
		return NewGoType(x.Underlying()).Kind()
	case *types.Interface:
		return reflect.Interface
	case *types.Pointer:
		return reflect.Pointer
	case *types.Struct:
		return reflect.Struct
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
	case *types.Basic:
		return TypesKindToReflectKind[x.Kind()]
	}
	return reflect.Invalid
}

func (t *GoType) Implements(u Type) bool {
	switch x := u.(type) {
	case *GoType:
		if xi, ok := x.Type.(*types.Interface); ok {
			return types.Implements(t.Type, xi)
		}
	case *ReflectType:
		var (
			tt  Type = t
			ptr      = false
		)
		for tt.Kind() == reflect.Pointer {
			tt = tt.Elem()
			ptr = true
		}
		if tt.PkgPath() == "" || x.PkgPath() == "" {
			return false
		}
		xi, ok := NewGoTypeFromReflectType(x.Type).Underlying().(*types.Interface)
		if !ok {
			return false
		}
		t2 := TypeByImportAndName(tt.PkgPath(), tt.Name())
		if ptr {
			t2 = types.NewPointer(t2)
		}
		return types.Implements(t2, xi)
	}
	return false
}

func (t *GoType) AssignableTo(u Type) bool {
	if x, ok := u.(*GoType); ok {
		return types.AssignableTo(t.Type, x.Type)
	}
	return false
}

func (t *GoType) ConvertibleTo(u Type) bool {
	if x, ok := u.(*GoType); ok {
		return types.ConvertibleTo(t.Type, x.Type)
	}
	return false
}

func (t *GoType) Comparable() bool {
	if x, ok := t.Type.(*types.Named); ok {
		return types.Comparable(ConstrainUnderlying(x.TypeParams(), x.Underlying()))
	}
	return t.Kind() == reflect.Struct ||
		types.Comparable(t.Type)
}

func (t *GoType) Key() Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(ConstrainUnderlying(x.TypeParams(), x.Underlying())).Key()
	case interface{ Key() types.Type }:
		return NewGoType(x.Key())
	}
	return nil
}

func (t *GoType) Elem() Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(ConstrainUnderlying(x.TypeParams(), x.Underlying())).Elem()
	case interface{ Elem() types.Type }:
		return NewGoType(x.Elem())
	}
	return nil
}

func (t *GoType) Len() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(x.Underlying()).Len()
	case *types.Array:
		return int(x.Len())
	}
	return 0
}

func (t *GoType) NumField() int {
	switch x := t.Type.(type) {
	case *types.Pointer:
		return NewGoType(x.Elem()).NumField()
	case *types.Named:
		return NewGoType(x.Underlying()).NumField()
	case *types.Struct:
		return x.NumFields()
	}
	return 0
}

func (t *GoType) Field(i int) StructField {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(ConstrainUnderlying(x.TypeParams(), x.Underlying())).Field(i)
	case *types.Struct:
		return &GoStructField{Var: x.Field(i), tag: x.Tag(i)}
	}
	return nil
}

func (t *GoType) FieldByName(name string) (StructField, bool) {
	return t.FieldByNameFunc(func(s string) bool { return name == s })
}

func (t *GoType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	for i := range t.NumField() {
		f := t.Field(i)
		if match(f.Name()) {
			return f, true
		}
		if f.Anonymous() {
			if sf, ok := f.Type().FieldByNameFunc(match); ok {
				return sf, true
			}
		}
	}
	return nil, false
}

func (t *GoType) NumMethod() int {
	if t.Kind() == reflect.Interface {
		switch x := t.Type.(type) {
		case *types.Named:
			return x.Underlying().(*types.Interface).NumMethods()
		case *types.Interface:
			return x.NumMethods()
		}
	}

	switch t.Type.(type) {
	case *types.Pointer:
		return len(t.ptrMethods)
	default:
		return len(t.methods)
	}
}

func (t *GoType) Method(i int) Method {
	if t.Kind() == reflect.Interface {
		switch x := t.Type.(type) {
		case *types.Named:
			return &GoMethod{fn: x.Underlying().(*types.Interface).Method(i)}
		case *types.Interface:
			return &GoMethod{fn: x.Method(i)}
		}
	}

	switch t.Type.(type) {
	case *types.Pointer:
		if i >= 0 && i < len(t.ptrMethods) {
			return &GoMethod{ptr: true, recv: t, fn: t.ptrMethods[i]}
		}
	default:
		if i >= 0 && i < len(t.methods) {
			return &GoMethod{recv: t, fn: t.methods[i]}
		}
	}
	return nil
}

func (t *GoType) MethodByName(name string) (Method, bool) {
	for i := range t.NumMethod() {
		if m := t.Method(i); m.Name() == name {
			return m, true
		}
	}
	return nil, false
}

func (t *GoType) IsVariadic() bool {
	if s, ok := t.Type.(*types.Signature); ok {
		return s.Variadic()
	}
	return false
}

func (t *GoType) NumIn() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(x.Underlying()).NumIn()
	case *types.Signature:
		return x.Params().Len()
	}
	return 0
}

func (t *GoType) In(i int) Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(x.Underlying()).In(i)
	case *types.Signature:
		return NewGoType(x.Params().At(i).Type())
	}
	return nil
}

func (t *GoType) NumOut() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(x.Underlying()).NumOut()
	case *types.Signature:
		return x.Results().Len()
	}
	return 0
}

func (t *GoType) Out(i int) Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return NewGoType(x.Underlying()).Out(i)
	case *types.Signature:
		return NewGoType(x.Results().At(i).Type())
	}
	return nil
}

type GoMethod struct {
	ptr  bool
	recv *GoType
	fn   *types.Func
}

var _ Method = (*GoMethod)(nil)

func (m *GoMethod) PkgPath() string {
	if ast.IsExported(m.Name()) {
		return ""
	}
	if pkg := m.fn.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (m *GoMethod) Name() string {
	return m.fn.Name()
}

func (m *GoMethod) Type() Type {
	sig := m.fn.Type().(*types.Signature)

	if m.recv == nil {
		return NewGoType(sig)
	}

	ins := make([]*types.Var, sig.Params().Len()+1)
	ins[0] = types.NewVar(0, nil, "", m.recv.Type)
	for i := range sig.Params().Len() {
		ins[i+1] = sig.Params().At(i)
	}
	return NewGoType(types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(ins...),
		sig.Results(),
		sig.Variadic(),
	))
}

type GoStructField struct {
	*types.Var
	tag string
}

var _ StructField = (*GoStructField)(nil)

func (f *GoStructField) PkgPath() string {
	if ast.IsExported(f.Name()) {
		return ""
	}
	if pkg := f.Var.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (f *GoStructField) Tag() reflect.StructTag {
	return reflect.StructTag(f.tag)
}

func (f *GoStructField) Type() Type {
	return NewGoType(f.Var.Type())
}

func ConstrainUnderlying(params *types.TypeParamList, underlying types.Type) types.Type {
	if params.Len() == 0 {
		return underlying
	}

	switch t := underlying.(type) {
	case *types.TypeParam:
		p := params.At(t.Index()).Constraint().(*types.Interface)
		if p.NumEmbeddeds() > 0 {
			return p.EmbeddedType(0)
		}
		return p
	case *types.Map:
		return types.NewMap(
			ConstrainUnderlying(params, t.Key()),
			ConstrainUnderlying(params, t.Elem()),
		)
	case *types.Slice:
		return types.NewSlice(
			ConstrainUnderlying(params, t.Elem()),
		)
	case *types.Array:
		return types.NewArray(
			ConstrainUnderlying(params, t.Elem()),
			t.Len(),
		)
	case *types.Struct:
		num := t.NumFields()
		tags, fields := make([]string, num), make([]*types.Var, num)
		for i := range num {
			f := t.Field(i)
			fields[i] = types.NewField(
				f.Pos(),
				f.Pkg(),
				f.Name(),
				ConstrainUnderlying(params, f.Type()),
				f.Embedded(),
			)
			tags[i] = t.Tag(i)
		}
		return types.NewStruct(fields, tags)
	}
	return underlying
}

var ReflectKindToTypesKind = map[reflect.Kind]types.BasicKind{
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

var TypesKindToReflectKind = map[types.BasicKind]reflect.Kind{
	types.Bool:           reflect.Bool,
	types.Int:            reflect.Int,
	types.Int8:           reflect.Int8,
	types.Int16:          reflect.Int16,
	types.Int32:          reflect.Int32,
	types.Int64:          reflect.Int64,
	types.Uint:           reflect.Uint,
	types.Uint8:          reflect.Uint8,
	types.Uint16:         reflect.Uint16,
	types.Uint32:         reflect.Uint32,
	types.Uint64:         reflect.Uint64,
	types.Uintptr:        reflect.Uintptr,
	types.Float32:        reflect.Float32,
	types.Float64:        reflect.Float64,
	types.Complex64:      reflect.Complex64,
	types.Complex128:     reflect.Complex128,
	types.String:         reflect.String,
	types.UnsafePointer:  reflect.UnsafePointer,
	types.UntypedBool:    reflect.Bool,
	types.UntypedInt:     reflect.Int,
	types.UntypedRune:    reflect.Int32,
	types.UntypedFloat:   reflect.Float32,
	types.UntypedComplex: reflect.Complex64,
	types.UntypedString:  reflect.String,
}
