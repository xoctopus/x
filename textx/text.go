package textx

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"

	"github.com/xoctopus/x/reflectx"
)

func MarshalText(v any) ([]byte, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	rv = reflectx.Indirect(rv)
	if rv == reflectx.InvalidValue {
		return nil, nil
	}

	if rv.CanInterface() {
		if marshaller, ok := rv.Interface().(encoding.TextMarshaler); ok {
			return marshaller.MarshalText()
		}
	}

	switch kind := rv.Kind(); kind {
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt([]byte{}, rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint([]byte{}, rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.AppendFloat([]byte{}, rv.Float(), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.AppendFloat([]byte{}, rv.Float(), 'g', -1, 64), nil
	case reflect.Bool:
		return strconv.AppendBool([]byte{}, rv.Bool()), nil
	default:
		if reflectx.IsBytes(rv) {
			return ToBase64(rv.Bytes()), nil
		}
		return nil, &ErrMarshalUnsupportedType{rv.Type()}
	}
}

func UnmarshalText(data []byte, v any) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if !rv.IsValid() {
		return &ErrInvalidUnmarshal{nil, "invalid value"}
	}

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() && rv.CanSet() {
			rv.Set(reflectx.New(rv.Type()))
		}
		return UnmarshalText(data, rv.Elem())
	}

	if !rv.CanSet() {
		return &ErrInvalidUnmarshal{rv.Type(), "cannot set"}
	}

	rt := rv.Type()
	if rt.Implements(rtTextUnmarshaller) {
		if err := rv.Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil {
			return &ErrUnmarshalFailed{data, rt, err.Error()}
		}
		return nil
	} else if reflect.PointerTo(rt).Implements(rtTextUnmarshaller) {
		if err := rv.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil {
			return &ErrUnmarshalFailed{data, rt, err.Error()}
		}
		return nil
	}

	switch rv.Kind() {
	case reflect.Slice:
		if reflectx.IsBytes(rv.Type()) {
			raw, err := FromBase64(data)
			if err != nil {
				return &ErrUnmarshalFailed{data, rt, err.Error()}
			}
			rv.SetBytes(raw)
			return nil
		}
		return &ErrUnmarshalUnsupportedType{rt}
	case reflect.String:
		rv.SetString(string(data))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := int64(0)
		if _, err := fmt.Sscan(string(data), &i64); err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui64 := uint64(0)
		if _, err := fmt.Sscan(string(data), &ui64); err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetUint(ui64)
		return nil
	case reflect.Float32, reflect.Float64:
		f64 := float64(0)
		if _, err := fmt.Sscan(string(data), &f64); err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetFloat(f64)
		return nil
	case reflect.Bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return &ErrUnmarshalFailed{data, rt, err.Error()}
		}
		rv.SetBool(v)
		return nil
	default:
		return &ErrUnmarshalUnsupportedType{Type: rt}
	}
}

func ToBase64(raw []byte) []byte {
	enc := base64.StdEncoding
	length := enc.EncodedLen(len(raw))
	d := make([]byte, length)
	enc.Encode(d, raw)
	return d
}

func FromBase64(data []byte) ([]byte, error) {
	length := base64.StdEncoding.DecodedLen(len(data))
	d := make([]byte, length)
	n, err := base64.StdEncoding.Decode(d, data)
	if err != nil {
		return nil, err
	}
	return d[:n], nil
}

type ErrMarshalUnsupportedType struct {
	Type reflect.Type
}

func (e *ErrMarshalUnsupportedType) Error() string {
	return "marshal text unsupported type `" + e.Type.String() + "`"
}

type ErrUnmarshalUnsupportedType struct {
	Type reflect.Type
}

func (e *ErrUnmarshalUnsupportedType) Error() string {
	return "unmarshal text unsupported type `" + e.Type.String() + "`"
}

type ErrInvalidUnmarshal struct {
	Type reflect.Type
	Err  string
}

func (e *ErrInvalidUnmarshal) Error() string {
	if e.Type == nil {
		return "unmarshal(nil): " + e.Err
	}
	return "unmarshal(`" + e.Type.String() + "`): " + e.Err
}

type ErrUnmarshalFailed struct {
	Data []byte
	Type reflect.Type
	Err  string
}

func (e *ErrUnmarshalFailed) Error() string {
	return "failed unmarshal from `" + string(e.Data) + "` to type `" + e.Type.String() + "`: " + e.Err
}

var rtTextUnmarshaller = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
