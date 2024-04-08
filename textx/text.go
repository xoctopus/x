package textx

import (
	"encoding"
	"encoding/base64"
	"reflect"
	"strconv"

	"github.com/sincospro/x/reflectx"
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

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		if v == nil {
			return &ErrInvalidUnmarshal{nil}
		}
		return &ErrInvalidUnmarshal{rv.Type()}
	}

	if rv.CanInterface() {
		if unmarshaler, ok := rv.Interface().(encoding.TextUnmarshaler); ok {
			if err := unmarshaler.UnmarshalText(data); err != nil {
				return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
			}
			return nil
		}
	}
	rv = rv.Elem()

	switch rv.Kind() {
	case reflect.Slice:
		if reflectx.IsBytes(rv.Type()) {
			raw, err := FromBase64(data)
			if err != nil {
				return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
			}
			rv.SetBytes(raw)
		} else {
			return &ErrUnmarshalUnsupportedType{rv.Type()}
		}
	case reflect.String:
		rv.SetString(string(data))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return &ErrUnmarshalFailed{data, rv.Type(), err.Error()}
		}
		rv.SetBool(v)
	case reflect.Pointer:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return UnmarshalText(data, rv.Elem().Addr())
	}
	return nil
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
	return "marshal text unsupported type " + e.Type.String()
}

type ErrUnmarshalUnsupportedType struct {
	Type reflect.Type
}

func (e *ErrUnmarshalUnsupportedType) Error() string {
	return "unmarshal text unsupported type " + e.Type.String()
}

type ErrInvalidUnmarshal struct {
	Type reflect.Type
}

func (e *ErrInvalidUnmarshal) Error() string {
	if e.Type == nil {
		return "unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "unmarshal(nil " + e.Type.String() + ")"
}

type ErrUnmarshalFailed struct {
	Data []byte
	Type reflect.Type
	Err  string
}

func (e *ErrUnmarshalFailed) Error() string {
	return "failed unmarshal from " + string(e.Data) + " to type " + e.Type.String() + ": " + e.Err
}
