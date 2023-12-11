package repo

type IBaseRepo interface {
	// TranslateError translates error from the database to
	// a more readable and understandable application error.
	TranslateError(err error) error
}
