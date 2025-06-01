package textx

import (
	"fmt"
	"go/ast"
	"net/url"
	"reflect"

	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/stringsx"
)

// MarshalURL encodes a struct value into url.Values.
// It supports struct tags `name:"customName" as the URL parameter name; It will
// skip unexported fields and fields with name tag "-". By default, the URL
// parameter names are converted to lowerCamelCase.
func MarshalURL(v any) (url.Values, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return url.Values{}, nil
	}

	if rv.Kind() != reflect.Struct {
		return nil, NewErrMarshalURLInvalidInput(v)
	}

	u := url.Values{}
	rt := rv.Type()

	for i := range rv.NumField() {
		sf := rt.Field(i)
		name := stringsx.LowerCamelCase(sf.Name)

		if !ast.IsExported(sf.Name) {
			continue
		}

		flags := reflectx.ParseFlags(sf.Tag)
		if tag := flags.Get("name"); tag != nil {
			if tag.Name == "-" {
				continue
			}
			if tag.Name != "" {
				name = tag.Name
			}
		}

		fi := rv.Field(i)
		if fi.IsZero() {
			if text := sf.Tag.Get("default"); len(text) > 0 {
				u[name] = append(u[name], text)
			}
			continue
		}
		if fi.Kind() == reflect.Slice && !reflectx.IsBytes(fi) {
			for idx := 0; idx < fi.Len(); idx++ {
				text, err := Marshal(fi.Index(idx))
				if err != nil {
					return nil, NewErrMarshalURLFailed(v, fmt.Sprintf("%s[%d]", sf.Name, idx), err)
				}
				u[name] = append(u[name], string(text))
			}
			continue
		}
		text, err := Marshal(rv.Field(i))
		if err != nil {
			return nil, NewErrMarshalURLFailed(v, sf.Name, err)
		}
		u[name] = append(u[name], string(text))
	}

	return u, nil
}

// UnmarshalURL decodes values from url.Values into a struct. The input MUST be
// a pointer to a struct or reflect.Value of a struct.
func UnmarshalURL(u url.Values, v any) error {
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
		return NewErrUnmarshalURLInvalidInput(v)
	}

	if rv.Kind() != reflect.Struct {
		return NewErrUnmarshalURLInvalidInput(v)
	}

	rt := rv.Type()
	for i := range rv.NumField() {
		sf := rt.Field(i)
		name := stringsx.LowerCamelCase(sf.Name)

		if !ast.IsExported(sf.Name) {
			continue
		}

		flags := reflectx.ParseFlags(sf.Tag)
		if tag := flags.Get("name"); tag != nil {
			if tag.Name == "-" {
				continue
			}
			if tag.Name != "" {
				name = tag.Name
			}
		}

		fi := rv.Field(i)
		if fi.Kind() == reflect.Slice && !reflectx.IsBytes(fi) {
			if len(u[name]) == 0 {
				continue
			}
			if fi.IsNil() {
				fi.Set(reflect.MakeSlice(fi.Type(), 0, len(u[name])))
			}
			for idx, text := range u[name] {
				elem := reflect.New(fi.Type().Elem()).Elem()
				if err := Unmarshal([]byte(text), elem); err != nil {
					return NewErrUnmarshalURLFailed(v, fmt.Sprintf("%s[%d]", sf.Name, idx), text, err)
				}
				fi.Set(reflect.Append(fi, elem))
			}
			continue
		}

		text := u.Get(name)
		if text == "" {
			text = sf.Tag.Get("default")
		}
		if err := Unmarshal([]byte(text), rv.Field(i)); err != nil {
			return NewErrUnmarshalURLFailed(v, sf.Name, text, err)
		}
	}
	return nil
}
