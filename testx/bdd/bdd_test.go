package bdd_test

import (
	"strconv"
	"testing"

	"github.com/xoctopus/x/testx/bdd"
)

func TestFeature(t *testing.T) {
	bdd.From(t).
		Given("initial v = 1", func(t bdd.T) {
			v := 1
			t.When("add 1", func(t bdd.T) {
				v += 1
				t.Then("should equal 2", bdd.Equal(2, v))
			})
			t.When("add 2", func(t bdd.T) {
				v += 2
				t.Then("should equal 4", bdd.Equal(4, v))
				t.Then("not equal 3", bdd.NotEqual(3, v))
			})
		})

	v := 0
	bdd.Given(func(t bdd.T) {
		v = 1

		t.Then("the string should equal '1'", bdd.Equal("1", strconv.Itoa(v)))
		t.Then("do nothing")
	})(t)
}
