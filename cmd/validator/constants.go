package validator

const (
	NOTBLANK_ERROR   = "Must not be blank"
	EMAIL_ERROR      = "Must be a valid email"
	NUMERIC_ERROR    = "Must be a valid string only with number"
	URL_ERROR        = "Must be a valid url"
	REQUIRED_ERROR   = "Must not be missing from body"
	NOTEMPTY_ERROR   = "Array must not be empty"
	LEN_STRING_ERROR = "Must have %d characters"
	LEN_SLICE_ERROR  = "Must have %d elements"

	NOTSTRING_ERROR          = "Invalid type. Must be string"
	NOTARRAY_ERROR           = "Invalid type. Must be array"
	NOTSTRINGORINTEGER_ERROR = "Invalid type. Must be string or integer"
	NOTSTRINGORSLICE_ERROR   = "Invalid type. Must be string or slice"

	MIN_STRING_ERROR  = "Must have more or equal than %d characters"
	MIN_INTEGER_ERROR = "Must be more or equal than %d"

	MAX_STRING_ERROR  = "Must have less or equal than %d characters"
	MAX_INTEGER_ERROR = "Must be less or equal than %d"
)
