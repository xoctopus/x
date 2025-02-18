package universe_test

import (
	"io"
	"iter"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/typex/testdata"
	"github.com/xoctopus/x/universe"
)

type MockInvalidType struct {
	reflect.Type
}

func (*MockInvalidType) Name() string       { return "" }
func (*MockInvalidType) PkgPath() string    { return "" }
func (*MockInvalidType) Kind() reflect.Kind { return reflect.Invalid }
func (*MockInvalidType) String() string     { return "" }

var cases = []struct {
	t       reflect.Type
	wrapped string
	id      string
}{
	{},
	{reflect.TypeFor[bool](), "bool", "bool"},
	{reflect.TypeFor[rune](), "int32", "int32"},
	{reflect.TypeFor[byte](), "uint8", "uint8"},
	{reflect.TypeFor[error](), "error", "error"},
	{
		reflect.TypeFor[testdata.SimpleStruct](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.SimpleStruct",
		"github.com/xoctopus/x/typex/testdata.SimpleStruct",
	},
	{reflect.TypeFor[iter.Seq[int]](), "iter.Seq[int]", "iter.Seq[int]"},
	{reflect.TypeFor[[3]int](), "[3]int", "[3]int"},
	{reflect.TypeFor[chan error](), "chan error", "chan error"},
	{reflect.TypeFor[chan<- error](), "chan<- error", "chan<- error"},
	{reflect.TypeFor[<-chan error](), "<-chan error", "<-chan error"},
	{reflect.TypeFor[func()](), "func()", "func()"},
	{
		reflect.TypeFor[func(int, ...string)](),
		"func(int, ...string)",
		"func(int, ...string)",
	},
	{
		reflect.TypeFor[func(int, ...string) error](),
		"func(int, ...string) error",
		"func(int, ...string) error",
	},
	{
		reflect.TypeFor[func(int, ...string) (error, bool)](),
		"func(int, ...string) (error, bool)",
		"func(int, ...string) (error, bool)",
	},
	{reflect.TypeFor[any](), "interface {}", "interface {}"},
	{
		reflect.TypeFor[interface{ io.ReadCloser }](),
		"interface { Close() error; Read([]uint8) (int, error) }",
		"interface { Close() error; Read([]uint8) (int, error) }",
	},
	{reflect.TypeFor[map[string]int](), "map[string]int", "map[string]int"},
	{reflect.TypeFor[*int](), "*int", "*int"},
	{reflect.TypeFor[[]int](), "[]int", "[]int"},
	{reflect.TypeFor[struct{}](), "struct {}", "struct {}"},
	{
		reflect.TypeFor[struct {
			A string `json:",;[](){}\"''"`
			testdata.SimpleStruct
			testdata.TypedMap[int, any]
		}](),
		`struct { A string "json:\",;[](){}\\\"''\""; xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.SimpleStruct; xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedMap[int,interface {}] }`,
		`struct { A string "json:\",;[](){}\\\"''\""; github.com/xoctopus/x/typex/testdata.SimpleStruct; github.com/xoctopus/x/typex/testdata.TypedMap[int,interface {}] }`,
	},
	{&MockInvalidType{}, "", ""},
	{
		reflect.TypeFor[testdata.TypedArray[bool]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[bool]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[bool]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[rune]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[int32]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[int32]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[byte]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[uint8]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[uint8]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[error]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[error]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[error]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[testdata.SimpleStruct]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.SimpleStruct]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[github.com/xoctopus/x/typex/testdata.SimpleStruct]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[iter.Seq[int]]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[iter.Seq[int]]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[iter.Seq[int]]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[[3]int]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[[3]int]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[[3]int]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[chan error]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[chan error]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[chan error]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[chan<- error]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[chan<- error]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[chan<- error]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[<-chan error]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[<-chan error]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[<-chan error]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[func()]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[func()]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[func()]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[func(int, ...string)]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[func(int, ...string)]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[func(int, ...string)]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[func(int, ...string) error]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[func(int, ...string) error]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[func(int, ...string) error]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[func(int, ...string) (error, bool)]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[func(int, ...string) (error, bool)]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[func(int, ...string) (error, bool)]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[any]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[interface {}]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[interface {}]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[interface{ io.ReadCloser }]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[interface { Close() error; Read([]uint8) (int, error) }]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[interface { Close() error; Read([]uint8) (int, error) }]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[map[string]int]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[map[string]int]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[map[string]int]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[*int]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[*int]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[*int]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[[]int]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[[]int]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[[]int]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[struct{}]](),
		"xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[struct {}]",
		"github.com/xoctopus/x/typex/testdata.TypedArray[struct {}]",
	},
	{
		reflect.TypeFor[testdata.TypedArray[struct {
			A string `json:",;[](){}\"''"`
			testdata.SimpleStruct
			testdata.TypedMap[int, any]
			C struct {
				testdata.TypedArray[struct{ testdata.TypedInterface[any] }]
			}
		}]](),
		`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[struct { ` +
			`A string "json:\",;[](){}\\\"''\""; ` +
			`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.SimpleStruct; ` +
			`xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedMap[int,interface {}]; ` +
			`C struct { xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedArray[struct { xoctopus__github__com_xoctopus_x_typex_testdata__xoctopus.TypedInterface[interface {}] }] } ` +
			`}]`,
		`github.com/xoctopus/x/typex/testdata.TypedArray[struct { ` +
			`A string "json:\",;[](){}\\\"''\""; ` +
			`github.com/xoctopus/x/typex/testdata.SimpleStruct; ` +
			`github.com/xoctopus/x/typex/testdata.TypedMap[int,interface {}]; ` +
			`C struct { github.com/xoctopus/x/typex/testdata.TypedArray[struct { github.com/xoctopus/x/typex/testdata.TypedInterface[interface {}] }] } ` +
			`}]`,
	},
}

func TestNewType(t *testing.T) {
	for _, c := range cases {
		wrapped := universe.WrapIDByT(c.t)
		NewWithT(t).Expect(wrapped).To(Equal(c.wrapped))
		tt := universe.NewUniverse(c.t).UnwrapOr(nil)
		if tt != nil {
			NewWithT(t).Expect(tt.String()).To(Equal(c.id))
		}
	}
	NewWithT(t).Expect((&universe.Named{}).UniverseKind()).To(Equal(universe.TypeName))
	NewWithT(t).Expect((&universe.Unnamed{}).UniverseKind()).To(Equal(universe.TypeLit))

	t.Run("InvalidUnnamedKind", func(t *testing.T) {
		defer func() {
			NewWithT(t).Expect(recover().(error).Error()).To(ContainSubstring("unexpected Unnamed kind"))
		}()
		_ = (&universe.Unnamed{}).String()
	})

	t.Run("InvalidID", func(t *testing.T) {
		defer func() {
			NewWithT(t).Expect(recover().(error).Error()).To(ContainSubstring("invalid id"))
		}()
		universe.NewUniverseByID("int(100.1)").Unwrap()
	})
}
