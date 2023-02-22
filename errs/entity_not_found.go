package errs

type EntityNotFoundError struct{}

func (e EntityNotFoundError) Error() string {
	return "entity not found"
}
