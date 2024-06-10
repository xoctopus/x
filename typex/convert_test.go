package typex_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"

	"github.com/xoctopus/x/ptrx"
	"github.com/xoctopus/x/typex"
	"github.com/xoctopus/x/typex/testdata"
	_ "github.com/xoctopus/x/typex/testdata/xoxo"
)

func TestTypeByImportAndName(t *testing.T) {
	cases := []struct {
		name string
		path string
		id   string
		want string
	}{
		{
			"NamedTypeXoxoPart",
			"github.com/xoctopus/x/typex/testdata/xoxo", "Part",
			"github.com/xoctopus/x/typex/testdata/xoxo.Part",
		}, {
			"NamedTypeString",
			"github.com/xoctopus/x/typex/testdata", "String",
			"github.com/xoctopus/x/typex/testdata.String",
		}, {
			"NamedTypeInt",
			"github.com/xoctopus/x/typex/testdata", "Int",
			"github.com/xoctopus/x/typex/testdata.Int",
		}, {
			"NamedTypeInt8",
			"github.com/xoctopus/x/typex/testdata", "Int8",
			"github.com/xoctopus/x/typex/testdata.Int8",
		}, {
			"NamedTypeInt16",
			"github.com/xoctopus/x/typex/testdata", "Int16",
			"github.com/xoctopus/x/typex/testdata.Int16",
		}, {
			"NamedTypeInt32",
			"github.com/xoctopus/x/typex/testdata", "Int32",
			"github.com/xoctopus/x/typex/testdata.Int32",
		}, {
			"NamedTypeInt64",
			"github.com/xoctopus/x/typex/testdata", "Int64",
			"github.com/xoctopus/x/typex/testdata.Int64",
		}, {
			"NamedTypeUint",
			"github.com/xoctopus/x/typex/testdata", "Uint",
			"github.com/xoctopus/x/typex/testdata.Uint",
		}, {
			"NamedTypeUint8",
			"github.com/xoctopus/x/typex/testdata", "Uint8",
			"github.com/xoctopus/x/typex/testdata.Uint8",
		}, {
			"NamedTypeUint16",
			"github.com/xoctopus/x/typex/testdata", "Uint16",
			"github.com/xoctopus/x/typex/testdata.Uint16",
		}, {
			"NamedTypeUint32",
			"github.com/xoctopus/x/typex/testdata", "Uint32",
			"github.com/xoctopus/x/typex/testdata.Uint32",
		}, {
			"NamedTypeUint64",
			"github.com/xoctopus/x/typex/testdata", "Uint64",
			"github.com/xoctopus/x/typex/testdata.Uint64",
		}, {
			"NamedTypeUintptr",
			"github.com/xoctopus/x/typex/testdata", "Uintptr",
			"github.com/xoctopus/x/typex/testdata.Uintptr",
		}, {
			"NamedTypeFloat32",
			"github.com/xoctopus/x/typex/testdata", "Float32",
			"github.com/xoctopus/x/typex/testdata.Float32",
		}, {
			"NamedTypeFloat64",
			"github.com/xoctopus/x/typex/testdata", "Float64",
			"github.com/xoctopus/x/typex/testdata.Float64",
		}, {
			"NamedTypeComplex64",
			"github.com/xoctopus/x/typex/testdata", "Complex64",
			"github.com/xoctopus/x/typex/testdata.Complex64",
		}, {
			"NamedTypeComplex128",
			"github.com/xoctopus/x/typex/testdata", "Complex128",
			"github.com/xoctopus/x/typex/testdata.Complex128",
		}, {
			"BasicInt",
			"", "int", "int",
		}, {
			"BasicString",
			"", "string", "string",
		}, {
			"BasicMap",
			"", "map[int]int", "map[int]int",
		}, {
			"BasicIntSlice",
			"", "[]int", "[]int",
		}, {
			"BasicIntArray",
			"", "[2]int", "[2]int",
		}, {
			"BasicError",
			"", "error", "error",
		}, {
			"BasicInterface",
			"", "interface {}", "interface{}",
		}, {
			"BasicAny",
			"", "any", "interface{}",
		}, {
			"BasicComparable",
			"", "comparable", "interface{}",
		}, {
			"GenericTypeAlias",
			"github.com/xoctopus/x/typex/testdata", "EnumMap",
			"github.com/xoctopus/x/typex/testdata.AnyMap[github.com/xoctopus/x/typex/testdata.Enum, any]",
		}, {
			"EmptyInvalid",
			"", "", "invalid type",
		}, {
			"InvalidMapMissingRsb",
			"", "map[int", "invalid type",
		}, {
			"InvalidMapMissingElemType",
			"", "map[int]", "invalid type",
		}, {
			"InvalidMapMissingKeyType",
			"", "map[]any", "invalid type",
		}, {
			"InvalidSliceMissingRsb",
			"", "[any", "invalid type",
		}, {
			"InvalidSliceMissingElem",
			"", "[]", "invalid type",
		}, {
			"InvalidArray",
			"", "[aaa]int", "invalid type",
		}, {
			"InvalidImport",
			"xxx", "any",
			"invalid type",
		}, {
			"InvalidType",
			"github.com/xoctopus/x/typex/testdata", "Invalid",
			"invalid type",
		}, {
			"InvalidTypeFromCache", // from cache
			"github.com/xoctopus/x/typex/testdata", "Invalid",
			"invalid type",
		}, {
			"InvalidImportedPath1",
			"", "abc", "invalid type",
		}, {
			"InvalidGenericTypeMissingRsb",
			"github.com/xoctopus/x/typex/testata", "AnyMap[",
			"invalid type",
		}, {
			"InvalidGenericTypeInvalidImportPath",
			"", "Any[comparable,any]",
			"invalid type",
		}, {
			"InvalidGenericTypeInvalidImportPath",
			"xxx", "Any[comparable,any]",
			"invalid type",
		}, {
			"InvalidGenericTypeInvalidTypename",
			"github.com/xoctopus/x/typex/testdata", "SomeType[comparable,any]",
			"invalid type",
		}, {
			"InvalidGenericTypeInvalidTypeParam",
			"github.com/xoctopus/x/typex/testdata", "AnyMap[Invalid,any]",
			"invalid type",
		}, {
			"InvalidGenericTypeInvalidTypeParamUnmatchedParamCount",
			"github.com/xoctopus/x/typex/testdata", "AnyMap[any]",
			"invalid type",
		}, {
			"GenericTypeWithoutTypeParams",
			"github.com/xoctopus/x/typex/testdata", "AnyMap",
			"github.com/xoctopus/x/typex/testdata.AnyMap[K comparable, V any]",
		}, {
			"GenericType",
			"github.com/xoctopus/x/typex/testdata", "AnyMap[int,string]",
			"github.com/xoctopus/x/typex/testdata.AnyMap[K int, V string]",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gt := typex.TypeByImportAndName(c.path, c.id)
			NewWithT(t).Expect(gt.String()).To(Equal(c.want))
		})
	}
}

func TestNewPackage(t *testing.T) {
	NewWithT(t).Expect(typex.NewPackage("")).To(BeNil())
	// ok
	pkg := typex.NewPackage("github.com/pkg/errors")
	NewWithT(t).Expect(pkg).NotTo(BeNil())
	// gopkg.in/yaml.v2 need import
	NewWithT(t).Expect(typex.NewPackage("gopkg.in/yaml.v2").Name()).To(BeEmpty())
}

func TestNewGoTypeFromReflectType(t *testing.T) {
	cases := []struct {
		name string
		v    reflect.Type
		want string
	}{
		{
			"Int",
			reflect.TypeOf(int(0)),
			"int",
		}, {
			"IntPtr",
			reflect.TypeOf(ptrx.Ptr(int(0))),
			"*int",
		}, {
			"IntPtrPtr",
			reflect.TypeOf(ptrx.Ptr(int(0))),
			"*int",
		}, {
			"PkgError",
			reflect.TypeOf(pkgerrors.New("any")),
			"*github.com/pkg/errors.fundamental",
		}, {
			"Error",
			reflect.TypeOf(errors.New("any")),
			"*errors.errorString",
		}, {
			"NilError",
			reflect.TypeOf(error(nil)),
			"nil",
		}, {
			"IntSlice",
			reflect.TypeOf([]int{}),
			"[]int",
		}, {
			"IntArray",
			reflect.TypeOf([3]int{}),
			"[3]int",
		}, {
			"ErrorSlice",
			reflect.TypeOf([]error{}),
			"[]error",
		}, {
			"ErrorArray",
			reflect.TypeOf([5]error{}),
			"[5]error",
		}, {
			"IntChanSendOnly",
			reflect.TypeOf(make(<-chan int)),
			"<-chan int",
		}, {
			"IntChanRecvOnly",
			reflect.TypeOf(make(chan<- int)),
			"chan<- int",
		}, {
			"IntChanBothDir",
			reflect.TypeOf(make(chan int)),
			"chan int",
		}, {
			"Map",
			reflect.TypeOf(map[int]string{}),
			"map[int]string",
		}, {
			"Func",
			reflect.TypeOf(func() {}),
			"func()",
		}, {
			"FuncWithInAndOut",
			reflect.TypeOf(func(a string, b reflect.Type, c testdata.Interface) (err error, res testdata.Interface) {
				return nil, nil
			}),
			"func(string, reflect.Type, github.com/xoctopus/x/typex/testdata.Interface) (error, github.com/xoctopus/x/typex/testdata.Interface)",
		}, {
			"MapInterface",
			reflect.TypeOf(map[int]testdata.Interface{}),
			"map[int]github.com/xoctopus/x/typex/testdata.Interface",
		}, {
			"MapStruct",
			reflect.TypeOf(map[int]*testdata.Struct{}),
			"map[int]*github.com/xoctopus/x/typex/testdata.Struct",
		}, {
			"StructPtr",
			reflect.TypeOf(&testdata.Struct{}),
			"*github.com/xoctopus/x/typex/testdata.Struct",
		}, {
			"StructSlice",
			reflect.TypeOf([]struct {
				A int
				B string
			}{}),
			"[]struct{A int; B string}",
		}, {
			"InterfaceSlicePtr",
			reflect.TypeOf(&[]interface {
				String() string
				Bytes() []byte
			}{}),
			"*[]interface{Bytes() []uint8; String() string}",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tt := typex.NewGoTypeFromReflectType(c.v)
			if tt == nil {
				NewWithT(t).Expect("nil").To(Equal(c.want))
			} else {
				NewWithT(t).Expect(tt.String()).To(Equal(c.want))
			}
		})
	}
}
