package store

// NOT_FOUND_ERR is the error message for "value not found in store"
const NOT_FOUND_ERR string = "value not found in store"

// NotFound is an error type representing a "value not found" error.
type NotFound struct {
	cause error
}

// NotFoundWithCause creates a new NotFound error with a cause.
func NotFoundWithCause(e error) error {
	err := NotFound{
		cause: e,
	}
	return &err
}

// Cause returns the underlying cause of the NotFound error.
func (e NotFound) Cause() error {
	return e.cause
}

// Is checks if the given error is a NotFound error.
func (e NotFound) Is(err error) bool {
	return err.Error() == NOT_FOUND_ERR
}

// Error returns the error message for the NotFound error.
func (e NotFound) Error() string {
	return NOT_FOUND_ERR
}

// Unwrap returns the underlying cause of the NotFound error.
func (e NotFound) Unwrap() error {
	return e.cause
}
