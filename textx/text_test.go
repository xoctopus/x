package textx_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/textx"
)

type Duration time.Duration

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	dur, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

var errInvalidDataLength = errors.New("data length should times 8")

type UnsignedIntegers []uint64

func (d UnsignedIntegers) UnmarshalText(data []byte) error {
	if len(data)%8 != 0 {
		return errInvalidDataLength
	}
	for offset := 0; len(data) >= 8; offset += 8 {
		v := binary.LittleEndian.Uint64(data)
		data = data[8:]
		d[offset/8] = v
	}
	return nil
}

type NamedString string
type NamedInt int

func Benchmark_ParseFloatToBytes(b *testing.B) {
	b.Run("append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.AppendFloat([]byte{}, 1.1, 'f', -1, 32)
		}
	})

	b.Run("format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = []byte(strconv.FormatFloat(1.1, 'f', -1, 32))
		}
	})
}

type MarshalCase struct {
	value any
	text  []byte
	err   error
}

var marshalCases = map[string]*MarshalCase{
	"NamedString":     {NamedString("any"), []byte("any"), nil},
	"NamedStringPtr":  {ptrx.Ptr(NamedString("any")), []byte("any"), nil},
	"NamedInt":        {NamedInt(1), []byte("1"), nil},
	"NamedIntPtr":     {ptrx.Ptr(NamedInt(1)), []byte("1"), nil},
	"Invalid":         {nil, nil, nil},
	"Integer":         {1, []byte("1"), nil},
	"UnsignedInteger": {uint(1), []byte("1"), nil},
	"Boolean":         {true, []byte("true"), nil},
	"Float32":         {float32(1.1), []byte("1.1"), nil},
	"Float64":         {1.1, []byte("1.1"), nil},
	"Bytes":           {[]byte("any"), ToBase64([]byte("any")), nil},
	"Marshaller":      {Duration(time.Second), []byte("1s"), nil},
	"Unsupported":     {struct{}{}, nil, &ErrMarshalUnsupportedType{Type: reflect.TypeOf(struct{}{})}},
}

func TestMarshalText(t *testing.T) {
	for name, c := range marshalCases {
		t.Run(name, func(t *testing.T) {
			text, err := MarshalText(c.value)
			if c.err != nil {
				NewWithT(t).Expect(err.Error()).To(Equal(c.err.Error()))
				return
			}
			NewWithT(t).Expect(bytes.Equal(text, c.text)).To(BeTrue())
		})
	}
}

type UnmarshalCase struct {
	name  string
	value any
	text  []byte
	err   error
}

var unmarshalCases = []*UnmarshalCase{
	{"Nil", nil, []byte("any"), &ErrInvalidUnmarshal{Type: nil, Err: "invalid value"}},
	{"CannotSet1", 1, []byte("1"), &ErrInvalidUnmarshal{Type: reflect.TypeOf(1), Err: "cannot set"}},
	{"CannotSet2", (*int)(nil), []byte("1"), &ErrInvalidUnmarshal{Type: nil, Err: "invalid value"}},
	{"CannotSet3", reflect.ValueOf(struct{ Int int }{}).Field(0), []byte("1"), &ErrInvalidUnmarshal{Type: reflect.TypeOf(0), Err: "cannot set"}},
	{"CannotSet4", reflect.ValueOf(&struct{ value int }{}).Elem().Field(0), []byte("1"), &ErrInvalidUnmarshal{Type: reflect.TypeOf(0), Err: "cannot set"}},
	{"StructField", reflect.ValueOf(&struct{ Int int }{}).Elem().Field(0), []byte("1"), nil},
	{"NamedString", new(NamedString), []byte("any"), nil},
	{"NamedInt", new(NamedInt), []byte("1"), nil},
	{"Integer", new(int), []byte("1"), nil},
	{"IntegerPtr", new(*int), []byte("1"), nil},
	{"IntegerPtrPtr", new(**int), []byte("1"), nil},
	{"IntegerFailed", new(int), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf(1)}},
	{"UnsignedInteger", new(uint32), []byte("1"), nil},
	{"UnsignedIntegerFailed", new(uint32), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf(uint32(1))}},
	{"Float", new(float64), []byte("1"), nil},
	{"FloatFailed", new(float64), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf(float64(1))}},
	{"Boolean", new(bool), []byte("true"), nil},
	{"BooleanFailed", new(bool), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf(true)}},
	{"Bytes", new([]byte), ToBase64([]byte("any")), nil},
	{"BytesFailed", new([]byte), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf([]byte{})}},
	{"Unmarshaler1", new(Duration), []byte("1s"), nil},
	{"UnmarshalerFailed1", new(Duration), []byte("any"), &ErrUnmarshalFailed{Data: []byte("any"), Type: reflect.TypeOf(Duration(0))}},
	{"Unmarshaler2", &UnsignedIntegers{0x02}, []byte{1, 0, 0, 0, 0, 0, 0, 0}, nil},
	{"UnmarshalerFailed2", &UnsignedIntegers{0x02}, []byte{1, 0, 0, 0, 0, 0, 0}, errInvalidDataLength},
	{"Unsupported1", new([]int), []byte("any"), &ErrUnmarshalUnsupportedType{Type: reflect.TypeOf([]int{})}},
	{"Unsupported2", &struct{}{}, []byte("any"), &ErrUnmarshalUnsupportedType{Type: reflect.TypeOf(struct{}{})}},
	{"HexInt", new(int), []byte("0xFF"), nil},
	{"OctUint", new(uint), []byte("077"), nil},
	{"InvalidNumeric", new(float64), []byte("abcdef"), &ErrUnmarshalFailed{Data: []byte("abcdef"), Type: reflect.TypeOf(float64(0))}},
}

func TestUnmarshalText(t *testing.T) {
	for _, c := range unmarshalCases {
		t.Run(c.name, func(t *testing.T) {
			err := UnmarshalText(c.text, c.value)
			if c.err != nil {
				NewWithT(t).Expect(err).NotTo(BeNil())
				NewWithT(t).Expect(err.Error()).To(ContainSubstring(c.err.Error()))
				return
			}
			NewWithT(t).Expect(err).To(BeNil())
		})
	}
}

func BenchmarkUnmarshalTextAndMarshalText(b *testing.B) {
	for _, c := range unmarshalCases {
		if c.err != nil {
			continue
		}
		b.Run("UnmarshalText_"+c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = UnmarshalText(c.text, c.value)
			}
		})
	}
}

func BenchmarkMarshalText(b *testing.B) {
	for name, c := range marshalCases {
		b.Run("MarshalText_"+name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = MarshalText(c.value)
			}
		})
	}
}
