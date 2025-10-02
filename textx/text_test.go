package textx_test

import (
	"bytes"
	"errors"
	"math/big"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/textx"
	"github.com/xoctopus/x/textx/testdata"
)

var (
	AsErrMarshalUnsupportedType   *ErrMarshalUnsupportedType
	AsErrMarshalFailed            *ErrMarshalFailed
	AsErrUnmarshalInvalidInput    *ErrUnmarshalInvalidInput
	AsErrUnmarshalUnsupportedType *ErrUnmarshalUnsupportedType
	AsErrUnmarshalParseFailed     *ErrUnmarshalParseFailed
	AsErrUnmarshalFailed          *ErrUnmarshalFailed
)

func TestMarshalText(t *testing.T) {
	cases := map[string]struct {
		input  any
		result any
	}{
		"Invalid":        {nil, []byte(nil)},
		"NamedString":    {testdata.String("any"), []byte("any")},
		"NamedStringPtr": {ptrx.Ptr(testdata.String("any")), []byte("any")},
		"Int":            {1, []byte("1")},
		"IntPtr":         {new(int), []byte("0")},
		"Uint":           {uint32(1), []byte("1")},
		"Float32":        {float32(1.1), []byte("1.1")},
		"Float64":        {1.1, []byte("1.1")},
		"Boolean":        {true, []byte("true")},
		"Bytes1":         {[]byte("xyz"), []byte("xyz")},
		"Bytes2":         {testdata.Bytes{'x', 'y', 'z'}, []byte("xyz")},
		"Arshaler1":      {testdata.Duration(100), []byte("100")},
		"Arshaler2": {
			testdata.Integers{0x0100000000000000, 0x00000000000000FF},
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		"Arshaler3": {
			reflect.ValueOf(&struct{ *big.Float }{big.NewFloat(100.001)}).Elem().Field(0),
			[]byte("100.001"),
		},
		"UnsupportedType":   {[]int{1, 2, 3}, AsErrMarshalUnsupportedType},
		"MarshalTextFailed": {testdata.MustFailedArshaler{}, AsErrMarshalFailed},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			text, err := Marshal(c.input)
			if err != nil {
				NewWithT(t).Expect(errors.As(err, &c.result)).To(BeTrue())
			} else {
				NewWithT(t).Expect(bytes.Equal(text, c.result.([]byte))).To(BeTrue())
			}
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	cases := map[string]struct {
		input  []byte
		value  any
		result any
	}{
		"InvalidInput1":       {nil, nil, AsErrUnmarshalInvalidInput},
		"InvalidInput2":       {nil, 1, AsErrUnmarshalInvalidInput},
		"IntFailed":           {[]byte("invalid"), new(int), AsErrUnmarshalParseFailed},
		"UintFailed":          {[]byte("invalid"), new(uint), AsErrUnmarshalParseFailed},
		"FloatFailed":         {[]byte("invalid"), new(float32), AsErrUnmarshalParseFailed},
		"BooleanFailed":       {[]byte("invalid"), new(bool), AsErrUnmarshalParseFailed},
		"Unsupported":         {[]byte("invalid"), new([]string), AsErrUnmarshalUnsupportedType},
		"MustUnmarshalFailed": {nil, new(testdata.MustFailedArshaler), AsErrUnmarshalFailed},
		"ArshalerFailed":      {[]byte{1, 2, 3, 4}, new(testdata.Integers), AsErrUnmarshalFailed},
		"Int":                 {[]byte("1"), new(int), ptrx.Ptr(1)},
		"Int2":                {[]byte("1"), new(*int), ptrx.Ptr(ptrx.Ptr(1))},
		"Uint":                {[]byte("1"), new(uint), ptrx.Ptr(uint(1))},
		"Float":               {[]byte("1.1"), new(float64), ptrx.Ptr(1.1)},
		"Boolean":             {[]byte("true"), new(bool), ptrx.Ptr(true)},
		"String":              {[]byte("any"), new(string), ptrx.Ptr("any")},
		"String2":             {[]byte("any"), new(testdata.String), ptrx.Ptr(testdata.String("any"))},
		"Bytes":               {[]byte("any"), new([]byte), ptrx.Ptr([]byte("any"))},
		"Bytes2":              {[]byte("any"), new(testdata.Bytes), ptrx.Ptr(testdata.Bytes("any"))},
		"Arshaler1":           {[]byte("1"), new(testdata.Duration), ptrx.Ptr(testdata.Duration(1))},
		"Arshaler2":           {[]byte{1, 0, 0, 0, 0, 0, 0, 0}, new(testdata.Integers), ptrx.Ptr(testdata.Integers{1})},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := Unmarshal(c.input, c.value)
			if err != nil {
				NewWithT(t).Expect(err).NotTo(BeNil())
				NewWithT(t).Expect(errors.As(err, &c.result)).To(BeTrue())
			} else {
				NewWithT(t).Expect(c.value).To(Equal(c.result))
			}
		})
	}
}

// func Benchmark_ParseFloatToBytes(b *testing.B) {
// 	b.Run("append", func(b *testing.B) {
// 		size := unsafe.Sizeof(1.1)
// 		data := make([]byte, 0, size)
// 		for i := 0; i < b.N; i++ {
// 			data = data[:]
// 			strconv.AppendFloat(data, 1.1, 'f', -1, 32)
// 		}
// 	})
//
// 	b.Run("format", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			_ = []byte(strconv.FormatFloat(1.1, 'f', -1, 32))
// 		}
// 	})
// }
// func BenchmarkUnmarshalTextAndMarshalText(b *testing.B) {
// 	for _, c := range unmarshalCases {
// 		if c.err != nil {
// 			continue
// 		}
// 		b.Run("UnmarshalText_"+c.name, func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				_ = Unmarshal(c.text, c.value)
// 			}
// 		})
// 	}
// }

// func BenchmarkMarshalText(b *testing.B) {
// 	for name, c := range marshalCases {
// 		b.Run("MarshalText_"+name, func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				_, _ = Marshal(c.value)
// 			}
// 		})
// 	}
// }
