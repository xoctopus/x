package stringsx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/stringsx"
)

func TestRegexps(t *testing.T) {
	t.Run("Identifier", func(t *testing.T) {
		for _, v := range []string{
			"PATH",
			"_var1",
			"Xyz_1",
			"_123_abc_",
		} {
			NewWithT(t).Expect(RegexpValidIdentifier.MatchString(v)).To(BeTrue())
		}
		for _, v := range []string{
			"1var",
			"foo-bar",
			"",
		} {
			NewWithT(t).Expect(RegexpValidIdentifier.MatchString(v)).To(BeFalse())
		}
	})

	t.Run("FlagKey", func(t *testing.T) {
		for _, v := range []string{
			"json",
			"xml",
			"db",
			"env",
			"cmd",
		} {
			NewWithT(t).Expect(RegexpValidFlagKey.MatchString(v)).To(BeTrue())
		}

		for _, v := range []string{
			"",
			"x_y_z",
			"-",
			"\t",
			"ABC",
		} {
			NewWithT(t).Expect(RegexpValidFlagKey.MatchString(v)).To(BeFalse())
		}
	})
	t.Run("FlagName", func(t *testing.T) {
		for _, v := range []string{
			"",  // undefine
			"-", // ignore
			"lowerCamel",
			"UpperCamel",
			"UPPER_SNAKE",
			"lower_snake",
			"UPPER-DASH",
			"lower-dash",
			"lower-dash-1",
		} {
			NewWithT(t).Expect(RegexpValidFlagName.MatchString(v)).To(BeTrue())
		}
		for _, v := range []string{
			" x",
			"1x",
			"\tx7f",
		} {
			NewWithT(t).Expect(RegexpValidFlagName.MatchString(v)).To(BeFalse())
		}
	})
	t.Run("FlagOptionKey", func(t *testing.T) {
		for _, v := range []string{
			"option",
			"option_x",
			"_",
		} {
			NewWithT(t).Expect(RegexpValidFlagOptionKey.MatchString(v)).To(BeTrue())
		}
		for _, v := range []string{
			"1",
			"",
			"A",
			"-",
		} {
			NewWithT(t).Expect(RegexpValidFlagOptionKey.MatchString(v)).To(BeFalse())
		}
	})
}
