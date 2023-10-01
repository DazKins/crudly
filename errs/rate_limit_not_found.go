package errs

type RateLimitNotFoundError struct{}

func (r RateLimitNotFoundError) Error() string {
	return "rate limit not found"
}
