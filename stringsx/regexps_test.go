package stringsx_test

import (
	"testing"

	. "github.com/xoctopus/x/stringsx"
	. "github.com/xoctopus/x/testx"
)

func TestRegexps(t *testing.T) {
	t.Run("Identifier", func(t *testing.T) {
		for _, v := range []string{
			"PATH",
			"_var1",
			"Xyz_1",
			"_123_abc_",
			"x",
		} {
			Expect(t, ValidIdentifier(v), BeTrue())
		}
		for _, v := range []string{
			"1var",
			"foo-bar",
			"",
		} {
			Expect(t, ValidIdentifier(v), BeFalse())
		}
	})

	t.Run("FlagKey", func(t *testing.T) {
		for _, v := range []string{
			"json",
			"xml",
			"db",
			"env",
			"cmd",
			"x",
		} {
			Expect(t, ValidFlagKey(v), BeTrue())
		}

		for _, v := range []string{
			"",
			"x_y_z",
			"-",
			"\t",
			"ABC",
		} {
			Expect(t, ValidFlagKey(v), BeFalse())
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
			Expect(t, ValidFlagName(v), BeTrue())
		}
		for _, v := range []string{
			" x",
			"1x",
			"\tx7f",
		} {
			Expect(t, ValidFlagName(v), BeFalse())
		}
	})
	t.Run("FlagOptionKey", func(t *testing.T) {
		for _, v := range []string{
			"option",
			"option_x",
			"_",
			"x",
		} {
			Expect(t, ValidFlagOptionKey(v), BeTrue())
		}
		for _, v := range []string{
			"1",
			"",
			"A",
			"-",
		} {
			Expect(t, ValidFlagOptionKey(v), BeFalse())
		}
	})
	t.Run("UnquotedOptionValue", func(t *testing.T) {
		for _, v := range []string{
			"abc",
			"100",
			"XyZ_0123",
			"",
		} {
			Expect(t, ValidUnquotedOptionValue(v), BeTrue())
		}
		for _, v := range []string{
			"1 2 3",
			"1,2,3",
			"\t",
			"\n",
			"\r",
			",",
		} {
			Expect(t, ValidUnquotedOptionValue(v), BeFalse())
		}
	})
}
