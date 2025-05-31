package textx

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"github.com/xoctopus/x/reflectx"
)

var (
	TextMarshaler   = reflect.TypeFor[encoding.TextMarshaler]()
	TextUnmarshaler = reflect.TypeFor[encoding.TextUnmarshaler]()
)

func Marshal(v any) ([]byte, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return nil, nil
	}

	rt := rv.Type()

	if rt.Implements(TextMarshaler) {
		output, err := rv.Interface().(encoding.TextMarshaler).MarshalText()
		return output, NewErrMarshalFailed(v, err)
	}

	output := make([]byte, 0, 8)

	switch kind := rv.Kind(); kind {
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(output, rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(output, rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.AppendFloat(output, rv.Float(), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.AppendFloat(output, rv.Float(), 'g', -1, 64), nil
	case reflect.Bool:
		return strconv.AppendBool(output, rv.Bool()), nil
	default:
		if reflectx.IsBytes(rv) {
			return rv.Bytes(), nil
		}
		return nil, NewErrMarshalUnsupportedType(v)
	}
}

func Unmarshal(data []byte, v any) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() && rv.CanSet() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	if !rv.CanSet() {
		return NewErrUnmarshalInvalidInput(v)
	}

	rt := rv.Type()
	if reflect.PointerTo(rt).Implements(TextUnmarshaler) {
		err := rv.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText(data)
		return NewErrUnmarshalFailed(data, v, err)
	}

	switch rv.Kind() {
	case reflect.String:
		rv.SetString(string(data))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := int64(0)
		if _, err := fmt.Sscan(string(data), &i64); err != nil {
			return NewErrUnmarshalParseFailed(data, v, err)
		}
		rv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui64 := uint64(0)
		if _, err := fmt.Sscan(string(data), &ui64); err != nil {
			return NewErrUnmarshalParseFailed(data, v, err)
		}
		rv.SetUint(ui64)
		return nil
	case reflect.Float32, reflect.Float64:
		f64 := float64(0)
		if _, err := fmt.Sscan(string(data), &f64); err != nil {
			return NewErrUnmarshalParseFailed(data, v, err)
		}
		rv.SetFloat(f64)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(string(data))
		if err != nil {
			return NewErrUnmarshalParseFailed(data, v, err)
		}
		rv.SetBool(b)
		return nil
	default:
		if reflectx.IsBytes(rv.Type()) {
			rv.SetBytes(data)
			return nil
		}
		return NewErrUnmarshalUnsupportedType(v)
	}
}
