package textx

// Ecode defines error code when Marshal/Unmarshal
// +genx:code_error
// @def textx
type Ecode int8

const (
	ECODE_UNDEFINED                     Ecode = iota
	ECODE__MARSHAL_TEXT_INVALID_INPUT         // marshal text got invalid input
	ECODE__MARSHAL_TEXT_FAILED                // failed to marshal text
	ECODE__MARSHAL_URL_INVALID_INPUT          // marshal url got invalid input
	ECODE__MARSHAL_URL_FAILED                 // failed to marshal url
	ECODE__UNMARSHAL_TEXT_INVALID_INPUT       // unmarshal text got invalid input
	ECODE__UNMARSHAL_TEXT_FAILED              // failed to unmarshal text
	ECODE__UNMARSHAL_URL_INVALID_INPUT        // unmarshal url got invalid input
	ECODE__UNMARSHAL_URL_FAILED               // failed to unmarshal url
)
