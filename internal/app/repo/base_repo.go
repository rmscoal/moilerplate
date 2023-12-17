package repo

type IBaseRepo interface {
	// DetectConstraintError detects whether the error from the database is a
	// constraint error. If so, it translates the error to a more readable error.
	// On the other hand, returns an unexpected error.
	DetectConstraintError(err error) error

	// DetectNotFoundError detects whether the error from the database is a
	// not found error. If so, it translates the error to a more readable error
	// On the other hand, returns an unexpected error.
	DetectNotFoundError(err error) error
}
