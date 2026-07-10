package textx_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/xoctopus/x/testx/bdd"
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
	Inlined    `url:",inline"`
	unexported any
}

type Inlined struct {
	K1 int    `url:",default=100"`
	K2 string `url:",default=v2"`
}

var DefaultValue = Value{
	Active:     5,
	IdleTasks:  3,
	DB:         "10",
	Timeout:    100,
	Ignored:    nil,
	NoTag:      1.1,
	unexported: nil,
	Inlined: Inlined{
		K1: 100,
		K2: "v2",
	},
}

func TestURLArshaler(t *testing.T) {
	bdd.From(t).When("marshaling", func(t bdd.T) {
		t.Given("nil", func(t bdd.T) {
			u, err := MarshalURL(nil)
			t.Then("success and got empty url",
				bdd.Succeed(err),
				bdd.Equal(u, url.Values{}),
			)
		})

		t.Given("not struct input", func(t bdd.T) {
			_, err := MarshalURL(1)
			t.Then("got invalid input error",
				bdd.IsCodeError(err, ECODE__MARSHAL_URL_INVALID_INPUT),
			)
		})

		t.Given("input must marshal failed", func(t bdd.T) {
			_, err1 := MarshalURL(struct {
				testdata.MustFailedArshaler
			}{
				MustFailedArshaler: testdata.MustFailedArshaler{V: 1},
			})
			_, err2 := MarshalURL(struct {
				V []testdata.MustFailedArshaler
			}{
				V: []testdata.MustFailedArshaler{{}},
			})
			t.Then("got marshal failed error",
				bdd.IsCodeError(err1, ECODE__MARSHAL_URL_FAILED),
				bdd.IsCodeError(err2, ECODE__MARSHAL_URL_FAILED),
			)
		})

		t.Given("value case for marshaling", func(t bdd.T) {
			u, err := MarshalURL(DefaultValue)
			t.Then("success and match expect",
				bdd.Succeed(err),
				bdd.Equal(u, url.Values{
					"activeTasks": {"5"},
					"idleTasks":   {"3"},
					"db":          {"10"},
					"timeout":     {"100"},
					"noTag":       {"1.1"},
					"k1":          {"100"},
					"k2":          {"v2"},
				}),
			)
		})
	})

	bdd.From(t).When("unmarshaling", func(t bdd.T) {
		t.Given("nil", func(t bdd.T) {
			err := UnmarshalURL(url.Values{}, nil)
			t.Then("got unmarshal failed error",
				bdd.IsCodeError(err, ECODE__UNMARSHAL_URL_INVALID_INPUT),
			)
		})

		t.Given("not struct value", func(t bdd.T) {
			err := UnmarshalURL(url.Values{}, nil)
			t.Then("got unmarshal failed error",
				bdd.IsCodeError(err, ECODE__UNMARSHAL_URL_INVALID_INPUT),
			)
		})

		t.Given("empty value pointer case", func(t bdd.T) {
			v1 := &Value{}
			err1 := UnmarshalURL(url.Values{}, v1)
			v2 := new(Value)
			err2 := UnmarshalURL(url.Values{}, v2)
			v3 := (**Value)(nil)
			err3 := UnmarshalURL(url.Values{}, &v3)
			t.Then("unmarshal succeed and equal default result",
				bdd.Succeed(err1),
				bdd.Equal(*v1, DefaultValue),
				bdd.Succeed(err2),
				bdd.Equal(*v2, DefaultValue),
				bdd.Succeed(err3),
				bdd.Equal(**v3, DefaultValue),
			)
		})

		t.Given("valued url.Values", func(t bdd.T) {
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
				"k1":          {"101"},
				"k2":          {"v3"},
			}, v)

			t.Then("should succeed and overwrite default value",
				bdd.Succeed(err),
				bdd.Equal(*v, Value{
					Active:    10,
					IdleTasks: 100,
					DB:        "database",
					Timeout:   500,
					Ignored:   nil,
					NoTag:     1.2,
					Codes:     []string{"d", "e", "f"},
					Inlined:   Inlined{K1: 101, K2: "v3"},
				}),
			)
		})

		t.Given("field unmarshal failed", func(t bdd.T) {
			v1 := &Value{}
			err1 := UnmarshalURL(url.Values{"noTag": {"abc"}}, v1)
			v2 := struct{ V []testdata.MustFailedArshaler }{}
			err2 := UnmarshalURL(url.Values{"v": {""}}, &v2)

			t.Then("got unmarshal failed error",
				bdd.IsCodeError(err1, ECODE__UNMARSHAL_URL_FAILED),
				bdd.IsCodeError(err2, ECODE__UNMARSHAL_URL_FAILED),
				bdd.IsCodeError(err2, ECODE__UNMARSHAL_TEXT_FAILED),
			)
		})
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
