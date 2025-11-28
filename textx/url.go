package textx

import (
	"go/ast"
	"net/url"
	"reflect"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/stringsx"
)

// URLTag describes url parameter in struct tag
const URLTag = "url"

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
		return nil, codex.Errorf(ECODE__MARSHAL_URL_INVALID_INPUT, "expect struct type")
	}

	u := url.Values{}
	rt := rv.Type()

	for i := range rv.NumField() {
		sf := rt.Field(i)
		name := stringsx.LowerCamelCase(sf.Name)

		if !ast.IsExported(sf.Name) {
			continue
		}

		flag := reflectx.ParseTag(sf.Tag).Get(URLTag)

		if flag != nil {
			if flag.Name() == "-" {
				continue
			}
			if flag.Name() != "" {
				name = flag.Name()
			}
		}

		fi := rv.Field(i)
		if fi.IsZero() {
			if flag != nil {
				if opt := flag.Option("default"); opt != nil {
					if text := opt.Unquoted(); len(text) > 0 {
						u[name] = append(u[name], text)
					}
				}
			}
			continue
		}
		if fi.Kind() == reflect.Slice && !reflectx.IsBytes(fi) {
			for idx := 0; idx < fi.Len(); idx++ {
				text, err := Marshal(fi.Index(idx))
				if err != nil {
					return nil, codex.Wrapf(ECODE__MARSHAL_URL_FAILED, err, "field %s[%d]", sf.Name, idx)
				}
				u[name] = append(u[name], string(text))
			}
			continue
		}
		text, err := Marshal(rv.Field(i))
		if err != nil {
			return nil, codex.Wrapf(ECODE__MARSHAL_URL_FAILED, err, "field %s", sf.Name)
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
		return codex.Errorf(ECODE__UNMARSHAL_URL_INVALID_INPUT, "must canbe set")
	}

	if rv.Kind() != reflect.Struct {
		return codex.Errorf(ECODE__UNMARSHAL_URL_INVALID_INPUT, "expect struct type")
	}

	rt := rv.Type()
	for i := range rv.NumField() {
		sf := rt.Field(i)
		name := stringsx.LowerCamelCase(sf.Name)

		if !ast.IsExported(sf.Name) {
			continue
		}

		flag := reflectx.ParseTag(sf.Tag).Get(URLTag)
		if flag != nil {
			if flag.Name() == "-" {
				continue
			}
			if flag.Name() != "" {
				name = flag.Name()
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
					return codex.Wrapf(ECODE__UNMARSHAL_URL_FAILED, err, "parse %s[%d] from `%s`", sf.Name, idx, text)
				}
				fi.Set(reflect.Append(fi, elem))
			}
			continue
		}

		text := u.Get(name)
		if text == "" && flag != nil {
			if opt := flag.Option("default"); opt != nil {
				text = opt.Unquoted()
			}
		}
		if len(text) > 0 {
			if err := Unmarshal([]byte(text), rv.Field(i)); err != nil {
				return codex.Wrapf(ECODE__UNMARSHAL_URL_FAILED, err, "parse field %s from %s", sf.Name, text)
			}
		}
	}
	return nil
}
