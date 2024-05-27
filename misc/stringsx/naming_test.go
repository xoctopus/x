package stringsx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/misc/stringsx"
)

func TestNaming(t *testing.T) {
	name := "i_am_a_10_years_senior"

	NewWithT(t).Expect(LowerCamelCase(name)).To(Equal("iAmA10YearsSenior"))
	NewWithT(t).Expect(LowerSnakeCase(name)).To(Equal("i_am_a_10_years_senior"))
	NewWithT(t).Expect(UpperCamelCase(name)).To(Equal("IAmA10YearsSenior"))
	NewWithT(t).Expect(UpperSnakeCase(name)).To(Equal("I_AM_A_10_YEARS_SENIOR"))
	NewWithT(t).Expect(UpperSnakeCase(name)).To(Equal("I_AM_A_10_YEARS_SENIOR"))
	NewWithT(t).Expect(LowerDashJoint(name)).To(Equal("i-am-a-10-years-senior"))

	NewWithT(t).Expect(UpperCamelCase("OrgID")).To(Equal("OrgID"))
}
