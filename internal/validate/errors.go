package validate

const (
	ErrMsgNotBlank           = "Must not be blank"
	ErrMsgInvalidEmail       = "Must be a valid email"
	ErrMsgNotNumeric         = "Must be a valid string with numbers only"
	ErrMsgInvalidURL         = "Must be a valid url"
	ErrMsgMissing            = "Must not be missing from body"
	ErrMsgEmptyArray         = "Array must not be empty"
	ErrMsgInvalidLength      = "Must have %d characters"
	ErrMsgInvalidLengthSlice = "Must have %d elements"

	ErrMsgNotStringType   = "Invalid type. Must be string"
	ErrMsgNotArrayType    = "Invalid type. Must be array"
	ErrMsgNotStrIntType   = "Invalid type. Must be string or integer"
	ErrMsgNotStrSliceType = "Invalid type. Must be string or slice"

	ErrMsgStrInvalidMin = "Must have more or equal than %d characters"
	ErrMsgStrInvalidInt = "Must be more or equal than %d"

	ErrMsgStrInvalidMax = "Must have less or equal than %d characters"
	ErrMsgIntInvalidMax = "Must be less or equal than %d"
)
