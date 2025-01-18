package internal

import (
	"go/types"
	"sort"

	"github.com/xoctopus/x/misc/must"
)

func combine[K comparable, V any](m1, m2 map[K][]V) map[K][]V {
	if len(m1) == 0 {
		return m2
	}
	if len(m2) == 0 {
		return m1
	}
	for k, v := range m2 {
		m1[k] = append(m1[k], v...)
	}
	return m1
}

func InspectMethods(t types.Type, names *NamedBacktrace) (methods map[string][]*types.Func) {
	must.BeTrue(names != nil)
	switch x := t.(type) {
	case HasMethods:
		methods = make(map[string][]*types.Func)
		for i := range x.NumMethods() {
			m := x.Method(i)
			methods[m.Name()] = append(methods[m.Name()], m)
		}
		if u, ok := x.(*types.Named); ok && names.Push(u) {
			methods = combine(methods, InspectMethods(u.Underlying(), names))
		}
	case *types.Struct:
		for i := range x.NumFields() {
			if f := x.Field(i); f.Anonymous() {
				methods = combine(methods, InspectMethods(f.Type(), names))
			}
		}
	case *types.Pointer:
		methods = InspectMethods(x.Elem(), names)
	case *types.Alias:
		methods = InspectMethods(x.Rhs(), names)
	}
	return methods
}

func InspectFields(t types.Type, names *NamedBacktrace) (fields map[string][]*types.Var) {
	must.BeTrue(names != nil)
	switch x := t.(type) {
	case *types.Named:
		if names.Push(x) {
			fields = InspectFields(x.Underlying(), names)
		}
	case *types.Struct:
		fields = make(map[string][]*types.Var)
		for i := range x.NumFields() {
			f := x.Field(i)
			fields[f.Name()] = append(fields[f.Name()], f)
			if f.Anonymous() {
				fields = combine(fields, InspectFields(f.Type(), names))
			}
		}
	case *types.Pointer:
		fields = InspectFields(x.Elem(), names)
	case *types.Alias:
		fields = InspectFields(x.Rhs(), names)
	}
	return fields
}

func DirectedMethod(m *types.Func, t types.Type) bool {
	r := m.Signature().Recv().Type()
	if x, ok := r.(*types.Pointer); ok {
		r = x.Elem()
	}
	if x, ok := t.(*types.Pointer); ok {
		t = x.Elem()
	}
	return types.Identical(r, t)
}

func RequirePointerRecv(m *types.Func, t types.Type) bool {
	_, ptrRecv := m.Signature().Recv().Type().(*types.Pointer)
	_, ptrType := t.(*types.Pointer)
	return !ptrRecv || ptrType
}

func ScanMethods(t types.Type) *Methods {
	var (
		methods = InspectMethods(t, &NamedBacktrace{})
		fields  = InspectFields(t, &NamedBacktrace{})
		finals  = make([]*types.Func, 0, len(methods))
	)
	for _, multi := range methods {
		var (
			final    *types.Func
			inherits = make([]*types.Func, 0, len(multi))
		)
		for _, m := range multi {
			if !m.Exported() {
				continue
			}
			if DirectedMethod(m, t) {
				final = m
				break
			}
			inherits = append(inherits, m)
		}
		if final != nil {
			if RequirePointerRecv(final, t) {
				finals = append(finals, final)
				continue
			}
			continue
		}
		if len(inherits) == 1 {
			final = inherits[0]
			if _, ok := fields[final.Name()]; !ok {
				finals = append(finals, final)
				continue
			}
			continue
		}
	}
	sort.Slice(finals, func(i, j int) bool {
		return finals[i].Name() < finals[j].Name()
	})

	v := &Methods{methods: make(map[string]int)}
	for i, f := range finals {
		v.methods[f.Name()] = i
		v.ordered = append(v.ordered, f)
	}

	return v
}

type Methods struct {
	methods map[string]int
	ordered []*types.Func
}

func (m *Methods) Method(i int) *types.Func {
	if i >= len(m.ordered) {
		return nil
	}
	return m.ordered[i]
}

func (m *Methods) MethodByName(name string) *types.Func {
	if i, ok := m.methods[name]; ok {
		return m.ordered[i]
	}
	return nil
}

func (m *Methods) NumMethod() int {
	return len(m.ordered)
}

func FieldByName(t types.Type, match func(string) bool) *Field {
	return FieldByNameFunc(t, match, &NamedBacktrace{}, 0)
}

func FieldByNameFunc(t types.Type, match func(string) bool, names *NamedBacktrace, entries int) *Field {
	switch x := t.(type) {
	case *types.Named:
		if names.Push(x) {
			return FieldByNameFunc(x.Underlying(), match, names, entries)
		}
		return nil
	case *types.Struct:
		var (
			direct    *Field
			embeddeds []*Field
		)
		for i := range x.NumFields() {
			f := x.Field(i)
			if match(f.Name()) {
				// potential match or matched multiple times
				if direct != nil {
					return nil
				}
				direct = &Field{f, x.Tag(i)}
			}
			if f.Anonymous() {
				embedded := FieldByNameFunc(f.Type(), match, names, entries+1)
				if embedded != nil {
					embeddeds = append(embeddeds, embedded)
				}
			}
		}
		if direct != nil {
			return direct
		}
		if len(embeddeds) == 1 {
			return embeddeds[0]
		}
		return nil
	case *types.Pointer:
		if entries > 0 {
			return FieldByNameFunc(x.Elem(), match, names, entries)
		}
		return nil
	default:
		return nil
	}
}

type Field struct {
	*types.Var
	Tag string
}
