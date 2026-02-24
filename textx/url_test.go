package textx_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/xoctopus/x/codex"
	. "github.com/xoctopus/x/testx"
	. "github.com/xoctopus/x/textx"
	"github.com/xoctopus/x/textx/testdata"
)

type Value struct {
	Active     int               `url:"activeTasks,default=5"`
	IdleTasks  int               `url:",default=3"`
	DB         string            `url:",default=10"`
	Timeout    testdata.Duration `url:",default=100"`
	Ignored    any               `url:"-,default=100"`
	NoTag      float32           `url:",default='1.1'"`
	Codes      []string
	unexported any
}

var DefaultValue = Value{
	Active:     5,
	IdleTasks:  3,
	DB:         "10",
	Timeout:    100,
	Ignored:    nil,
	NoTag:      1.1,
	unexported: nil,
}

func TestMarshalURL(t *testing.T) {
	t.Run("InvalidInput", func(t *testing.T) {
		u, err := MarshalURL(nil)
		Expect(t, err, BeNil[error]())
		Expect(t, u, Equal(url.Values{}))

		u, err = MarshalURL((*struct{})(nil))
		Expect(t, err, BeNil[error]())
		Expect(t, u, Equal(url.Values{}))

		_, err = MarshalURL(1)
		Expect(t, err, IsError(codex.New(ECODE__MARSHAL_URL_INVALID_INPUT)))
	})

	t.Run("FailedMarshal", func(t *testing.T) {
		_, err := MarshalURL(struct{ testdata.MustFailedArshaler }{testdata.MustFailedArshaler{V: 1}})
		Expect(t, err, IsError(codex.New(ECODE__MARSHAL_URL_FAILED)))

		_, err = MarshalURL(struct{ V []testdata.MustFailedArshaler }{V: []testdata.MustFailedArshaler{{}}})
		Expect(t, err, IsError(codex.New(ECODE__MARSHAL_URL_FAILED)))
	})

	u, err := MarshalURL(DefaultValue)
	Expect(t, err, BeNil[error]())
	Expect(t, u, Equal(url.Values{
		"activeTasks": {"5"},
		"idleTasks":   {"3"},
		"db":          {"10"},
		"timeout":     {"100"},
		"noTag":       {"1.1"},
	}))
}

func TestUnmarshalURL(t *testing.T) {
	t.Run("InvalidInput", func(t *testing.T) {
		for _, v := range []any{nil, new(int)} {
			err := UnmarshalURL(url.Values{}, v)
			Expect(t, err, IsError(codex.New(ECODE__UNMARSHAL_URL_INVALID_INPUT)))
		}
	})

	t.Run("UnmarshalDefault", func(t *testing.T) {
		v1 := &Value{}
		err := UnmarshalURL(url.Values{}, v1)
		Expect(t, err, Succeed())
		Expect(t, *v1, Equal(DefaultValue))

		v2 := (*Value)(nil)
		_ = UnmarshalURL(url.Values{}, &v2)
		Expect(t, *v2, Equal(DefaultValue))

		v3 := (**Value)(nil)
		_ = UnmarshalURL(url.Values{}, &v3)
		Expect(t, **v3, Equal(DefaultValue))
	})

	t.Run("OverwriteDefault", func(t *testing.T) {
		v := &Value{}
		err := UnmarshalURL(url.Values{
			"activeTasks": {"10"},
			"idleTasks":   {"100"},
			"db":          {"database"},
			"timeout":     {"500"},
			"ignored":     {"ignored"},
			"noTag":       {"1.2"},
			"unexported":  {"unexported"},
			"codes":       {"d", "e", "f"},
		}, v)
		Expect(t, err, Succeed())
		Expect(t, *v, Equal(Value{
			Active:    10,
			IdleTasks: 100,
			DB:        "database",
			Timeout:   500,
			Ignored:   nil,
			NoTag:     1.2,
			Codes:     []string{"d", "e", "f"},
		}))
	})

	t.Run("FailedUnmarshalText", func(t *testing.T) {
		v := &Value{}
		err := UnmarshalURL(url.Values{"noTag": {"abc"}}, v)
		Expect(t, err, Failed())
		Expect(t, err, IsError(codex.New(ECODE__UNMARSHAL_URL_FAILED)))

		v2 := struct{ V []testdata.MustFailedArshaler }{}
		err = UnmarshalURL(url.Values{"v": {""}}, &v2)
		Expect(t, err, Failed())
		Expect(t, err, IsError(codex.New(ECODE__UNMARSHAL_URL_FAILED)))
		Expect(t, err, IsError(codex.New(ECODE__UNMARSHAL_TEXT_FAILED)))
	})
}

func ExampleMarshalURL() {
	u, err := MarshalURL(struct {
		Name    string `url:"id"`
		Age     int
		Gender  int8
		Country string `url:"country,default='cn'"`
		Codes   []string
	}{
		Name:  "Alex",
		Age:   30,
		Codes: []string{"a", "b", "c"},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(u.Encode())

	// Output:
	// age=30&codes=a&codes=b&codes=c&country=cn&id=Alex
}

func ExampleSetDefault() {
	fmt.Println("=> need struct value for marshaling")
	_, err := SetDefault(1)
	fmt.Println(err.Error())

	var config = struct {
		K1 int    `url:",default=100"`
		K2 string `url:",default=v2"`
	}{K1: 101}

	fmt.Println("=> initialize a mutable(pointer) struct")
	u, _ := SetDefault(&config)
	fmt.Printf("url parameter values: %s", u.Encode())

	// Output:
	// => need struct value for marshaling
	// [textx:3] marshal url got invalid input. expect struct type
	// => initialize a mutable(pointer) struct
	// url parameter values: k1=101&k2=v2
}
