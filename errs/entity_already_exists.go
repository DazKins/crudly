package errs

type EntityAlreadyExistsError struct{}

func (e EntityAlreadyExistsError) Error() string {
	return "entity already exists"
}
