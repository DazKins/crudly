package errs

import "fmt"

type InvalidEntityError struct {
	validationError error
}

func NewInvalidEntityError(validationError error) InvalidEntityError {
	return InvalidEntityError{
		validationError,
	}
}

func (i InvalidEntityError) Error() string {
	return fmt.Sprintf("entity is not valid: %s", i.validationError)
}
