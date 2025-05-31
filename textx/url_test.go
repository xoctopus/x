package textx_test

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/textx"
	"github.com/xoctopus/x/textx/testdata"
)

var (
	AsErrUnmarshalURLInvalidInput *ErrUnmarshalURLInvalidInput
	AsErrMarshalURLInvalidInput   *ErrMarshalURLInvalidInput
	AsErrMarshalURLFailed         *ErrMarshalURLFailed
	AsErrUnmarshalURLFailed       *ErrUnmarshalURLFailed
)

type Value struct {
	Active     int               `name:"activeTasks" default:"5"`
	IdleTasks  int               `                   default:"3"`
	DB         string            `                   default:"10"`
	Timeout    testdata.Duration `                   default:"100"`
	Ignored    any               `name:"-"           default:"100"`
	NoTag      float32           `                   default:"1.1"`
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
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(u).To(Equal(url.Values{}))

		u, err = MarshalURL((*struct{})(nil))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(u).To(Equal(url.Values{}))

		u, err = MarshalURL(1)
		NewWithT(t).Expect(errors.As(err, &AsErrMarshalURLInvalidInput)).To(BeTrue())
	})

	t.Run("FailedMarshalTest", func(t *testing.T) {
		_, err := MarshalURL(struct{ testdata.MustFailedArshaler }{testdata.MustFailedArshaler{V: 1}})
		NewWithT(t).Expect(errors.As(err, &AsErrMarshalURLFailed)).To(BeTrue())

		_, err = MarshalURL(struct{ V []testdata.MustFailedArshaler }{V: []testdata.MustFailedArshaler{{}}})
		NewWithT(t).Expect(errors.As(err, &AsErrMarshalURLFailed)).To(BeTrue())
	})

	u, err := MarshalURL(DefaultValue)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(u).To(Equal(url.Values{
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
			NewWithT(t).Expect(errors.As(err, &AsErrUnmarshalURLInvalidInput)).To(BeTrue())
		}
	})

	t.Run("UnmarshalDefault", func(t *testing.T) {
		v1 := &Value{}
		err := UnmarshalURL(url.Values{}, v1)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(*v1).To(Equal(DefaultValue))

		v2 := (*Value)(nil)
		err = UnmarshalURL(url.Values{}, &v2)
		NewWithT(t).Expect(*v2).To(Equal(DefaultValue))

		v3 := (**Value)(nil)
		err = UnmarshalURL(url.Values{}, &v3)
		NewWithT(t).Expect(**v3).To(Equal(DefaultValue))
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
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(*v).To(Equal(Value{
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
		NewWithT(t).Expect(err).NotTo(BeNil())
		NewWithT(t).Expect(errors.As(err, &AsErrUnmarshalURLFailed)).To(BeTrue())
		NewWithT(t).Expect(errors.As(err, &AsErrUnmarshalParseFailed)).To(BeTrue())

		v2 := struct{ V []testdata.MustFailedArshaler }{}
		err = UnmarshalURL(url.Values{"v": {""}}, &v2)
		NewWithT(t).Expect(errors.As(err, &AsErrUnmarshalURLFailed)).To(BeTrue())
		NewWithT(t).Expect(errors.As(err, &AsErrUnmarshalFailed)).To(BeTrue())
	})
}

func ExampleMarshalURL() {
	u, err := MarshalURL(struct {
		Name    string `name:"id"`
		Age     int
		Gender  int8
		Country string `name:"country" default:"cn"`
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
