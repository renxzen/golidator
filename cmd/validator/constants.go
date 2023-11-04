package validator

const (
	NOTBLANK_ERROR = "Must not be blank"
	EMAIL_ERROR = "Must be a valid email"
	URL_ERROR = "Must be a valid url"
	REQUIRED_ERROR = "Must not be missing from body"
	NOTEMPTY_ERROR = "Array must not be empty"

	NOTSTRING_ERROR = "Invalid type. Must be string"
	NOTARRAY_ERROR = "Invalid type. Must be array"
	NOTSTRINGORNUMERIC_ERROR = "Invalid type. Must be string or numeric"

	MIN_STRING_ERROR = "Must have more or equal than %d characters"
	MIN_NUMERIC_ERROR = "Must be more or equal than %d"

	MAX_STRING_ERROR = "Must have less or equal than %d characters"
	MAX_NUMERIC_ERROR = "Must be less or equal than %d"
)
