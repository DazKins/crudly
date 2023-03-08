package errs

import "fmt"

type InvalidEntityFilterError struct {
	validationError error
}

func NewInvalidEntityFilterError(validationError error) InvalidEntityFilterError {
	return InvalidEntityFilterError{
		validationError,
	}
}

func (i InvalidEntityFilterError) Error() string {
	return fmt.Sprintf("entity filter is not valid: %s", i.validationError)
}
