package bdd_test

import (
	"strconv"
	"testing"

	"github.com/xoctopus/x/testx/bdd"
)

func TestFeature(t *testing.T) {
	bdd.From(t).Given("ValueInitializedWith1", func(t bdd.T) {
		v := 1
		t.When("Add1", func(t bdd.T) {
			v += 1
			t.Then("Equal2", bdd.Equal(2, v))
		})
		t.When("ThenAdd2", func(t bdd.T) {
			v += 2
			t.Then("Equal4", bdd.Equal(4, v))
			t.Then("NotEqual3", bdd.NotEqual(3, v))
		})
	})

	t.Run("GIVEN v=1", bdd.Given(func(t bdd.T) {
		v := 1
		t.Then("StringIs", bdd.Equal("1", strconv.Itoa(v)))
	}))
}
