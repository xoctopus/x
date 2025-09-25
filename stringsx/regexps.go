package stringsx

import "regexp"

var (
	RegexpValidIdentifier    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	RegexpValidFlagKey       = regexp.MustCompile(`^[a-z]+$`)
	RegexpValidFlagName      = regexp.MustCompile(`^(?:-|[A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])?$`)
	RegexpValidFlagOptionKey = regexp.MustCompile(`^[a-z_]+$`)
)
