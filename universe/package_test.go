package universe_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/universe"
)

func TestNewPackage(t *testing.T) {
	cases := []struct {
		path string
		pkg  string
		name string
		id   string
	}{
		{
			path: "github.com/xoctopus/x/internal",
			pkg:  "github.com/xoctopus/x/internal",
			name: "internal",
			id:   "xoctopus__github__com_xoctopus_x_internal__xoctopus",
		},
		{
			path: "github.com/xoctopus/x/internal",
			pkg:  "github.com/xoctopus/x/internal",
			name: "internal",
			id:   "xoctopus__github__com_xoctopus_x_internal__xoctopus",
		},
		{
			path: "xoctopus__github__com_xoctopus_x_internal__xoctopus",
			pkg:  "github.com/xoctopus/x/internal",
			name: "internal",
			id:   "xoctopus__github__com_xoctopus_x_internal__xoctopus",
		},
		{
			path: "encoding/json",
			pkg:  "encoding/json",
			name: "json",
			id:   "xoctopus__encoding_json__xoctopus",
		},
		{
			path: "net",
			pkg:  "net",
			name: "net",
			id:   "net",
		},
		{
			path: "xoctopus__net__xoctopus",
			pkg:  "net",
			name: "net",
			id:   "net",
		},
	}

	for _, c := range cases {
		p := NewPackage(c.path)
		NewWithT(t).Expect(p.Name()).To(Equal(c.name))
		NewWithT(t).Expect(p.Path()).To(Equal(c.pkg))
		NewWithT(t).Expect(p.ID()).To(Equal(c.id))
	}

	NewWithT(t).Expect(NewPackage("")).To(BeNil())
}
