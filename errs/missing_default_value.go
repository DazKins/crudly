package errs

type MissingDefaultValue struct{}

func (m MissingDefaultValue) Error() string {
	return "missing default value"
}
