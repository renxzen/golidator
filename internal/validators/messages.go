package validators

const (
	MessageNotBlank           = "must not be blank"
	MessageInvalidEmail       = "must be a valid email"
	MessageNotNumeric         = "must be a valid string with numbers only"
	MessageInvalidURL         = "must be a valid url"
	MessageMissing            = "must not be missing from body"
	MessageEmptyArray         = "array must not be empty"
	MessageInvalidLength      = "must have %d characters"
	MessageInvalidLengthSlice = "must have %d elements"
	MessageNotStringType      = "invalid type. must be string"
	MessageNotArrayType       = "invalid type. must be array"
	MessageNotStrIntType      = "invalid type. must be string or integer"
	MessageNotStrSliceType    = "invalid type. must be string or slice"
	MessageStrInvalidMin      = "must have more or equal than %d characters"
	MessageStrInvalidInt      = "must be more or equal than %d"
	MessageStrInvalidMax      = "must have less or equal than %d characters"
	MessageIntInvalidMax      = "must be less or equal than %d"
)
