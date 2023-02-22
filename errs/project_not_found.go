package errs

type ProjectNotFoundError struct{}

func (p ProjectNotFoundError) Error() string {
	return "project not found"
}
