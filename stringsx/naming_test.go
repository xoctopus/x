package stringsx_test

import (
	"testing"

	. "github.com/xoctopus/x/stringsx"
	. "github.com/xoctopus/x/testx"
)

func TestNaming(t *testing.T) {
	t.Run("Common", func(t *testing.T) {
	})
	name := "i_am_a_10_years_senior"

	Expect(t, LowerCamelCase(name), Equal("iAmA10YearsSenior"))
	Expect(t, LowerSnakeCase(name), Equal("i_am_a_10_years_senior"))
	Expect(t, UpperCamelCase(name), Equal("IAmA10YearsSenior"))
	Expect(t, UpperSnakeCase(name), Equal("I_AM_A_10_YEARS_SENIOR"))
	Expect(t, UpperSnakeCase(name), Equal("I_AM_A_10_YEARS_SENIOR"))
	Expect(t, LowerDashJoint(name), Equal("i-am-a-10-years-senior"))

	Expect(t, UpperCamelCase("OrgID"), Equal("OrgID"))
	Expect(t, LowerCamelCase("OrgID"), Equal("orgId"))
	Expect(t, LowerSnakeCase("OrgID"), Equal("org_id"))
	Expect(t, UpperSnakeCase("OrgID"), Equal("ORG_ID"))
	Expect(t, LowerDashJoint("OrgID"), Equal("org-id"))
}
