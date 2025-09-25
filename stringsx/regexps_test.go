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
			NewWithT(t).Expect(ValidIdentifier(v)).To(BeTrue())
		}
		for _, v := range []string{
			"1var",
			"foo-bar",
			"",
		} {
			NewWithT(t).Expect(ValidIdentifier(v)).To(BeFalse())
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
			NewWithT(t).Expect(ValidFlagKey(v)).To(BeTrue())
		}

		for _, v := range []string{
			"",
			"x_y_z",
			"-",
			"\t",
			"ABC",
		} {
			NewWithT(t).Expect(ValidFlagKey(v)).To(BeFalse())
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
			"b",
		} {
			NewWithT(t).Expect(ValidFlagName(v)).To(BeTrue())
		}
		for _, v := range []string{
			" x",
			"1x",
			"\tx7f",
		} {
			NewWithT(t).Expect(ValidFlagName(v)).To(BeFalse())
		}
	})
	t.Run("FlagOptionKey", func(t *testing.T) {
		for _, v := range []string{
			"option",
			"option_x",
			"_",
		} {
			NewWithT(t).Expect(ValidFlagOptionKey(v)).To(BeTrue())
		}
		for _, v := range []string{
			"1",
			"",
			"A",
			"-",
		} {
			NewWithT(t).Expect(ValidFlagOptionKey(v)).To(BeFalse())
		}
	})
}
