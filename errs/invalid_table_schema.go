package errs

import "fmt"

type InvalidTableError struct {
	validationError error
}

func NewInvalidTableError(validationError error) InvalidTableError {
	return InvalidTableError{
		validationError,
	}
}

func (i InvalidTableError) Error() string {
	return fmt.Sprintf("table schema is not valid: %s", i.validationError)
}
