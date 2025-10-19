package reflectx

// Ecode defines error codes for tag parsing
// +genx:code_error
// @def PARSING_TAG
type Ecode int8

const (
	ECODE_UNDEFINED                Ecode = iota
	ECODE__INVALID_FLAG_KEY              // invalid flag key
	ECODE__INVALID_FLAG_VALUE            // invalid flag value
	ECODE__INVALID_FLAG_NAME             // invalid flag flag name
	ECODE__INVALID_OPTION_KEY            // invalid option key
	ECODE__INVALID_OPTION_VALUE          // invalid option value
	ECODE__INVALID_OPTION_UNQUOTED       // invalid option unquoted
)
