package internal_test

import (
	"go/types"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/typex/internal"
	"github.com/xoctopus/x/typex/testdata"
)

func TestInstTypeID(t *testing.T) {
	pkg := NewPackage("github.com/xoctopus/x/typex/testdata")

	serialized := pkg.Scope().Lookup("Serialized").Type()
	serializedstring, err := types.Instantiate(nil, serialized, []types.Type{types.Typ[types.String]}, true)
	NewWithT(t).Expect(err).To(BeNil())
	serializedbytes, err := types.Instantiate(nil, serialized, []types.Type{types.NewSlice(types.Typ[types.Byte])}, true)
	NewWithT(t).Expect(err).To(BeNil())

	cases := []struct {
		name        string
		typ         types.Type
		args        []types.Type
		typenames   [2]string
		underlyings [2]string
		methods     map[string][2]string
	}{
		{
			name: "TypeParamPass__string_net_Addr",
			typ:  pkg.Scope().Lookup("TypeParamPass").Type().(*types.Named),
			args: []types.Type{
				types.Typ[types.String],
				NewPackage("net").Scope().Lookup("Addr").Type(),
			},
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.TypeParamPass[string,net.Addr]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypeParamPass[string,net.Addr]",
			},
			underlyings: [2]string{
				"struct { t1 string; t2 net.Addr; *github.com/xoctopus/x/typex/testdata.BTreeNode[net.Addr] }",
				"struct { t1 string; t2 net.Addr; *xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.BTreeNode[net.Addr] }",
			},
			methods: map[string][2]string{
				"Deal": {"func(string) net.Addr", "func(string) net.Addr"},
			},
		},
		{
			name: "Serialized_string",
			typ:  serialized.(*types.Named),
			args: []types.Type{types.Typ[types.String]},
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.Serialized[string]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]",
			},
			underlyings: [2]string{"struct { data string }", "struct { data string }"},
			methods: map[string][2]string{
				"String": {"func() string", "func() string"},
				"Bytes":  {"func() []uint8", "func() []uint8"},
				"Data":   {"func() string", "func() string"},
			},
		},
		{
			name: "Serialized_bytes",
			typ:  serialized.(*types.Named),
			args: []types.Type{
				types.NewSlice(types.Typ[types.Byte]),
			},
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]",
			},
			underlyings: [2]string{"struct { data []uint8 }", "struct { data []uint8 }"},
			methods: map[string][2]string{
				"String": {"func() string", "func() string"},
				"Bytes":  {"func() []uint8", "func() []uint8"},
				"Data":   {"func() []uint8", "func() []uint8"},
			},
		},
		{
			name: "TypeParamPass__SerializedString_SerializedBytes",
			typ:  pkg.Scope().Lookup("TypeParamPass").Type().(*types.Named),
			args: []types.Type{serializedstring, serializedbytes},
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.TypeParamPass[github.com/xoctopus/x/typex/testdata.Serialized[string],github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypeParamPass[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string],xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]]",
			},
			underlyings: [2]string{
				"struct { t1 github.com/xoctopus/x/typex/testdata.Serialized[string]; t2 github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]; *github.com/xoctopus/x/typex/testdata.BTreeNode[github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]] }",
				"struct { t1 xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]; t2 xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]; *xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.BTreeNode[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]] }",
			},
			methods: map[string][2]string{
				"Deal": {
					"func(github.com/xoctopus/x/typex/testdata.Serialized[string]) github.com/xoctopus/x/typex/testdata.Serialized[[]uint8]",
					"func(xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[string]) xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Serialized[[]uint8]",
				},
			},
		},
		{
			name: "MapBTreeNode__string",
			typ:  pkg.Scope().Lookup("MapBTreeNode").Type().(*types.Named),
			args: []types.Type{types.Typ[types.String]},
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.MapBTreeNode[string]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.MapBTreeNode[string]",
			},
			underlyings: [2]string{
				"map[string]*github.com/xoctopus/x/typex/testdata.BTreeNode[string]",
				"map[string]*xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.BTreeNode[string]",
			},
		},
		{
			name: "TypedSliceAlias",
			typ:  pkg.Scope().Lookup("TypedSliceAlias").Type().(*types.Alias),
			args: nil,
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.TypedSlice[net.Addr]",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedSlice[net.Addr]",
			},
			underlyings: [2]string{"[]net.Addr", "[]net.Addr"},
			methods:     map[string][2]string{"Len": {"func() int", "func() int"}},
		},
		{
			name: "SimpleStruct",
			typ:  pkg.Scope().Lookup("SimpleStruct").Type().(*types.Named),
			args: nil,
			typenames: [2]string{
				"github.com/xoctopus/x/typex/testdata.SimpleStruct",
				"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.SimpleStruct",
			},
			underlyings: [2]string{
				"struct { A string; B github.com/xoctopus/x/typex/testdata.String; unexported interface {}; Name fmt.Stringer; HasTag github.com/xoctopus/x/typex/testdata.Int tag:\"tagKey,otherLabel\"; github.com/xoctopus/x/typex/testdata.EmptyInterface }",
				"struct { A string; B xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.String; unexported interface {}; Name fmt.Stringer; HasTag xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.Int tag:\"tagKey,otherLabel\"; xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.EmptyInterface }",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			typenames := [2]string{
				InstTypeID(c.typ, c.args),
				InstTypeID(c.typ, c.args, true),
			}
			NewWithT(t).Expect(typenames).To(Equal(c.typenames))

			underlyings := [2]string{
				InstTypeID(c.typ.Underlying(), c.args),
				InstTypeID(c.typ.Underlying(), c.args, true),
			}
			NewWithT(t).Expect(underlyings).To(Equal(c.underlyings))

			if len(c.args) > 0 {
				inst, _ := types.Instantiate(nil, c.typ, c.args, false)
				NewWithT(t).Expect(inst).NotTo(BeNil())
				underlyings = [2]string{
					TypesTypeID(InstantiateUnderlying(inst)),
					TypesTypeID(InstantiateUnderlying(inst), true),
				}
				NewWithT(t).Expect(underlyings).To(Equal(c.underlyings))
			}

			named, ok := c.typ.(*types.Named)
			if !ok {
				named = c.typ.(*types.Alias).Rhs().(*types.Named)
			}

			methods := make(map[string][2]string)
			for i := range named.NumMethods() {
				m := named.Method(i)
				methods[m.Name()] = [2]string{
					InstTypeID(m.Signature(), c.args),
					InstTypeID(m.Signature(), c.args, true),
				}
			}
			NewWithT(t).Expect(len(methods)).To(Equal(len(c.methods)))
			for name := range methods {
				names, ok := c.methods[name]
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(names).To(Equal(methods[name]))
			}
		})
	}

	t.Run("InstTypesTypeID", func(t *testing.T) {
		t.Run("InvalidType", func(t *testing.T) {
			defer func() {
				NewWithT(t).Expect(recover()).NotTo(BeNil())
			}()
			InstTypeID(types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.Int])), nil)
		})
	})

	t.Run("NewTypesTypeID", func(t *testing.T) {
		t.Run("WithEmpty", func(t *testing.T) {
			NewWithT(t).Expect(types.Identical(NewTypesTypeByID(""), types.Typ[types.Invalid])).To(BeTrue())
		})
	})
}

func TestInstantiate(t *testing.T) {
	t.Run("InstantiateUnderlying", func(t *testing.T) {
		t.Run("Instantiated", func(t *testing.T) {
			gt := NewTypesType(reflect.TypeOf(testdata.T[int]{}))
			InstantiateUnderlying(gt)
		})
		t.Run("Uninstantiated", func(t *testing.T) {
			gt := NewPackage("github.com/xoctopus/x/typex/testdata").Scope().Lookup("T").Type().Underlying()
			Instantiate(gt.Underlying(), []types.Type{types.Typ[types.Int]})
		})
	})
	t.Run("NoTypeArgs", func(t *testing.T) {
		InstantiateUnderlying(types.Typ[types.Int])
		Instantiate(types.Typ[types.Int], nil)
	})
	t.Run("InvalidType", func(t *testing.T) {
		defer func() {
			NewWithT(t).Expect(recover()).NotTo(BeNil())
		}()
		Instantiate(types.NewUnion([]*types.Term{
			types.NewTerm(false, types.Typ[types.String]),
			types.NewTerm(false, types.NewSlice(types.Typ[types.Byte])),
		}), []types.Type{types.Typ[types.Int]})
	})
}
