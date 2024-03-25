package slice_test

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/sincospro/x/misc/slice"
)

func TestToAnySlice(t *testing.T) {
	strings := []string{"a", "b"}
	gomega.NewWithT(t).Expect(slice.ToAnySlice(strings...)).
		To(gomega.Equal([]interface{}{"a", "b"}))
}
