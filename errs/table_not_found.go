package errs

type TableNotFoundError struct{}

func (t TableNotFoundError) Error() string {
	return "table not found"
}
