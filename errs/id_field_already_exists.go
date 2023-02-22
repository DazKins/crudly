package errs

type IdFieldAlreadyExistsError struct{}

func (i IdFieldAlreadyExistsError) Error() string {
	return "id field already exists in table"
}
