package errs

type UserAlreadyExistsError struct{}

func (u UserAlreadyExistsError) Error() string {
	return "user already exists"
}
