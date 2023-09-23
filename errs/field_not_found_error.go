package errs

type FieldNotFoundError struct{}

func (f FieldNotFoundError) Error() string {
	return "field not found"
}
