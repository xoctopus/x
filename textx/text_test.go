package textx_test

import (
	"math/big"
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/testx/bdd"
	. "github.com/xoctopus/x/textx"
	"github.com/xoctopus/x/textx/testdata"
)

func TestTextArshaler(t *testing.T) {
	integersArshaler := testdata.Integers{0x0100000000000000, 0x00000000000000FF}
	integersArshalerResult := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	pointerValueArshaler := reflect.ValueOf(&struct{ *big.Float }{big.NewFloat(100.001)}).Elem().Field(0)
	pointerValueArshalerResult := []byte("100.001")

	cases1 := map[string]struct {
		input  any
		result any
	}{
		"invalid input":        {nil, []byte(nil)},
		"named string":         {testdata.String("any"), []byte("any")},
		"named string pointer": {new(testdata.String("any")), []byte("any")},
		"integer":              {1, []byte("1")},
		"integer pointer":      {new(int), []byte("0")},
		"unsigned integer":     {uint32(1), []byte("1")},
		"float32":              {float32(1.1), []byte("1.1")},
		"float64":              {1.1, []byte("1.1")},
		"boolean":              {true, []byte("true")},
		"bytes":                {[]byte("xyz"), []byte("xyz")},
		"bytes underlying":     {testdata.Bytes{'x', 'y', 'z'}, []byte("xyz")},
		"Arshaler1":            {testdata.Duration(100), []byte("100")},
		"Arshaler2":            {integersArshaler, integersArshalerResult},
		"Arshaler3":            {pointerValueArshaler, pointerValueArshalerResult},
		"unsupported input":    {[]int{1, 2, 3}, codex.New(ECODE__MARSHAL_TEXT_INVALID_INPUT)},
		"must failed arshaler": {testdata.MustFailedArshaler{}, codex.New(ECODE__MARSHAL_TEXT_FAILED)},
	}

	bdd.From(t).When("marshaling text", func(t bdd.T) {
		for summary, c := range cases1 {
			text, err := Marshal(c.input)

			t.Given(summary, func(t bdd.T) {
				if err == nil {
					t.Then("should equal expect", bdd.Equal(c.result.([]byte), text))
				} else {
					t.Then("should match error", bdd.IsError(c.result.(error), err))
				}
			})
		}
	})

	cases2 := map[string]struct {
		input  []byte
		value  any
		result any
	}{
		"nil input":                        {nil, nil, codex.New(ECODE__UNMARSHAL_TEXT_INVALID_INPUT)},
		"InvalidInput2":                    {nil, 1, codex.New(ECODE__UNMARSHAL_TEXT_INVALID_INPUT)},
		"invalid integer input":            {[]byte("invalid"), new(int), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"invalid unsigned integer input":   {[]byte("invalid"), new(uint), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"invalid float32 input":            {[]byte("invalid"), new(float32), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"invalid boolean input":            {[]byte("invalid"), new(bool), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"unsupported type input([]string)": {[]byte("invalid"), new([]string), codex.New(ECODE__UNMARSHAL_TEXT_INVALID_INPUT)},
		"must failed arshaler":             {nil, new(testdata.MustFailedArshaler), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"another must failed arshaler":     {[]byte{1, 2, 3, 4}, new(testdata.Integers), codex.New(ECODE__UNMARSHAL_TEXT_FAILED)},
		"integer":                          {[]byte("1"), new(int), new(1)},
		"integer pointer":                  {[]byte("1"), new(*int), new(new(1))},
		"unsigned":                         {[]byte("1"), new(uint), new(uint(1))},
		"float64":                          {[]byte("1.1"), new(float64), new(1.1)},
		"boolean":                          {[]byte("true"), new(bool), new(true)},
		"string":                           {[]byte("any"), new(string), new("any")},
		"string underlying":                {[]byte("any"), new(testdata.String), new(testdata.String("any"))},
		"bytes":                            {[]byte("any"), new([]byte), new([]byte("any"))},
		"bytes underlying":                 {[]byte("any"), new(testdata.Bytes), new(testdata.Bytes("any"))},
		"arshaler pointer":                 {[]byte("1"), new(testdata.Duration), new(testdata.Duration(1))},
		"another arshaler pointer":         {[]byte{1, 0, 0, 0, 0, 0, 0, 0}, new(testdata.Integers), new(testdata.Integers{1})},
	}

	bdd.From(t).When("unmarshaling text", func(t bdd.T) {
		for summary, c := range cases2 {
			err := Unmarshal(c.input, c.value)
			t.Given(summary, func(t bdd.T) {
				if err == nil {
					t.Then("should equal expect", bdd.Equal(c.result, c.value))
				} else {
					t.Then("should match error", bdd.IsError(c.result.(error), err))
				}
			})
		}
	})
}

func Benchmark_ParseFloatToBytes(b *testing.B) {
	b.Run("append", func(b *testing.B) {
		size := unsafe.Sizeof(1.1)
		data := make([]byte, 0, size)
		for i := 0; i < b.N; i++ {
			data = data[:]
			_ = strconv.AppendFloat(data, 1.1, 'f', -1, 32)
		}
	})

	b.Run("format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = []byte(strconv.FormatFloat(1.1, 'f', -1, 32))
		}
	})
}

//func BenchmarkUnmarshalTextAndMarshalText(b *testing.B) {
//	for _, c := range unmarshalCases {
//		if c.err != nil {
//			continue
//		}
//		b.Run("UnmarshalText_"+c.name, func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				_ = Unmarshal(c.text, c.value)
//			}
//		})
//	}
//}
//
//func BenchmarkMarshalText(b *testing.B) {
//	for name, c := range marshalCases {
//		b.Run("MarshalText_"+name, func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				_, _ = Marshal(c.value)
//			}
//		})
//	}
//}
