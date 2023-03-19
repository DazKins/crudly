package errs

import "fmt"

type InvalidEntityOrderError struct {
	validationError error
}

func NewInvalidEntityOrderError(validationError error) InvalidEntityOrderError {
	return InvalidEntityOrderError{
		validationError,
	}
}

func (i InvalidEntityOrderError) Error() string {
	return fmt.Sprintf("entity order is not valid: %s", i.validationError)
}
