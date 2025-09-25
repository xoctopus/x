package stringsx

import "regexp"

var (
	regexValidIdentifier    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	regexValidFlagKey       = regexp.MustCompile(`^[a-z]+$`)
	regexValidFlagName      = regexp.MustCompile(`^(?:-|[A-Za-z](?:[A-Za-z0-9_-]*[A-Za-z0-9])?)?$`)
	regexValidFlagOptionKey = regexp.MustCompile(`^[a-z_]+$`)
)

func ValidIdentifier(s string) bool {
	return regexValidIdentifier.MatchString(s)
}

func ValidFlagKey(s string) bool {
	return regexValidFlagKey.MatchString(s)
}

func ValidFlagName(s string) bool {
	return regexValidFlagName.MatchString(s)
}

func ValidFlagOptionKey(s string) bool {
	return regexValidFlagOptionKey.MatchString(s)
}
