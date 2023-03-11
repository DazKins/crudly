package errs

import "fmt"

type InvalidPartialEntityError struct {
	validationError error
}

func NewInvalidPartialEntityError(validationError error) InvalidPartialEntityError {
	return InvalidPartialEntityError{
		validationError,
	}
}

func (i InvalidPartialEntityError) Error() string {
	return fmt.Sprintf("partial entity is not valid: %s", i.validationError)
}
